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
	"math/rand"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	core "k8s.io/api/core/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	smeeproxyv1 "github.com/evertras/concourseci-sandbox/smeeproxy/api/v1"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	logf.SetLogger(zap.LoggerTo(GinkgoWriter, true))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
	}

	var err error
	cfg, err = testEnv.Start()
	Expect(err).ToNot(HaveOccurred())
	Expect(cfg).ToNot(BeNil())

	err = smeeproxyv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	// +kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).ToNot(HaveOccurred())
	Expect(k8sClient).ToNot(BeNil())

	close(done)
}, 60)

var _ = Context("Inside of a new namespace", func() {
	ctx := context.Background()
	ns := SetupTest(ctx)

	Describe("when no existing pods exist", func() {
		It("creates a new pod with the Smee client with the correct arguments", func() {
			smeeURL := "http://doesnt-exist-smee/abcdef"
			targetURL := "http://localhost:3808/xyz"
			testName := "test-proxy"
			smeeProxy := &smeeproxyv1.SmeeProxy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      testName,
					Namespace: ns.Name,
				},
				Spec: smeeproxyv1.SmeeProxySpec{
					SmeeURL:   smeeURL,
					TargetURL: targetURL,
				},
			}

			err := k8sClient.Create(ctx, smeeProxy)

			Expect(err).NotTo(HaveOccurred(), "failed to create test SmeeProxy resource")

			pod := &corev1.Pod{}

			Eventually(
				func() error {
					return k8sClient.Get(
						ctx,
						client.ObjectKey{
							// This is fragile, how to do this more cleanly with this setup?
							Name:      "smeeproxy-" + testName,
							Namespace: ns.Name,
						},
						pod,
					)
				},
				time.Second*5,
				time.Millisecond*500,
			).Should(BeNil())

			args := pod.Spec.Containers[0].Args
			foundURL := false
			foundTarget := false
			for i := 0; i < len(args)-1; i++ {
				if args[i] == "-u" || args[i] == "--url" {
					foundURL = true
					Expect(args[i+1]).To(Equal(smeeURL))
				} else if args[i] == "-t" || args[i] == "--target" {
					foundTarget = true
					Expect(args[i+1]).To(Equal(targetURL))
				}
			}

			Expect(foundURL).To(BeTrue())
			Expect(foundTarget).To(BeTrue())
		})

		It("deletes a pod when the SmeeProxy resource is deleted", func() {
			testName := "test-proxy"
			smeeProxy := &smeeproxyv1.SmeeProxy{
				ObjectMeta: metav1.ObjectMeta{
					Name:      testName,
					Namespace: ns.Name,
				},
				Spec: smeeproxyv1.SmeeProxySpec{
					SmeeURL:   "http://doesntexistsmee",
					TargetURL: "http://targetalsodoesntexist",
				},
			}

			err := k8sClient.Create(ctx, smeeProxy)

			Expect(err).NotTo(HaveOccurred(), "failed to create test SmeeProxy resource")

			pod := &corev1.Pod{}

			// This is fragile, how to do this more cleanly with this setup?
			podName := "smeeproxy-" + testName

			Eventually(
				func() error {
					return k8sClient.Get(
						ctx,
						client.ObjectKey{
							Name:      podName,
							Namespace: ns.Name,
						},
						pod,
					)
				},
				time.Second*5,
				time.Millisecond*500,
			).Should(BeNil())

			k8sClient.Delete(ctx, smeeProxy)

			Eventually(
				func() error {
					return k8sClient.Get(
						ctx,
						client.ObjectKey{
							Name:      podName,
							Namespace: ns.Name,
						},
						pod,
					)
				},
				time.Second*5,
				time.Millisecond*500,
			).ShouldNot(BeNil())
		})
	})
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).ToNot(HaveOccurred())
})

// SetupTest will set up a testing environment.
// This includes:
// * creating a Namespace to be used during the test
// * starting the 'MyKindReconciler'
// * stopping the 'MyKindReconciler" after the test ends
// Call this function at the start of each of your tests.
func SetupTest(ctx context.Context) *core.Namespace {
	var stopCh chan struct{}
	ns := &core.Namespace{}

	BeforeEach(func() {
		stopCh = make(chan struct{})
		*ns = core.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: "testns-" + randStringRunes(5)},
		}

		err := k8sClient.Create(ctx, ns)
		Expect(err).NotTo(HaveOccurred(), "failed to create test namespace")

		mgr, err := ctrl.NewManager(cfg, ctrl.Options{})
		Expect(err).NotTo(HaveOccurred(), "failed to create manager")

		controller := &SmeeProxyReconciler{
			Client: mgr.GetClient(),
			Log:    logf.Log,
			//Recorder: mgr.GetEventRecorderFor("mykind-controller"),
		}
		err = controller.SetupWithManager(mgr)
		Expect(err).NotTo(HaveOccurred(), "failed to setup controller")

		go func() {
			err := mgr.Start(stopCh)
			Expect(err).NotTo(HaveOccurred(), "failed to start manager")
		}()
	})

	AfterEach(func() {
		close(stopCh)

		err := k8sClient.Delete(ctx, ns)
		Expect(err).NotTo(HaveOccurred(), "failed to delete test namespace")
	})

	return ns
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
