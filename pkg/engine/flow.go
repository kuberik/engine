package engine

import (
	"fmt"

	corev1alpha1 "github.com/kuberik/engine/api/v1alpha1"
	"github.com/kuberik/engine/pkg/engine/scheduler"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

const (
	frameCopyIndexVar  = "FRAME_COPY_INDEX"
	mainScreenplayName = "main"
)

// Flow implements ordered exeuction of Actions in a Play
// Sceneres are executed one after another.
// For scene to be completed, all its frames need to be completed
// Frames of a Scene are executed in parallel
type Flow struct {
	scheduler scheduler.Scheduler
}

// NewFlow creates a new Flow that executes actions with given Scheduler
func NewFlow(scheduler scheduler.Scheduler) Flow {
	return Flow{
		scheduler: scheduler,
	}
}

// PlayNext executes all actions that are possible to play at the current stage
func (f *Flow) PlayNext(play *corev1alpha1.Play) error {
	// Expand definition
	populateVars(play, play.Status.VarsConfigMap)
	expandCopies(&play.Spec)
	expandProvisionedVolumes(play)
	return f.playScreenplay(play, mainScreenplayName)
}

func framesFinished(status *corev1alpha1.PlayStatus, frames []corev1alpha1.Frame) bool {
	sceneFinished := true
	for _, frame := range frames {
		_, ok := status.Frames[frame.ID]
		sceneFinished = sceneFinished && ok
	}
	return sceneFinished
}

func (f *Flow) playScreenplay(play *corev1alpha1.Play, name string) error {
	screenplay := play.Screenplay(name)
	if screenplay == nil {
		return fmt.Errorf("Screenplay '%s' not found in the Play", name)
	}

	if !play.Status.Failed() {
		if screenplay.Credits != nil && !framesFinished(&play.Status, screenplay.Credits.Opening) {
			return f.playFrames(play, screenplay.Credits.Opening)
		}

		for si := range screenplay.Scenes {
			if framesFinished(&play.Status, screenplay.Scenes[si].Frames) {
				continue
			}

			return f.playFrames(play, screenplay.Scenes[si].Frames)
		}
	}

	if screenplay.Credits != nil && !framesFinished(&play.Status, screenplay.Credits.Closing) {
		addScreenplayResult(screenplay.Credits.Closing, play, screenplay.Name)
		return f.playFrames(play, screenplay.Credits.Closing)
	}

	return NewError(NoMoreFrames)
}

func (f *Flow) playFrames(play *corev1alpha1.Play, frames []corev1alpha1.Frame) error {
	for _, frame := range frames {
		if _, ok := play.Status.Frames[frame.ID]; ok {
			continue
		}
		err := f.playFrame(play, frame.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Flow) playFrame(play *corev1alpha1.Play, frameID string) error {
	err := f.scheduler.Run(play, frameID)
	if err != nil {
		log.Errorf("Failed to play %s from %s: %s", frameID, play.Name, err)
	}
	return err
}

func expandCopies(playSpec *corev1alpha1.PlaySpec) {
	for k := range playSpec.Screenplays {
		for si := range playSpec.Screenplays[k].Scenes {
			var frames []corev1alpha1.Frame
			for _, f := range playSpec.Screenplays[k].Scenes[si].Frames {
				if f.Copies > 1 {
					for i := 0; i < f.Copies; i++ {
						fc := f.Copy()

						fc.ID = fmt.Sprintf("%s-%v", fc.ID, i)
						fc.Name = fmt.Sprintf("%s-%v", fc.Name, i)
						for ci := range fc.Action.Template.Spec.Containers {
							fc.Action.Template.Spec.Containers[ci].Env = append(fc.Action.Template.Spec.Containers[ci].Env, corev1.EnvVar{
								Name:  frameCopyIndexVar,
								Value: fmt.Sprintf("%v", i),
							})
						}
						frames = append(frames, fc)
					}
				} else {
					frames = append(frames, f)
				}
			}
			playSpec.Screenplays[k].Scenes[si].Frames = frames
		}
	}
}

func expandProvisionedVolumes(play *corev1alpha1.Play) {
	volumes := play.Status.ProvisionedVolumes
	for k := range play.Spec.Screenplays {
		for si := range play.Spec.Screenplays[k].Scenes {
			for fi := range play.Spec.Screenplays[k].Scenes[si].Frames {
			volumes:
				for volumeName, provisionedVolumeName := range volumes {
					for _, container := range play.Spec.Screenplays[k].Scenes[si].Frames[fi].Action.Template.Spec.Containers {
						for _, m := range container.VolumeMounts {
							if m.Name == volumeName {
								play.Spec.Screenplays[k].Scenes[si].Frames[fi].Action.Template.Spec.Volumes = append(
									play.Spec.Screenplays[k].Scenes[si].Frames[fi].Action.Template.Spec.Volumes,
									corev1.Volume{
										Name: volumeName,
										VolumeSource: corev1.VolumeSource{
											PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
												ClaimName: provisionedVolumeName,
											},
										},
									},
								)
								continue volumes
							}
						}
					}
					for _, container := range play.Spec.Screenplays[k].Scenes[si].Frames[fi].Action.Template.Spec.InitContainers {
						for _, m := range container.VolumeMounts {
							if m.Name == volumeName {
								play.Spec.Screenplays[k].Scenes[si].Frames[fi].Action.Template.Spec.Volumes = append(
									play.Spec.Screenplays[k].Scenes[si].Frames[fi].Action.Template.Spec.Volumes,
									corev1.Volume{
										Name: volumeName,
										VolumeSource: corev1.VolumeSource{
											PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
												ClaimName: provisionedVolumeName,
											},
										},
									},
								)
								continue volumes
							}
						}
					}
				}
			}
		}
	}
}

func populateVars(play *corev1alpha1.Play, varsConfigMap string) {
	if varsConfigMap == "" {
		return
	}
	mountName := "kuberik-vars"
	mountPath := "/kuberik/vars"
	frames := play.AllFrames()
	for fi := range frames {
		frames[fi].Action.Template.Spec.Volumes = append(
			frames[fi].Action.Template.Spec.Volumes,
			corev1.Volume{
				Name: mountName,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: varsConfigMap,
						},
					},
				},
			},
		)
		for ci := range frames[fi].Action.Template.Spec.Containers {
			frames[fi].Action.Template.Spec.Containers[ci].EnvFrom = append(
				frames[fi].Action.Template.Spec.Containers[ci].EnvFrom,
				corev1.EnvFromSource{
					ConfigMapRef: &corev1.ConfigMapEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: varsConfigMap,
						},
					},
				},
			)
			frames[fi].Action.Template.Spec.Containers[ci].VolumeMounts = append(
				frames[fi].Action.Template.Spec.Containers[ci].VolumeMounts,
				corev1.VolumeMount{
					Name:      mountName,
					MountPath: mountPath,
				},
			)
		}
		for ci := range frames[fi].Action.Template.Spec.InitContainers {
			frames[fi].Action.Template.Spec.InitContainers[ci].EnvFrom = append(
				frames[fi].Action.Template.Spec.InitContainers[ci].EnvFrom,
				corev1.EnvFromSource{
					ConfigMapRef: &corev1.ConfigMapEnvSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: varsConfigMap,
						},
					},
				},
			)
			frames[fi].Action.Template.Spec.InitContainers[ci].VolumeMounts = append(
				frames[fi].Action.Template.Spec.InitContainers[ci].VolumeMounts,
				corev1.VolumeMount{
					Name:      mountName,
					MountPath: mountPath,
				},
			)
		}
	}
}

const (
	kuberikScreenplayResultEnv         = "KUBERIK_SCREENPLAY_RESULT"
	kuberikScreenplayResultValueSucces = "success"
	kuberikScreenplayResultValueFail   = "fail"
)

func addScreenplayResult(frames []corev1alpha1.Frame, play *corev1alpha1.Play, screenplayName string) {
	var result string
	for _, s := range play.Screenplay(screenplayName).Scenes {
		for _, f := range s.Frames {
			if play.Status.Frames[f.ID] == corev1alpha1.FrameResultFailed {
				result = kuberikScreenplayResultValueFail
			}
		}
	}
	if result == "" {
		result = kuberikScreenplayResultValueSucces
	}
	for fi := range frames {
		mutateContainers := func(containers []corev1.Container) {
			for ci := range containers {
				containers[ci].Env = append(containers[ci].Env, corev1.EnvVar{
					Name:  kuberikScreenplayResultEnv,
					Value: result,
				})
			}
		}
		mutateContainers(frames[fi].Action.Template.Spec.Containers)
		mutateContainers(frames[fi].Action.Template.Spec.InitContainers)
	}
}
