package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/launchboxio/cript/internal/admission"
	"github.com/spf13/cobra"
	admissionv1 "k8s.io/api/admission/v1"
	"net/http"
)

var (
	webhookCmd = &cobra.Command{
		Use:   "webhook",
		Short: "Run a webhook server for validation scans",
		Run:   webhookServer,
	}
)

func init() {
	webhookCmd.Flags().Int16("port", 8080, "The port to serve webhooks on")
}

func validatePod(w http.ResponseWriter, r *http.Request) {
	in, err := parseRequest(*r)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	adm := admission.Admission{
		Logger:  logger,
		Request: in.Request,
	}

	out, err := adm.RunPodReview()
	if err != nil {
		e := fmt.Sprintf("could not generate admission response: %v", err)
		logger.Error(e)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jout, err := json.Marshal(out)
	if err != nil {
		e := fmt.Sprintf("could not parse admission response: %v", err)
		logger.Error(e)
		http.Error(w, e, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", jout)
}

func webhookServer(cmd *cobra.Command, args []string) {
	port, _ := cmd.Flags().GetInt8("port")
	http.HandleFunc("/pods", validatePod)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		logger.Fatal(err)
	}
}

func parseRequest(r http.Request) (*admissionv1.AdmissionReview, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("Content-Type: %q should be %q",
			r.Header.Get("Content-Type"), "application/json")
	}

	bodybuf := new(bytes.Buffer)
	bodybuf.ReadFrom(r.Body)
	body := bodybuf.Bytes()

	if len(body) == 0 {
		return nil, fmt.Errorf("admission request body is empty")
	}

	var a admissionv1.AdmissionReview

	if err := json.Unmarshal(body, &a); err != nil {
		return nil, fmt.Errorf("could not parse admission review request: %v", err)
	}

	if a.Request == nil {
		return nil, fmt.Errorf("admission review can't be used: Request field is nil")
	}

	return &a, nil
}
