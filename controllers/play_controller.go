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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1alpha1 "github.com/kuberik/engine/api/v1alpha1"
	"github.com/kuberik/engine/pkg/engine"
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
	_ = r.Log.WithValues("play", req.NamespacedName)

	instance := &corev1alpha1.Play{}
	ctx := context.TODO()
	err := r.Client.Get(ctx, req.NamespacedName, instance)
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
	r.populateRandomIDs(instance)
	err := r.Client.Update(context.TODO(), instance)
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

	err := r.Flow.Next(instance)
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
	err := r.Flow.Next(instance)
	if err != nil && !engine.IsPlayEndedErorr(err) {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, nil
}

func (r *PlayReconciler) populateRandomIDs(play *corev1alpha1.Play) {
	frames := play.AllFrames()
	randomIDs := randutils.RandList(len(frames))
	for i, f := range frames {
		f.ID = randomIDs[i]
	}
}

func (r *PlayReconciler) updateStatus(play *corev1alpha1.Play) error {
	jobs := &batchv1.JobList{}
	r.Client.List(context.TODO(), jobs, &client.ListOptions{
		LabelSelector: engine.JobLabelSelector(play),
	})

	var updated bool
	for _, j := range jobs.Items {
		frameID := j.Annotations[engine.ActionAnnotationFrameID]
		if _, ok := play.Status.Frames[frameID]; ok {
			continue
		}

		updated = true
		if status := frameStatus(&j); status != corev1alpha1.FrameStatusRunning {
			play.Status.SetFrameStatus(frameID, status)
		}
	}

	if updated {
		return r.Client.Status().Update(context.TODO(), play)
	}

	return nil
}

func frameStatus(job *batchv1.Job) corev1alpha1.FrameStatus {
	// Successfully completed a single instance of a job
	for _, condition := range job.Status.Conditions {
		if condition.Type == batchv1.JobFailed || condition.Type == batchv1.JobComplete {
			if condition.Type == batchv1.JobComplete {
				return corev1alpha1.FrameStatusSuccessful
			} else {
				return corev1alpha1.FrameStatusFailed
			}
		}
	}
	return corev1alpha1.FrameStatusRunning
}

func (r *PlayReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.Play{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}
