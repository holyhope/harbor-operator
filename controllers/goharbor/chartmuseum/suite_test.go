package chartmuseum_test

import (
	"context"
	"path"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/goharbor/harbor-operator/controllers"
	"github.com/goharbor/harbor-operator/controllers/goharbor/chartmuseum"
	"github.com/goharbor/harbor-operator/controllers/goharbor/internal/test"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	stopCh      chan struct{}
	ctx         context.Context
	reconciler  *chartmuseum.Reconciler
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

	name := controllers.ChartMuseum.String()

	configStore, _ := test.NewConfig(ctx, chartmuseum.ConfigTemplatePathKey, path.Base(chartmuseum.DefaultConfigTemplatePath))
	/* TODO: Enable this when HarborClass is fixed
	provider.Add(configstore.NewItem(config.HarborClassKey, harborClass, 100))
	*/
	configStore.Env(name)

	reconciler, err := chartmuseum.New(ctx, name, configStore)
	Expect(err).ToNot(HaveOccurred())

	ctx, harborClass, stopCh = test.StartController(ctx, reconciler)

	close(done)
}, 60)

var _ = AfterSuite(func() {
	defer test.AfterSuite(ctx)

	close(stopCh)
})
