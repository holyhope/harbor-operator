/*
Copyright 2019 The Kubernetes Authors.

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

package core_test

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	goharborv1alpha2 "github.com/goharbor/harbor-operator/apis/goharbor.io/v1alpha2"
	harbormetav1 "github.com/goharbor/harbor-operator/apis/meta/v1alpha1"
	"github.com/goharbor/harbor-operator/controllers/goharbor/internal/test"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

const defaultGenerationNumber int64 = 1

var _ = Describe("Core", func() {
	var (
		ns   = test.InitNamespace(func() context.Context { return ctx })
		core goharborv1alpha2.Core
	)

	BeforeEach(func() {
		core.ObjectMeta = metav1.ObjectMeta{
			Name:      test.NewName("core"),
			Namespace: ns.GetName(),
			/* TODO: Enable this when HarborClass is fixed
			Annotations: map[string]string{
				goharborv1alpha2.HarborClassAnnotation: harborClass,
			},
			*/
		}
	})

	JustAfterEach(test.LogsAll(&ctx, func() types.NamespacedName {
		return types.NamespacedName{
			Name:      reconciler.NormalizeName(ctx, core.GetName()),
			Namespace: core.GetNamespace(),
		}
	}))

	Context("Without TLS", func() {
		BeforeEach(func() {
			core.Spec = goharborv1alpha2.CoreSpec{}
		})

		It("Should works", func() {
			By("Creating new resource", func() {
				Ω(test.GetClient(ctx).Create(ctx, &core)).
					Should(test.SuccessOrExists)

				Eventually(func() error { return test.GetClient(ctx).Get(ctx, test.GetNamespacedName(&core), &core) }, time.Minute, 5*time.Second).
					Should(Succeed(), "resource should exists")

				Ω(core.GetGeneration()).
					Should(Equal(defaultGenerationNumber), "Generation should not be updated")

				test.EnsureReady(ctx, &core, time.Minute, 5*time.Second)

				IntegTest(ctx, &core)
			})

			By("Updating resource spec", func() {
				oldGeneration := core.GetGeneration()

				test.ScaleUp(ctx, &core)

				Ω(core.GetGeneration()).
					Should(BeNumerically(">", oldGeneration), "ObservedGeneration should be updated")

				Ω(test.GetClient(ctx).Get(ctx, test.GetNamespacedName(&core), &core)).
					Should(Succeed(), "resource should still be accessible")

				test.EnsureReady(ctx, &core, time.Minute, 5*time.Second)

				IntegTest(ctx, &core)
			})

			By("Deleting resource", func() {
				Ω(test.GetClient(ctx).Delete(ctx, &core)).
					Should(Succeed())

				Eventually(func() error {
					return test.GetClient(ctx).Get(ctx, test.GetNamespacedName(&core), &core)
				}, time.Minute, 5*time.Second).
					ShouldNot(Succeed(), "Resource should no more exist")
			})
		})
	})
})

func IntegTest(ctx context.Context, core *goharborv1alpha2.Core) {
	client, err := rest.UnversionedRESTClientFor(test.NewRestConfig(ctx))
	Expect(err).ToNot(HaveOccurred())

	namespacedName := types.NamespacedName{
		Name:      reconciler.NormalizeName(ctx, core.GetName()),
		Namespace: core.GetNamespace(),
	}

	proxyReq := client.Get().
		Resource("services").
		Namespace(namespacedName.Namespace).
		Name(fmt.Sprintf("%s:%s", namespacedName.Name, harbormetav1.CoreHTTPPortName)).
		SubResource("proxy").
		Suffix("api/v2.0/health")

	Ω(proxyReq.DoRaw(ctx)).
		Should(WithTransform(func(result []byte) string {
			var health struct {
				Status string `json:"status"`
			}

			Ω(json.Unmarshal(result, &health)).
				Should(Succeed())

			return health.Status
		}, BeEquivalentTo("healthy")))
}
