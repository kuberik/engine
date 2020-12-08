package controllers

import (
	"context"
	"fmt"
	"testing"

	corev1alpha1 "github.com/kuberik/engine/api/v1alpha1"
	"github.com/kuberik/engine/pkg/engine"
	"github.com/kuberik/engine/pkg/engine/scheduler/k8s"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	playClient    client.Client
	reconcilePlay *PlayReconciler
)

func init() {
	play := &corev1alpha1.Play{}
	scheme := scheme.Scheme
	scheme.AddKnownTypes(corev1alpha1.GroupVersion, play)
	playClient = fake.NewFakeClientWithScheme(scheme)
	reconcilePlay = &PlayReconciler{
		Client: playClient,
		Scheme: scheme,
		Log:    ctrl.Log.WithName("controllers").WithName("Play"),
		Flow:   engine.NewFlow(k8s.NewKubernetesScheduler(playClient)),
	}
}

func TestPlayCreate(t *testing.T) {
	var (
		name      = "hello-world-run"
		namespace = "default"
	)
	play := &corev1alpha1.Play{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1alpha1.PlaySpec{
			Screenplays: []corev1alpha1.Screenplay{{
				Name: "main",
				Scenes: []corev1alpha1.Scene{{
					Frames: []corev1alpha1.Frame{{
						Name: "test",
					}},
				}},
			}},
		},
	}
	playClient.Create(context.TODO(), play)

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	nn := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	req := reconcile.Request{
		NamespacedName: nn,
	}
	_, err := reconcilePlay.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	playClient.Get(context.TODO(), nn, play)
	if play.Status.Phase != corev1alpha1.PlayPhaseInit {
		t.Error("Reconcile create play didn't reach expected phase")
	}
}

func TestPlayInit(t *testing.T) {
	var (
		name      = "hello-world-init"
		namespace = "default"
	)
	play := &corev1alpha1.Play{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1alpha1.PlaySpec{
			Screenplays: []corev1alpha1.Screenplay{{
				Name: "main",
				Scenes: []corev1alpha1.Scene{{
					Frames: []corev1alpha1.Frame{{
						Name: "test",
						Action: &corev1alpha1.Action{
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{{
										Name:    "test",
										Command: []string{"echo", "test"},
										Image:   "alpine",
									}},
								},
							},
						},
					}},
				}},
			}},
		},
		Status: corev1alpha1.PlayStatus{
			Phase: corev1alpha1.PlayPhaseRunning,
		},
	}
	playClient.Create(context.TODO(), play)

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	nn := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	req := reconcile.Request{
		NamespacedName: nn,
	}
	_, err := reconcilePlay.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	playClient.Get(context.TODO(), nn, play)
	if play.Status.Phase != corev1alpha1.PlayPhaseRunning {
		t.Error("Reconcile init play didn't reach expected phase")
	}
}

func TestPlayRunning(t *testing.T) {
	var (
		name      = "hello-world-running"
		namespace = "default"
	)
	play := &corev1alpha1.Play{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1alpha1.PlaySpec{
			Screenplays: []corev1alpha1.Screenplay{{
				Name: "main",
				Scenes: []corev1alpha1.Scene{{
					Name: "test",
					Frames: []corev1alpha1.Frame{{
						Name: "test",
						Action: &corev1alpha1.Action{
							Template: corev1.PodTemplateSpec{
								Spec: corev1.PodSpec{
									Containers: []corev1.Container{{
										Name:    "test",
										Command: []string{"echo", "test"},
										Image:   "alpine",
									}},
								},
							},
						},
					}},
				}},
			}},
		},
		Status: corev1alpha1.PlayStatus{
			Phase: corev1alpha1.PlayPhaseRunning,
		},
	}
	playClient.Create(context.TODO(), play)

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource.
	nn := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	req := reconcile.Request{
		NamespacedName: nn,
	}
	_, err := reconcilePlay.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	playClient.Get(context.TODO(), nn, play)
	if play.Status.Phase != corev1alpha1.PlayPhaseRunning {
		t.Error("Initialize play didn't reach expected phase")
	}

	job := &batchv1.Job{}
	err = playClient.Get(context.TODO(), types.NamespacedName{
		Name:      fmt.Sprintf("%s-%s", play.Screenplay("main").Scenes[0].Frames[0].Name, play.Name),
		Namespace: play.Namespace,
	}, job)
	if err != nil {
		t.Errorf("Failed to find a job created by the Play: %s", err)
	}

	job.Status.Conditions = append(job.Status.Conditions, batchv1.JobCondition{
		Type: batchv1.JobComplete,
	})
	playClient.Status().Update(context.TODO(), job)
	_, err = reconcilePlay.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	playClient.Get(context.TODO(), nn, play)
	if play.Status.Phase != corev1alpha1.PlayPhaseComplete {
		t.Errorf("Play state want %s, got %s", corev1alpha1.PlayPhaseComplete, play.Status.Phase)
	}
}

