package cmd

import (
	"bytes"
	"context"
	"github.com/anchore/grype/grype/presenter/models"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/launchboxio/cript/internal/storage"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/json"
	"os/exec"
	"strings"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan an image against a declaration",
	Run:   runScan,
}

func init() {
	scanCmd.Flags().String("image", "", "Container image to scan")
	scanCmd.Flags().String("bucket", "image-scans", "S3 bucket for storing raw scan output")
	scanCmd.MarkFlagRequired("image")
}

func runScan(cmd *cobra.Command, args []string) {

	image, _ := cmd.Flags().GetString("image")

	store, err := storage.NewForConfig(conf)
	if err != nil {
		logger.Fatal(err)
	}

	chunks := strings.Split(image, ":")
	if len(chunks) == 1 || chunks[1] == "latest" {
		logger.Info("We're scanning a latest image.")
		logger.Info("Figure out what we want to do here...")
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("Pulling image %s\n", image)

	_, err = cli.ImagePull(context.TODO(), image, types.ImagePullOptions{})
	if err != nil {
		logger.Fatal(err)
	}
	_, res, err := getGrypeInformation(image)
	if err != nil {
		logger.Fatalf("Grype scan failed: %v\n", err)
	}
	err = store.StoreVulnerabilityReport(image, res)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Scan completed. Starting inspection")

	inspect, _, err := cli.ImageInspectWithRaw(context.TODO(), image)
	if err != nil {
		logger.Fatal(err)
	}

	err = store.StoreInspection(image, inspect)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Analysis complete")
}

func getGrypeInformation(image string) (string, models.Document, error) {
	var outb, errb bytes.Buffer
	var res models.Document

	grypeCmd := exec.Command("grype", image, "-q", "-o=json")
	grypeCmd.Stdout = &outb
	grypeCmd.Stderr = &errb

	err := grypeCmd.Run()
	if err != nil {
		return "", res, err
	}

	err = json.Unmarshal(outb.Bytes(), &res)
	return outb.String(), res, err
}
