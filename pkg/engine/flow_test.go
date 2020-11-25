package engine

import (
	"testing"

	corev1alpha1 "github.com/kuberik/engine/pkg/apis/core/v1alpha1"
	"github.com/kuberik/engine/pkg/engine/scheduler"
	corev1 "k8s.io/api/core/v1"
)

var (
	failed  = corev1alpha1.FrameResultFailed
	success = corev1alpha1.FrameResultSuccessful
)

func helloWorldAction() *corev1alpha1.Exec {
	return &corev1alpha1.Exec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{{
					Name:    "hello",
					Image:   "alpine",
					Command: []string{"echo", "Hello world!"},
				}},
			},
		},
	}
}

func assertFrameState(t *testing.T, play *corev1alpha1.Play, states map[string]*corev1alpha1.FrameResult) {
	for k, v := range states {
		if v == nil {
			if _, ok := play.Status.Frames[k]; ok {
				t.Errorf("Excpected %s to not be played yet", k)
			}
		} else {
			if r, ok := play.Status.Frames[k]; !ok || r != *v {
				t.Errorf("Excpected %s to played and finished with status '%s', but got '%s'", k, v, r)
			}
		}
	}
}

func TestExpandLoops(t *testing.T) {
	screenplay := corev1alpha1.Screenplay{
		Scenes: []corev1alpha1.Scene{{
			Frames: []corev1alpha1.Frame{{
				Copies: 3,
				Action: &corev1alpha1.Exec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{},
							},
						},
					},
				},
			},
			},
		},
		},
	}
	expandCopies(&corev1alpha1.PlaySpec{
		Screenplays: []corev1alpha1.Screenplay{
			screenplay,
		},
	})

	if len(screenplay.Scenes[0].Frames) != 3 {
		t.Errorf("Expand loop doesn't add new frames")
	}
	found := false
	for _, e := range screenplay.Scenes[0].Frames[1].Action.Template.Spec.Containers[0].Env {
		if e.Name == frameCopyIndexVar {
			found = true
			if e.Value != "1" {
				t.Errorf("Index variable is not correctly populated")
			}
		}
	}
	if !found {
		t.Errorf("Index variable is not injected")
	}
}

func TestPlayNextLoop(t *testing.T) {
	play := &corev1alpha1.Play{
		Spec: corev1alpha1.PlaySpec{
			Screenplays: []corev1alpha1.Screenplay{{
				Name: "main",
				Scenes: []corev1alpha1.Scene{
					{
						Name: "first-scene",
						Frames: []corev1alpha1.Frame{
							{
								ID:     "a",
								Name:   "first-hello-a",
								Action: helloWorldAction(),
							},
							{
								ID:     "b",
								Name:   "first-hello-b",
								Action: helloWorldAction(),
							},
						},
					},
					{
						Name: "second-scene",
						Frames: []corev1alpha1.Frame{
							{
								ID:     "c",
								Name:   "second-hello-a",
								Action: helloWorldAction(),
							},
							{
								ID:     "d",
								Name:   "second-hello-b",
								Action: helloWorldAction(),
							},
						},
					},
				},
			}},
		},
	}

	flow := NewFlow(&scheduler.DummyScheduler{})
	flow.PlayNext(play)
	// Mark "a" as not played
	delete(play.Status.Frames, "a")
	assertFrameState(t, play, map[string]*corev1alpha1.FrameResult{
		"a": nil,
		"b": &success,
		"c": nil,
		"d": nil,
	})

	flow.PlayNext(play)
	assertFrameState(t, play, map[string]*corev1alpha1.FrameResult{
		"a": &success,
		"b": &success,
		"c": nil,
		"d": nil,
	})

	flow.PlayNext(play)
	assertFrameState(t, play, map[string]*corev1alpha1.FrameResult{
		"a": &success,
		"b": &success,
		"c": &success,
		"d": &success,
	})
}

func TestPlayNextWithCredits(t *testing.T) {
	play := &corev1alpha1.Play{
		Spec: corev1alpha1.PlaySpec{
			Screenplays: []corev1alpha1.Screenplay{{
				Name: "main",
				Credits: &corev1alpha1.Credits{
					Opening: []corev1alpha1.Frame{{
						Name:   "opening",
						ID:     "a",
						Action: helloWorldAction(),
					}},
					Closing: []corev1alpha1.Frame{{
						Name:   "closing",
						ID:     "d",
						Action: helloWorldAction(),
					}},
				},
				Scenes: []corev1alpha1.Scene{
					{
						Name: "first-scene",
						Frames: []corev1alpha1.Frame{
							{
								ID:     "b",
								Name:   "first-hello-a",
								Action: helloWorldAction(),
							},
						},
					},
					{
						Name: "second-scene",
						Frames: []corev1alpha1.Frame{
							{
								ID:     "c",
								Name:   "second-hello-a",
								Action: helloWorldAction(),
							},
						},
					},
				},
			}},
		},
	}

	flow := NewFlow(&scheduler.DummyScheduler{})
	flow.PlayNext(play)
	assertFrameState(t, play, map[string]*corev1alpha1.FrameResult{
		"a": &success,
		"b": nil,
		"c": nil,
		"d": nil,
	})

	flow.PlayNext(play)
	assertFrameState(t, play, map[string]*corev1alpha1.FrameResult{
		"a": &success,
		"b": &success,
		"c": nil,
		"d": nil,
	})

	flow.PlayNext(play)
	assertFrameState(t, play, map[string]*corev1alpha1.FrameResult{
		"a": &success,
		"b": &success,
		"c": &success,
		"d": nil,
	})

	flow.PlayNext(play)
	assertFrameState(t, play, map[string]*corev1alpha1.FrameResult{
		"a": &success,
		"b": &success,
		"c": &success,
		"d": &success,
	})
}

