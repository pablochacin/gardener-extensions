package actuator

import (
	"context"

	"github.com/gardener/gardener-extensions/controllers/os-common/pkg/generator"
	"github.com/gardener/gardener-extensions/pkg/controller/operatingsystemconfig"

	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

type actuator struct {
	scheme    *runtime.Scheme
	client    client.Client
	logger    logr.Logger
	osName    string
	generator generator.Generator
}

// NewActuator creates a new actuator with the given logger.
func NewActuator(osName string, generator generator.Generator) operatingsystemconfig.Actuator {
	return &actuator{
		logger:    log.Log.WithName(osName + "-operatingsystemconfig-actuator"),
		osName:    osName,
		generator: generator,
	}
}

func (a *actuator) InjectScheme(scheme *runtime.Scheme) error {
	a.scheme = scheme
	return nil
}

func (a *actuator) InjectClient(client client.Client) error {
	a.client = client
	return nil
}

func (a *actuator) Reconcile(ctx context.Context, config *extensionsv1alpha1.OperatingSystemConfig) ([]byte, *string, []string, error) {
	return a.reconcile(ctx, config)
}

func (a *actuator) Delete(ctx context.Context, config *extensionsv1alpha1.OperatingSystemConfig) error {
	return a.delete(ctx, config)
}