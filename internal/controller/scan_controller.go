/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	securityv1alpha1 "github.com/launchboxio/cript/api/v1alpha1"
)

// ScanReconciler reconciles a Scan object
type ScanReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=security.cript.dev,resources=scans,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=security.cript.dev,resources=scans/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=security.cript.dev,resources=scans/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Scan object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ScanReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	scan := &securityv1alpha1.Scan{}

	// Fetch our CRD resource
	err := r.Get(ctx, req.NamespacedName, scan)
	if err != nil {
		if errors.IsNotFound(err) {
			ctrl.Log.Info("Scan resource not found")
			return ctrl.Result{}, nil
		}
		ctrl.Log.Error(err, "Failed getting scan resource")
		return ctrl.Result{}, err
	}

	// Kick off the job to get scan output
	found := &batchv1.Job{}
	err = r.Get(ctx, types.NamespacedName{Name: scan.Name, Namespace: scan.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new deployment
		job := r.jobForScan(scan)
		err = r.Create(ctx, job)
		if err != nil {
			ctrl.Log.Error(err, "Failed creating new job resource")
			return ctrl.Result{}, err
		}
		// Deployment created successfully - return and requeue
		ctrl.Log.Info("Job created, requeueing")
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		ctrl.Log.Error(err, "Failed getting job resource")
		return ctrl.Result{}, err
	}

	// Handle running status of the job

	// Once the job has completed (successfully), analyze the report
	if found.Status.Succeeded == 0 {
		// Continue polling, waiting for completion
		ctrl.Log.Info("Waiting for job to complete")
		return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
	}

	ctrl.Log.Info("Image scan completed, we would now run analysis")

	// Job succeeded
	// Update the scan status with the output from the report

	// Grab the data from the store, and cache it in our CRD as well

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ScanReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&securityv1alpha1.Scan{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}

func (r *ScanReconciler) jobForScan(scan *securityv1alpha1.Scan) *batchv1.Job {
	parallelism := int32(1)
	// TODO: Move to -rootless, and remove securityContext.privileged: true
	privileged := true
	// TODO: Mount imagePullSecrets if needed
	// TODO: Mount the config file for cript
	volumes := []corev1.Volume{{
		Name: "docker-graph-storage",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}, {
		Name: "dind-certs",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}}
	serviceAccount := ""
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        scan.Name,
			Namespace:   scan.Namespace,
			Labels:      map[string]string{},
			Annotations: map[string]string{},
		},
		Spec: batchv1.JobSpec{
			Parallelism: &parallelism,
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Volumes:        volumes,
					InitContainers: []corev1.Container{},
					Containers: []corev1.Container{{
						Name:            "scanner",
						Image:           "docker.io/library/cript:latest",
						ImagePullPolicy: "Never",
						Command:         []string{"/cript", "scan"},
						Args: []string{
							fmt.Sprintf("--image=%s", scan.Spec.ImageUri),
						},
						Env: []corev1.EnvVar{{
							Name:  "DOCKER_HOST",
							Value: "tcp://localhost:2376",
						}, {
							Name:  "DOCKER_TLS_VERIFY",
							Value: "1",
						}, {
							Name:  "DOCKER_CERT_PATH",
							Value: "/certs/client",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "dind-certs",
							MountPath: "/certs/client",
						}},
					}, {
						Name:  "dind-daemon",
						Image: "docker:24.0.4-dind",
						SecurityContext: &corev1.SecurityContext{
							Privileged: &privileged,
						},
						Args: []string{
							"--userland-proxy=false",
							"--debug",
						},
						Env: []corev1.EnvVar{{
							Name:  "DOCKER_TLS_CERTDIR",
							Value: "/certs",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "docker-graph-storage",
							MountPath: "/var/lib/docker",
						}, {
							Name:      "dind-certs",
							MountPath: "/certs/client",
						}},
						ReadinessProbe: &corev1.Probe{
							PeriodSeconds: int32(1),
							ProbeHandler: corev1.ProbeHandler{
								Exec: &corev1.ExecAction{
									Command: []string{
										"ls",
										"/certs/client/ca.pem",
									},
								},
							},
						},
					}},
					RestartPolicy:      corev1.RestartPolicyOnFailure,
					ServiceAccountName: serviceAccount,
				},
			},
		},
	}
}
