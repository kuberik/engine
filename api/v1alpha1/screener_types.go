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
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ScreenerSpec defines the desired state of Screener
type ScreenerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Class string `json:"class"`
	Type  string `json:"type"`
	// Movie is referencing a Movie for which this Screener creates Events
	Movie  string               `json:"movie"`
	Config runtime.RawExtension `json:"config"`
}

// ScreenerStatus defines the observed state of Screener
type ScreenerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Screener is the Schema for the screeners API
type Screener struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScreenerSpec   `json:"spec,omitempty"`
	Status ScreenerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ScreenerList contains a list of Screener
type ScreenerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Screener `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Screener{}, &ScreenerList{})
}
