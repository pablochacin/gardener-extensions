package actuator

import (
	"github.com/gardener/gardener-extensions/controllers/os-common/pkg/generator"
	"github.com/gardener/gardener-extensions/pkg/controller/operatingsystemconfig"

	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Default options
var options = controller.Options{}

// AddWithOptions adds a controller with the given Options to the given manager.
// The opts.Reconciler is being set with a newly instantiated actuator.
func AddWithOptions(mgr manager.Manager, os string, generator generator.Generator, opts controller.Options) error {
	return operatingsystemconfig.Add(mgr, operatingsystemconfig.AddArgs{
		Actuator:          NewActuator(os, generator),
		Type:              os,
		ControllerOptions: opts,
	})
}

// Add adds a controller with the default Options.
func Add(mgr manager.Manager, os string, generator generator.Generator) error {
	return AddWithOptions(mgr, os, generator, options)
}
