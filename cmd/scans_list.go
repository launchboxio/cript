package cmd

import (
	"context"
	"fmt"
	"github.com/launchboxio/cript/api/v1alpha1"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	scansListCmd = &cobra.Command{
		Use:   "scans list",
		Short: "Show current scans in cluster",
		Run:   runScansList,
	}
)

func init() {
	scansListCmd.Flags().String("kubeconfig", "", "Path to kubeconfig file")

	rootCmd.AddCommand(scansListCmd)
}

func runScansList(cmd *cobra.Command, args []string) {
	fmt.Println(conf)
	var config *rest.Config
	var err error

	kubeconfig, _ := cmd.Flags().GetString("kubeconfig")

	if kubeconfig == "" {
		logger.Info("using in-cluster configuration")
		config, err = rest.InClusterConfig()
	} else {
		logger.Infof("using configuration from '%s'", kubeconfig)
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		logger.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatal(err)
	}

	result := v1alpha1.ScanList{}

	d, err := clientset.RESTClient().Get().AbsPath("/apis/security.cript.dev/v1alpha1/scans").DoRaw(context.TODO())
	if err != nil {
		logger.Fatal(err)
	}
	if err = json.Unmarshal(d, &result); err != nil {
		logger.Fatal(err)
	}

	fmt.Println(result)
}
