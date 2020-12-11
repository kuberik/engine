package v1alpha1

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// Screenplay describes how pipeline execution will look like.
type Screenplay struct {
	Name      string `json:"name,omitempty"`
	Provision `json:"provision,omitempty"`
	Scenes    []Scene  `json:"scenes,omitempty"`
	Credits   *Credits `json:"credits,omitempty"`
}

// Credits describe actions that need to be run at the start or at the end of a screenplay.
// Actions are ran in parallel
type Credits struct {
	// Opening credits are played before anything else in the scene.
	Opening []Frame `json:"opening,omitempty"`

	// Closing credits are played after screenplay is finished.
	// Finished in this case means started and ended with any result.
	// This provides a way to run some tasks even if some frames failed.
	Closing []Frame `json:"closing,omitempty"`
}

type Provision struct {
	Resources []runtime.RawExtension `json:"resources,omitempty"`
}

// Scene finds a scene by name
func (s *Screenplay) Scene(name string) (*Scene, error) {
	for _, a := range s.Scenes {
		if a.Name == name {
			return &a, nil
		}
	}
	return &Scene{}, fmt.Errorf("Scene not found")
}

// Scene describes a collection of frames that need to be executed in parallel
type Scene struct {
	Name   string  `json:"name"`
	Frames []Frame `json:"frames"`
}

// Frame describes either an action or story that needs to be executed
type Frame struct {
	ID     string  `json:"id,omitempty"`
	Name   string  `json:"name,omitempty"`
	Copies int     `json:"copies,omitempty"`
	Action *Action `json:"action,omitempty"`
	Story  *string `json:"story,omitempty"`
}

// FrameStatus represents end result of a frame
type FrameStatus int

const (
	// FrameStatusSuccessful indicates that frame ended sucessfully
	FrameStatusSuccessful FrameStatus = iota
	// FrameStatusFailed indicates that frame finished with an error
	FrameStatusFailed
	// FrameStatusRunning indicates that frame is running
	FrameStatusRunning
)

func (fr FrameStatus) String() string {
	switch fr {
	case FrameStatusSuccessful:
		return "success"
	case FrameStatusFailed:
		return "failed"
	case FrameStatusRunning:
		return "running"
	}
	panic("Frame result not defined")
}

// +kubebuilder:object:generate=false

// Action represents an action that needs to be executed
type Action = batchv1.JobSpec

// func (a *Action) FrameID() {

// }
