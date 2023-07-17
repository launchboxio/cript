package admission

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Admission struct {
	Logger  *zap.SugaredLogger
	Request *admissionv1.AdmissionRequest
}

func (a *Admission) RunPodReview() *admissionv1.AdmissionReview {
	res := &admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AdmissionReview",
			APIVersion: "admission.k8s.io/v1",
		},
		Response: &admissionv1.AdmissionResponse{
			UID:     a.Request.UID,
			Allowed: true,
			Result: &metav1.Status{
				Code:    200,
				Message: "No vulnerabilities found",
			},
		},
	}
	pod, err := a.Pod()
	if err != nil {
		res.Response.Allowed = false
		res.Response.Result = &metav1.Status{
			Code:    400,
			Message: "Payload was invalid",
		}
		return res
	}

	images := []string{}
	for _, container := range pod.Spec.Containers {
		images = append(images, container.Image)
	}
	for _, initContainer := range pod.Spec.InitContainers {
		images = append(images, initContainer.Image)
	}

	// Verify all images are stored, scanned, and have passed
	// any required validation checks
	//for image := range images {
	// Figure out what to do here if image is ":latest"
	// Get the CRD
	// If image hasn't been scanned yet, return error
	// If image isn't approved, return error as well
	//res.Response.Allowed = false
	//res.Response.Result = &metav1.Status{
	//	Code:   400,
	//	Message: fmt.Sprintf("Image %s hasn't been scanned", image),
	//}
	//}

	return res
}

func (a Admission) Pod() (*corev1.Pod, error) {
	if a.Request.Kind.Kind != "Pod" {
		return nil, fmt.Errorf("only pods are supported here")
	}

	p := corev1.Pod{}
	if err := json.Unmarshal(a.Request.Object.Raw, &p); err != nil {
		return nil, err
	}

	return &p, nil
}
