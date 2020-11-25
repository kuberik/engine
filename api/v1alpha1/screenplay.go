package v1alpha1

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
)

// Screenplay describes how pipeline execution will look like.
type Screenplay struct {
	Name    string   `json:"name,omitempty"`
	Scenes  []Scene  `json:"scenes,omitempty"`
	Credits *Credits `json:"credits,omitempty"`
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

// Var is a parametrizable variable for the screenplay shared between all jobs.
// +k8s:openapi-gen=true
type Var struct {
	Name  string  `json:"name"`
	Value *string `json:"value,omitempty"`
}

type Vars []Var

func (vars Vars) Get(name string) (string, error) {
	for _, v := range vars {
		if v.Name == name {
			return *v.Value, nil
		}
	}
	return "", fmt.Errorf("Variable not found")
}

func (vars Vars) Set(name, value string) error {
	for i, v := range vars {
		if v.Name == name {
			vars[i].Value = &value
			return nil
		}
	}
	return fmt.Errorf("Variable not declared")
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
	Name         string    `json:"name"`
	Frames       []Frame   `json:"frames"`
	Pass         Condition `json:"pass,omitempty"`
	IgnoreErrors bool      `json:"ignore_errors,omitempty" yaml:"ignore_errors"`
}

// Condition describes a logical filter which controls execution of the pipeline
type Condition []map[string]string

// Evaluate returns the result of condition filter
func (c Condition) Evaluate(vars Vars) bool {
	var pass bool
	for _, conditions := range c {
		conditionPass := true
		for variable, v := range conditions {
			varValue, err := vars.Get(variable)

			if err != nil {
				conditionPass = conditionPass && false
				// TODO process error
				break
			}

			if varValue != v {
				conditionPass = conditionPass && false
				break
			}
		}
		pass = pass || conditionPass
	}
	return pass
}

// Frame describes either an action or story that needs to be executed
type Frame struct {
	ID            string    `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	IgnoreErrors  bool      `json:"ignoreErrors,omitempty"`
	Copies        int       `json:"copies,omitempty"`
	SkipCondition Condition `json:"skipCondition,omitempty"`
	Action        *Exec     `json:"action,omitempty"`
	Story         *string   `json:"story,omitempty"`
}

// FrameResult represents end result of a frame
type FrameResult int

const (
	// FrameResultSuccessful indicates that frame ended sucessfully
	FrameResultSuccessful FrameResult = iota
	// FrameResultFailed indicates that frame finished with an error
	FrameResultFailed
)

func (fr FrameResult) String() string {
	switch fr {
	case FrameResultSuccessful:
		return "success"
	case FrameResultFailed:
		return "failed"
	}
	panic("Frame result not defined")
}

// Exec Represents a running container
type Exec = batchv1.JobSpec

// Copy makes a copy of the frame
func (f *Frame) Copy() Frame {
	return Frame{
		ID:     f.ID,
		Name:   f.Name,
		Action: f.Action.DeepCopy(),
		Copies: f.Copies,
	}
}