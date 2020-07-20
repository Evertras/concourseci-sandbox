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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SmeeProxySpec defines the desired state of SmeeProxy
type SmeeProxySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// SmeeURL points to a smee url like https://smee.io/abcdef
	SmeeURL string `json:"smeeURL"`

	// TargetURL is where to point the webhook events to, such as http://internal.svc.cluster/hook
	TargetURL string `json:"targetURL"`
}

// SmeeProxyStatus defines the observed state of SmeeProxy
type SmeeProxyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Active is the current pod managing this proxy
	Active string `json:"active"`
}

// +kubebuilder:object:root=true

// SmeeProxy is the Schema for the smeeproxies API
type SmeeProxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SmeeProxySpec   `json:"spec,omitempty"`
	Status SmeeProxyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SmeeProxyList contains a list of SmeeProxy
type SmeeProxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SmeeProxy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SmeeProxy{}, &SmeeProxyList{})
}
