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

package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/prometheus/common/log"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1alpha1 "github.com/kuberik/engine/api/v1alpha1"
	"github.com/kuberik/engine/pkg/engine"
	"github.com/kuberik/engine/pkg/engine/scheduler"
	"github.com/kuberik/engine/pkg/randutils"

	batchv1 "k8s.io/api/batch/v1"
)

// PlayReconciler reconciles a Play object
type PlayReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme

	Flow engine.Flow
}

// +kubebuilder:rbac:groups=core.kuberik.io,resources=plays,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.kuberik.io,resources=plays/status,verbs=get;update;patch

func (r *PlayReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("play", req.NamespacedName)

	instance := &corev1alpha1.Play{}
	ctx := context.TODO()
	err := r.Client.Get(ctx, request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	switch instance.Status.Phase {
	case "", corev1alpha1.PlayPhaseCreated:
		return r.reconcileCreated(instance)
	case corev1alpha1.PlayPhaseInit:
		return r.reconcileInit(instance)
	case corev1alpha1.PlayPhaseRunning:
		return r.reconcileRunning(instance)
	case corev1alpha1.PlayPhaseComplete, corev1alpha1.PlayPhaseFailed, corev1alpha1.PlayPhaseError:
		return r.reconcileComplete(instance)
	}
	return reconcile.Result{}, nil
}

func (r *PlayReconciler) reconcileCreated(instance *corev1alpha1.Play) (reconcile.Result, error) {
	instance.Status.Phase = corev1alpha1.PlayPhaseInit
	err := r.Client.Status().Update(context.TODO(), instance)
	return reconcile.Result{}, err
}

func (r *PlayReconciler) reconcileInit(instance *corev1alpha1.Play) (reconcile.Result, error) {
	err := func() error {
		err := r.provisionVarsConfigMap(instance)
		if err != nil {
			return err
		}
		err = r.provisionVolumes(instance)
		return err
	}()

	if err != nil {
		instance.Status.Phase = corev1alpha1.PlayPhaseError
		if errUpdate := r.Client.Status().Update(context.TODO(), instance); errUpdate != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, err
	}

	r.populateRandomIDs(instance)
	err = r.Client.Update(context.TODO(), instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	log.Info(fmt.Sprintf("Running play %s", instance.Name))
	instance.Status.Phase = corev1alpha1.PlayPhaseRunning
	err = r.Client.Status().Update(context.TODO(), instance)
	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *PlayReconciler) reconcileRunning(instance *corev1alpha1.Play) (reconcile.Result, error) {
	if err := r.updateStatus(instance); err != nil {
		return reconcile.Result{}, err
	}

	err := r.Flow.PlayNext(instance)
	if engine.IsPlayEndedErorr(err) {
		if instance.Status.Failed() {
			instance.Status.Phase = corev1alpha1.PlayPhaseFailed
		} else {
			instance.Status.Phase = corev1alpha1.PlayPhaseComplete
		}
		return reconcile.Result{}, r.Client.Status().Update(context.TODO(), instance)
	}
	return reconcile.Result{}, err
}

func (r *PlayReconciler) reconcileComplete(instance *corev1alpha1.Play) (reconcile.Result, error) {
	for _, pvcName := range instance.Status.ProvisionedVolumes {
		r.Client.Delete(context.TODO(), &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pvcName,
				Namespace: instance.Namespace,
			},
		})
	}
	instance.Status.ProvisionedVolumes = make(map[string]string)
	err := r.Client.Status().Update(context.TODO(), instance)
	log.Info(fmt.Sprintf("Play %s competed with status: %s", instance.Name, instance.Status.Phase))
	return reconcile.Result{}, err
}

func (r *PlayReconciler) populateRandomIDs(play *corev1alpha1.Play) {
	frames := play.AllFrames()
	randomIDs := randutils.RandList(len(frames))
	for i, f := range frames {
		f.ID = randomIDs[i]
	}
}

func (r *PlayReconciler) provisionVarsConfigMap(instance *corev1alpha1.Play) error {
	varsConfigMapName := fmt.Sprintf("%s-vars", instance.Name)
	configMapValues := make(map[string]string)
	for _, v := range instance.Spec.Vars {
		configMapValues[v.Name] = *v.Value
	}
	varsConfigMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      varsConfigMapName,
			Namespace: instance.Namespace,
		},
		Data: configMapValues,
	}

	err := r.Client.Create(context.TODO(), varsConfigMap)
	if err != nil && !errors.IsAlreadyExists(err) {
		return err
	}
	instance.Status.VarsConfigMap = varsConfigMapName
	return nil
}

func (r *PlayReconciler) updateStatus(play *corev1alpha1.Play) error {
	jobs := &batchv1.JobList{}
	r.Client.List(context.TODO(), jobs, &client.ListOptions{
		LabelSelector: func() labels.Selector {
			ls, _ := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
				MatchLabels: map[string]string{
					scheduler.JobLabelPlay: play.Name,
				},
			})
			return ls
		}(),
	})

	var updated bool
	for _, j := range jobs.Items {
		frameID := j.Annotations[scheduler.JobAnnotationFrameID]
		if _, ok := play.Status.Frames[frameID]; ok {
			continue
		}

		updated = true
		if finished, exit := jobResult(&j); finished {
			play.Status.SetFrameStatus(frameID, exit)
		}
	}

	if !updated {
		return nil
	}

	return r.Client.Status().Update(context.TODO(), play)
}

// ProvisionVolumes provisions volumes for the duration of the play
func (r *PlayReconciler) provisionVolumes(play *corev1alpha1.Play) (err error) {
	if play.Status.ProvisionedVolumes == nil {
		play.Status.ProvisionedVolumes = make(map[string]string)
	}

	for _, volumeClaimTemplate := range play.Spec.VolumeClaimTemplates {
		pvcName := fmt.Sprintf("%s-%s", play.Name, volumeClaimTemplate.Name)

		err = r.Client.Create(context.TODO(), &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      pvcName,
				Namespace: play.Namespace,
				Labels: map[string]string{
					"core.kuberik.io/play": play.Name,
				},
			},
			Spec: volumeClaimTemplate.Spec,
		})

		if err != nil && !errors.IsAlreadyExists(err) {
			return
		}
		play.Status.ProvisionedVolumes[volumeClaimTemplate.Name] = pvcName
	}
	return
}

func jobResult(job *batchv1.Job) (finished bool, exit corev1alpha1.FrameResult) {
	// Successfully completed a single instance of a job
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobFailed || condition.Type == batchv1.JobComplete {
			finished = true
			if condition.Type == batchv1.JobComplete {
				exit = corev1alpha1.FrameResultSuccessful
			} else {
				exit = corev1alpha1.FrameResultFailed
			}
		}
	}
	return
}

func (r *PlayReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.Play{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
