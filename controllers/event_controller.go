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
	"reflect"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1alpha1 "github.com/kuberik/engine/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EventReconciler reconciles a Event object
type EventReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core.kuberik.io,resources=events,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.kuberik.io,resources=events/status,verbs=get;update;patch

func (r *EventReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("event", req.NamespacedName)

	// Fetch the Event instance
	instance := &corev1alpha1.Event{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	movie := &corev1alpha1.Movie{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{
		Name:      instance.Spec.Movie,
		Namespace: instance.Namespace,
	}, movie)
	if err != nil {
		// TODO: update status to error
		// TODO: this should not happen if event validation hook is deployed
		return reconcile.Result{}, err
	}

	// TODO: test the GeneratePlay method
	// TODO: test using operator-sdk e2e testing
	p := generateEventPlay(*movie, *instance)
	r.Client.Create(context.TODO(), &p)
	if err != nil && !errors.IsAlreadyExists(err) {
		return reconcile.Result{Requeue: true}, err
	}
	return reconcile.Result{}, nil
}

func (r *EventReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.Event{}).
		Complete(r)
}

func generatePlay(movie corev1alpha1.Movie) corev1alpha1.Play {
	return corev1alpha1.Play{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: movie.Namespace,
			OwnerReferences: []metav1.OwnerReference{*metav1.NewControllerRef(
				&movie, corev1alpha1.GroupVersion.WithKind(reflect.TypeOf(corev1alpha1.Movie{}).Name()),
			)},
		},
		Spec: movie.Spec.Template.Spec,
	}
}

func generateEventPlay(movie corev1alpha1.Movie, event corev1alpha1.Event) corev1alpha1.Play {
	play := generatePlay(movie)
	play.OwnerReferences = append(play.OwnerReferences, *metav1.NewControllerRef(
		&event, corev1alpha1.GroupVersion.WithKind(reflect.TypeOf(corev1alpha1.Event{}).Name())),
	)
	play.Name = fmt.Sprintf("%s-%s", movie.Name, event.Name)
	eventDataConfigMap := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "event-data",
		},
		Data: event.Spec.Data,
	}
	play.Spec.Screenplays[0].Provision.Resources = append(
		play.Spec.Screenplays[0].Provision.Resources,
		runtime.RawExtension{Object: &eventDataConfigMap},
	)
	return play
}
