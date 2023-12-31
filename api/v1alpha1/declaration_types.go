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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DeclarationSpec defines the desired state of Declaration
type DeclarationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Vulnerabilities to ignore
	Ignore []string `json:"ignore,omitempty"`

	// Criteria for when a scan should be marked as failed
	FailurePolicy string `json:"failurePolicy,omitempty"`

	Manifest Manifest `json:"manifest,omitempty"`
}

type Manifest struct {
	Rules []ManifestRule `json:"rules,omitempty"`
}

type ManifestRule struct {
	Key         string   `json:"key"`
	Equals      string   `json:"equals,omitempty"`
	OneOf       []string `json:"oneOf,omitempty"`
	ArrayEquals []string `json:"arrayEquals,omitempty"`
}

// DeclarationStatus defines the observed state of Declaration
type DeclarationStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Declaration is the Schema for the declarations API
type Declaration struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeclarationSpec   `json:"spec,omitempty"`
	Status DeclarationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DeclarationList contains a list of Declaration
type DeclarationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Declaration `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Declaration{}, &DeclarationList{})
}
