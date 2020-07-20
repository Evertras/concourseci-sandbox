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
	"crypto/md5"
	"fmt"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
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

func genPodName(smeeURL, targetURL string) string {
	data := []byte(smeeURL + targetURL)
	return "smeeproxy-" + fmt.Sprintf("%x", md5.Sum(data))
}

// +kubebuilder:rbac:groups=smeeproxy.evertras.com,resources=smeeproxies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=smeeproxy.evertras.com,resources=smeeproxies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=batch,resources=pods,verbs=get;list;watch;create;update;patch;delete

func (r *SmeeProxyReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("smeeproxy", req.NamespacedName)

	var smeeProxy smeeproxyv1.SmeeProxy

	if err := r.Get(ctx, req.NamespacedName, &smeeProxy); err != nil {
		log.Error(err, "Unable to fetch SmeeProxy")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	podName := genPodName(smeeProxy.Spec.SmeeURL, smeeProxy.Spec.TargetURL)
	log = log.WithValues("podName", podName)

	log.Info("Looking for existing Pod")

	pod := corev1.Pod{}
	if err := r.Client.Get(ctx, client.ObjectKey{Namespace: smeeProxy.Namespace, Name: podName}, &pod); err != nil {
		log.Info("Nope")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *SmeeProxyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&smeeproxyv1.SmeeProxy{}).
		Complete(r)
}
