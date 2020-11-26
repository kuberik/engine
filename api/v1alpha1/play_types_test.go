package v1alpha1

import "testing"

func TestGetCreditsFrame(t *testing.T) {
	play := Play{
		Spec: PlaySpec{
			Screenplays: []Screenplay{{
				Credits: &Credits{
					Opening: []Frame{{
						Name: "opening",
						ID:   "a",
					}},
					Closing: []Frame{{
						Name: "closing",
						ID:   "b",
					}},
				},
			}},
		},
	}
	if f := play.Frame("a"); f == nil || f.Name != "opening" {
		t.Errorf("Failed to retrieve opening frame")
	}
	if f := play.Frame("b"); f == nil || f.Name != "closing" {
		t.Errorf("Failed to retrieve closing frame")
	}
}

func TestGetFrame(t *testing.T) {
	play := Play{
		Spec: PlaySpec{
			Screenplays: []Screenplay{{
				Scenes: []Scene{{
					Frames: []Frame{{
						Name: "first",
						ID:   "a",
					}},
				}, {
					Frames: []Frame{{
						Name: "second",
						ID:   "b",
					}},
				}},
			}},
		},
	}
	if f := play.Frame("a"); f == nil || f.Name != "first" {
		t.Errorf("Failed to retrieve opening frame")
	}
	if f := play.Frame("b"); f == nil || f.Name != "second" {
		t.Errorf("Failed to retrieve closing frame")
	}
}

func TestPlayStatusFailed(t *testing.T) {
	var status PlayStatus

	status = PlayStatus{
		Frames: map[string]FrameStatus{},
	}
	if status.Failed() {
		t.Errorf("Status is failed but it's shouldn't be")
	}

	status = PlayStatus{
		Frames: map[string]FrameStatus{
			"a": FrameStatusSuccessful,
			"b": FrameStatusSuccessful,
		},
	}
	if status.Failed() {
		t.Errorf("Status is failed but it's shouldn't be")
	}

	status = PlayStatus{
		Frames: map[string]FrameStatus{
			"a": FrameStatusSuccessful,
			"b": FrameStatusSuccessful,
			"c": FrameStatusFailed,
		},
	}
	if !status.Failed() {
		t.Errorf("Status is not failed but should be")
	}
}
