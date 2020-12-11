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

// PlaySpec defines the desired state of Play
type PlaySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Screenplays []Screenplay `json:"screenplays"`
}

// PlayStatus defines the observed state of Play
type PlayStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Frames map[string]FrameStatus `json:"frames,omitempty"`
	Phase  PlayPhaseType          `json:"phase,omitempty"`
}

// SetFrameStatus sets result of a frame
func (ps *PlayStatus) SetFrameStatus(frameID string, result FrameStatus) {
	if ps.Frames == nil {
		ps.Frames = make(map[string]FrameStatus)
	}
	ps.Frames[frameID] = result
}

// Failed checks if a play failed
func (ps *PlayStatus) Failed() bool {
	if ps.Frames == nil {
		return false
	}
	for _, r := range ps.Frames {
		if r == FrameStatusFailed {
			return true
		}
	}
	return false
}

// PlayPhaseType defines the phase of a Play
type PlayPhaseType string

// These are valid phases of a play.
const (
	// PlayPhaseComplete means the play has completed its execution.
	PlayPhaseComplete PlayPhaseType = "Complete"
	// PlayPhaseInit means the play is in initializing phase
	PlayPhaseInit PlayPhaseType = "Init"
	// PlayPhaseFailed means the play has failed its execution.
	PlayPhaseFailed PlayPhaseType = "Failed"
	// PlayPhaseRunning means the play is executing.
	PlayPhaseRunning PlayPhaseType = "Running"
	// PlayPhaseRunning means the play has been created.
	PlayPhaseCreated PlayPhaseType = "Created"
	// PlayPhaseError means the play ended because of an error.
	PlayPhaseError PlayPhaseType = "Error"
)

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Play is the Schema for the plays API
type Play struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PlaySpec   `json:"spec,omitempty"`
	Status PlayStatus `json:"status,omitempty"`
}

// Frame gets a frame with specified identifier
func (p *Play) Frame(frameID string) *Frame {
	for spi, screenplay := range p.Spec.Screenplays {
		for sci, scene := range screenplay.Scenes {
			for fi, frame := range scene.Frames {
				if frame.ID == frameID {
					return &p.Spec.Screenplays[spi].Scenes[sci].Frames[fi]
				}
			}
		}
		for fi, frame := range screenplay.Credits.Opening {
			if frame.ID == frameID {
				return &p.Spec.Screenplays[spi].Credits.Opening[fi]
			}
		}
		for fi, frame := range screenplay.Credits.Closing {
			if frame.ID == frameID {
				return &p.Spec.Screenplays[spi].Credits.Closing[fi]
			}
		}
	}
	return nil
}

// Screenplay gets a Screenplay with specified name
func (p *Play) Screenplay(name string) *Screenplay {
	for spi, screenplay := range p.Spec.Screenplays {
		if screenplay.Name == name {
			return &p.Spec.Screenplays[spi]
		}
	}
	return nil
}

// AllFrames returns references to all Frames in the Play
func (p *Play) AllFrames() (frames []*Frame) {
	playSpec := &p.Spec
	for k := range playSpec.Screenplays {
		for i := range playSpec.Screenplays[k].Scenes {
			for j := range playSpec.Screenplays[k].Scenes[i].Frames {
				frames = append(frames, &playSpec.Screenplays[k].Scenes[i].Frames[j])
			}
		}
		if playSpec.Screenplays[k].Credits != nil {
			for i := range playSpec.Screenplays[k].Credits.Opening {
				frames = append(frames, &playSpec.Screenplays[k].Credits.Opening[i])
			}
			for i := range playSpec.Screenplays[k].Credits.Closing {
				frames = append(frames, &playSpec.Screenplays[k].Credits.Closing[i])
			}
		}
	}
	return
}

// +kubebuilder:object:root=true

// PlayList contains a list of Play
type PlayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Play `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Play{}, &PlayList{})
}
