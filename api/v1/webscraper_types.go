/*
Copyright 2024.

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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WebScraperSpec defines the desired state of WebScraper
type WebScraperSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Schedule          string                      `json:"schedule,omitempty"`
	Image             string                      `json:"image,omitempty"`
	Command           []string                    `json:"command,omitempty"`
	Retries           int32                       `json:"retries,omitempty"`
	ConcurrencyPolicy string                      `json:"concurrencyPolicy,omitempty"`
	Resources         corev1.ResourceRequirements `json:"resources,omitempty"`
}

// WebScraperStatus defines the observed state of WebScraper
type WebScraperStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	LastRunTime metav1.Time `json:"lastRunTime,omitempty"`
	Success     bool        `json:"success,omitempty"`
	Message     string      `json:"message,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// WebScraper is the Schema for the webscrapers API
type WebScraper struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebScraperSpec   `json:"spec,omitempty"`
	Status WebScraperStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WebScraperList contains a list of WebScraper
type WebScraperList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WebScraper `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WebScraper{}, &WebScraperList{})
}