func TestGetAllFramesWithCredits(t *testing.T) {
	frames := []corev1alpha1.Frame{{
		Name: "a",
	}, {
		Name: "b",
	}, {
		Name: "c",
	}, {
		Name: "d",
	}, {
		Name: "e",
	}, {
		Name: "f",
	}}

	play := &corev1alpha1.Play{Spec: corev1alpha1.PlaySpec{
		Screenplays: []corev1alpha1.Screenplay{{
			Name: "main",
			Credits: &corev1alpha1.Credits{
				Opening: frames[0:1],
				Closing: frames[1:2],
			},
			Scenes: []corev1alpha1.Scene{{
				Frames: frames[2:4],
			}, {
				Frames: frames[4:6],
			}},
		}},
	}}

	allFrames := play.AllFrames()
	if l := len(allFrames); l != 6 {
		t.Errorf("Expected to retrieve %v frames, got %v", len(frames), l)
	}
}

func TestGetAllFrames(t *testing.T) {
	frames := []corev1alpha1.Frame{{
		Name: "a",
	}, {
		Name: "b",
	}, {
		Name: "c",
	}, {
		Name: "d",
	}}

	play := &corev1alpha1.Play{Spec: corev1alpha1.PlaySpec{
		Screenplays: []corev1alpha1.Screenplay{{
			Name: "main",
			Scenes: []corev1alpha1.Scene{{
				Frames: frames[0:2],
			}, {
				Frames: frames[2:4],
			}},
		}},
	}}

	allFrames := play.AllFrames()
	if l := len(allFrames); l != 4 {
		t.Errorf("Expected to retrieve %v frames, got %v", len(frames), l)
	}
}

var _ = Describe("Play controller", func() {
	const (
		PlayName      = "test"
		PlayNamespace = "default"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	provisionedCM := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "foo",
		},
	}

	Context("When creating Play with jobs referencing provisioned objects", func() {
		It("Should create jobs referencing the names of the provisioned objects with their suffixes", func() {
			By("By creating a new Play")
			ctx := context.Background()
			play := &corev1alpha1.Play{
				TypeMeta: metav1.TypeMeta{
					APIVersion: corev1alpha1.GroupVersion.String(),
					Kind:       "Play",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      PlayName,
					Namespace: PlayNamespace,
				},
				Spec: corev1alpha1.PlaySpec{
					Screenplays: []corev1alpha1.Screenplay{{
						Name: "main",
						Provision: corev1alpha1.Provision{
							Resources: []runtime.RawExtension{{Object: &provisionedCM}},
						},
						Scenes: []corev1alpha1.Scene{{
							Frames: []corev1alpha1.Frame{{
								Action: &corev1alpha1.Action{Template: corev1.PodTemplateSpec{Spec: corev1.PodSpec{Containers: []corev1.Container{{}}}}},
							}},
						}},
					}},
				},
			}
			Expect(k8sClient.Create(ctx, play)).Should(Succeed())

			playLookupKey := types.NamespacedName{Name: play.Name, Namespace: play.Namespace}
			createdPlay := &corev1alpha1.Play{}

			Eventually(func() bool {
				err := k8sClient.Get(ctx, playLookupKey, createdPlay)
				if err != nil {
					return false
				}
				return true
			}, timeout, interval).Should(BeTrue())
		})
	})

})
