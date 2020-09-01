/*


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

type SecretRef struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
	DataKey   string `json:"key,omitempty"`
}

type Secret struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

// GoogleSecretEntrySpec defines the desired state of GoogleSecretEntry
type GoogleSecretEntrySpec struct {
	SecretRef *SecretRef `json:"secretRef,omitempty"`
	Secrets   []Secret   `json:"secrets"`
}

// GoogleSecretEntryStatus defines the observed state of GoogleSecretEntry
type GoogleSecretEntryStatus struct {
	ServiceAccountEmail string `json:"serviceAccountEmail"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// GoogleSecretEntry is the Schema for the googlesecretentries API
type GoogleSecretEntry struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GoogleSecretEntrySpec   `json:"spec,omitempty"`
	Status GoogleSecretEntryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GoogleSecretEntryList contains a list of GoogleSecretEntry
type GoogleSecretEntryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GoogleSecretEntry `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GoogleSecretEntry{}, &GoogleSecretEntryList{})
}
