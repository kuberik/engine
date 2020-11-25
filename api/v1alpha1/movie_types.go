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
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MovieSpec defines the desired state of Movie
type MovieSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// TODO remove this
	Template PlayTemplate `json:"template"`
	// +optional
	FailedJobsHistoryLimit int `json:"failedJobsHistoryLimit"`
	// +optional
	SuccessfulJobsHistoryLimit int `json:"successfulJobsHistoryLimit"`
}

// PlayTemplate defines a template of Play to be created from a Movie
type PlayTemplate struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PlaySpec `json:"spec,omitempty"`
}

// MovieStatus defines the observed state of Movie
type MovieStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Movie is the Schema for the movies API
type Movie struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MovieSpec   `json:"spec,omitempty"`
	Status MovieStatus `json:"status,omitempty"`
}

// TODO remove
func (m *Movie) GeneratePlay(vars ...Var) (*Play, error) {
	play := &Play{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-", m.Name),
			Namespace:    m.Namespace,
		},
		Spec: m.Spec.Template.Spec,
	}

	for _, v := range vars {
		if err := play.Spec.Vars.Set(v.Name, *v.Value); err != nil {
			// There might be more provided vars than play requires
			continue
		}
	}
	for _, v := range vars {
		if v.Value == nil {
			return nil, fmt.Errorf("Missing required variable: %s", v.Name)
		}
	}

	return play, nil
}

// +kubebuilder:object:root=true

// MovieList contains a list of Movie
type MovieList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Movie `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Movie{}, &MovieList{})
}