func TestPlayNextFailedPlay(t *testing.T) {
	play := &corev1alpha1.Play{
		Spec: corev1alpha1.PlaySpec{
			Screenplays: []corev1alpha1.Screenplay{{
				Name: "main",
				Scenes: []corev1alpha1.Scene{
					{
						Name: "first-scene",
						Frames: []corev1alpha1.Frame{
							{
								ID:     "a",
								Name:   "first-hello-a",
								Action: helloWorldAction(),
							},
							{
								ID:     "b",
								Name:   "first-hello-b",
								Action: helloWorldAction(),
							},
						},
					},
					{
						Name: "second-scene",
						Frames: []corev1alpha1.Frame{
							{
								ID:     "c",
								Name:   "second-hello-a",
								Action: helloWorldAction(),
							},
							{
								ID:     "d",
								Name:   "second-hello-b",
								Action: helloWorldAction(),
							},
						},
					},
				},
			}},
		},
	}

	flow := NewFlow(&scheduler.DummyScheduler{})
	flow.PlayNext(play)
	// Mark "a" as not played
	delete(play.Status.Frames, "a")
	assertFrameState(t, play, map[string]*corev1alpha1.FrameResult{
		"a": nil,
		"b": &success,
		"c": nil,
		"d": nil,
	})

	flow = NewFlow(&scheduler.DummyScheduler{Result: 1})
	flow.PlayNext(play)
	assertFrameState(t, play, map[string]*corev1alpha1.FrameResult{
		"a": &failed,
		"b": &success,
		"c": nil,
		"d": nil,
	})

	err := flow.PlayNext(play)
	if !IsPlayEndedErorr(err) {
		t.Errorf("Play should have ended")
	}
	assertFrameState(t, play, map[string]*corev1alpha1.FrameResult{
		"a": &failed,
		"b": &success,
		"c": nil,
		"d": nil,
	})
}

func TestAddScreenplayResult(t *testing.T) {
	play := &corev1alpha1.Play{
		Spec: corev1alpha1.PlaySpec{
			Screenplays: []corev1alpha1.Screenplay{{
				Name: "main",
				Scenes: []corev1alpha1.Scene{{
					Name:   "second-scene",
					Frames: []corev1alpha1.Frame{{ID: "a"}, {ID: "b"}},
				}, {
					Name:   "second-scene",
					Frames: []corev1alpha1.Frame{{ID: "c"}},
				}},
			}, {
				Name: "second",
				Scenes: []corev1alpha1.Scene{{
					Name:   "second-scene",
					Frames: []corev1alpha1.Frame{{ID: "d"}, {ID: "e"}},
				}},
			}},
		},
	}

	var closingFrames []corev1alpha1.Frame
	play.Status = corev1alpha1.PlayStatus{Frames: map[string]corev1alpha1.FrameResult{
		"a": corev1alpha1.FrameResultSuccessful,
		"b": corev1alpha1.FrameResultFailed,
		"d": corev1alpha1.FrameResultFailed,
		"e": corev1alpha1.FrameResultFailed,
	}}
	closingFrames = []corev1alpha1.Frame{{
		Action: helloWorldAction(),
	}}
	addScreenplayResult(closingFrames, play, "main")
	if env := closingFrames[0].Action.Template.Spec.Containers[0].Env[0]; env.Name != kuberikScreenplayResultEnv || env.Value != kuberikScreenplayResultValueFail {
		t.Errorf("Expected to find %s env with status %s", kuberikScreenplayResultEnv, kuberikScreenplayResultValueFail)
	}

	play.Status = corev1alpha1.PlayStatus{Frames: map[string]corev1alpha1.FrameResult{
		"a": corev1alpha1.FrameResultSuccessful,
		"b": corev1alpha1.FrameResultSuccessful,
		"c": corev1alpha1.FrameResultSuccessful,
		"d": corev1alpha1.FrameResultFailed,
		"e": corev1alpha1.FrameResultFailed,
	}}
	closingFrames = []corev1alpha1.Frame{{
		Action: helloWorldAction(),
	}}
	addScreenplayResult(closingFrames, play, "main")
	if env := closingFrames[0].Action.Template.Spec.Containers[0].Env[0]; env.Name != kuberikScreenplayResultEnv || env.Value != kuberikScreenplayResultValueSucces {
		t.Errorf("Expected to find %s env with status %s", kuberikScreenplayResultEnv, kuberikScreenplayResultValueSucces)
	}
}
