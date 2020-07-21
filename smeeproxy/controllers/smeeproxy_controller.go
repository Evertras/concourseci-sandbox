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

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	smeeproxyv1 "github.com/evertras/concourseci-sandbox/smeeproxy/api/v1"
)

// SmeeProxyReconciler reconciles a SmeeProxy object
type SmeeProxyReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=smeeproxy.evertras.com,resources=smeeproxies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=smeeproxy.evertras.com,resources=smeeproxies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch,resources=pods,verbs=get;list;watch;create;update;patch;delete

func (r *SmeeProxyReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("smeeproxy", req.NamespacedName)

	podName := "smeeproxy-" + req.Name
	log = log.WithValues("podName", podName)

	var smeeProxy smeeproxyv1.SmeeProxy

	if err := r.Get(ctx, req.NamespacedName, &smeeProxy); err != nil {
		if errors.IsNotFound(err) {
			log.Info("SmeeProxy deleted, deleting any pods it may have had")
			pod := corev1.Pod{}
			if err := r.Client.Get(ctx, client.ObjectKey{Namespace: smeeProxy.Namespace, Name: podName}, &pod); err == nil {
				log.Info("Actually deleting")
				return ctrl.Result{}, r.Client.Delete(ctx, &pod)
			}

			return ctrl.Result{}, nil
		}

		log.Error(err, "SmeeProxy not found")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Looking for existing Pod")

	pod := corev1.Pod{}
	if err := r.Client.Get(ctx, client.ObjectKey{Namespace: smeeProxy.Namespace, Name: podName}, &pod); apierrors.IsNotFound(err) {
		log.Info("Not found, creating")

		pod := buildPod(&smeeProxy, podName)

		if err := r.Client.Create(ctx, pod); err != nil {
			return ctrl.Result{}, err
		}

		log.Info("Created pod")

		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

func buildPod(smeeProxy *smeeproxyv1.SmeeProxy, podName string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: smeeProxy.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(smeeProxy, smeeproxyv1.GroupVersion.WithKind("SmeeProxy")),
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				corev1.Container{
					Name:  "smee-client",
					Image: "deltaprojects/smee-client",
					Args: []string{
						"--url",
						smeeProxy.Spec.SmeeURL,
						"--target",
						smeeProxy.Spec.TargetURL,
					},
				},
			},
		},
	}
}

func (r *SmeeProxyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&smeeproxyv1.SmeeProxy{}).
		Complete(r)
}
