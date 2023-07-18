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

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type DeclarationRef struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// ScanSpec defines the desired state of Scan
type ScanSpec struct {
	ImageUri string `json:"imageUri,omitempty"`

	// Reference to the declaration to validate packages during a scan
	DeclarationRef DeclarationRef `json:"declarationRef,omitempty"`

	// For private images, we'll need a secretRef to use for pulling them
	ImagePullSecretsRef v1.SecretReference `json:"imagePullSecretsRef,omitempty"`
}

// ScanStatus defines the observed state of Scan
type ScanStatus struct {
	State string `json:"state,omitempty"`
	Job   string `json:"job,omitempty"`

	//Vulnerabilities []interface{} `json:"vulnerabilities,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Scan is the Schema for the scans API
type Scan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScanSpec   `json:"spec,omitempty"`
	Status ScanStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ScanList contains a list of Scan
type ScanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Scan `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Scan{}, &ScanList{})
}
