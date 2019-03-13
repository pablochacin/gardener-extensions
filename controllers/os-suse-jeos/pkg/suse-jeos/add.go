package jeos

import (
	extensionscontroller "github.com/gardener/gardener-extensions/pkg/controller"
	"github.com/gardener/gardener-extensions/pkg/controller/operatingsystemconfig"

	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Type is the type of operating system configs the controller monitors.
const Type = "suse-jeos"

func init() {
	addToManagerBuilder.Register(Add)
}

var (
	addToManagerBuilder = extensionscontroller.NewAddToManagerBuilder()
	// AddToManager adds all coreos-alicloud controllers to the given manager.
	AddToManager = addToManagerBuilder.AddToManager

	// Options are the default controller.Options for Add.
	Options = controller.Options{}
)

// AddWithOptions adds a controller with the given Options to the given manager.
// The opts.Reconciler is being set with a newly instantiated actuator.
func AddWithOptions(mgr manager.Manager, opts controller.Options) error {
	return operatingsystemconfig.Add(mgr, operatingsystemconfig.AddArgs{
		Actuator:          NewActuator(),
		Type:              Type,
		ControllerOptions: opts,
	})
}

// Add adds a controller with the default Options.
func Add(mgr manager.Manager) error {
	return AddWithOptions(mgr, Options)
}
