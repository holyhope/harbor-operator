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
	"path"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/goharbor/harbor-operator/controllers"
	"github.com/goharbor/harbor-operator/controllers/goharbor/core"
	"github.com/goharbor/harbor-operator/controllers/goharbor/internal/test"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	stopCh      chan struct{}
	ctx         context.Context
	reconciler  *core.Reconciler
	harborClass string
)

func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecsWithDefaultAndCustomReporters(t,
		"Controller Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func(done Done) {
	ctx = test.InitSuite()

	name := controllers.Core.String()

	configStore, _ := test.NewConfig(ctx, core.ConfigTemplatePathKey, path.Base(core.DefaultConfigTemplatePath))
	/* TODO: Enable this when HarborClass is fixed
	provider.Add(configstore.NewItem(config.HarborClassKey, harborClass, 100))
	*/
	configStore.Env(name)

	reconciler, err := core.New(ctx, name, configStore)
	Expect(err).ToNot(HaveOccurred())

	ctx, harborClass, stopCh = test.StartController(ctx, reconciler)

	close(done)
}, 60)

var _ = AfterSuite(func() {
	defer test.AfterSuite(ctx)

	close(stopCh)
})
