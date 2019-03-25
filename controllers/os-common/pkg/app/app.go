package app

import (
	"context"
	"github.com/gardener/gardener-extensions/controllers/os-common/pkg/actuator"
	"github.com/gardener/gardener-extensions/controllers/os-common/pkg/generator"
	extcontroller "github.com/gardener/gardener-extensions/pkg/controller"
	controllercmd "github.com/gardener/gardener-extensions/pkg/controller/cmd"
	"github.com/spf13/cobra"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// NewControllerCommand creates a new command for running an OS controller.
func NewControllerCommand(ctx context.Context, osName string, generator generator.Generator, options *controller.Options) *cobra.Command {
	var (
		restOpts = &controllercmd.RESTOptions{}
		mgrOpts  = &controllercmd.ManagerOptions{
			LeaderElection:          true,
			LeaderElectionID:        controllercmd.LeaderElectionNameID(osName),
			LeaderElectionNamespace: os.Getenv("LEADER_ELECTION_NAMESPACE"),
		}
		ctrlOpts = &controllercmd.ControllerOptions{
			MaxConcurrentReconciles: 5,
		}

		aggOption = controllercmd.NewOptionAggregator(restOpts, mgrOpts, ctrlOpts)
	)

	cmd := &cobra.Command{
		Use: "os-" + osName + "-controller-manager",

		Run: func(cmd *cobra.Command, args []string) {
			if err := aggOption.Complete(); err != nil {
				controllercmd.LogErrAndExit(err, "Error completing options")
			}

			mgr, err := manager.New(restOpts.Completed().Config, mgrOpts.Completed().Options())
			if err != nil {
				controllercmd.LogErrAndExit(err, "Could not instantiate manager")
			}

			if err := extcontroller.AddToScheme(mgr.GetScheme()); err != nil {
				controllercmd.LogErrAndExit(err, "Could not update manager scheme")
			}

			ctrlOpts.Completed().Apply(options)

			if err := actuator.Add(mgr, osName, generator); err != nil {
				controllercmd.LogErrAndExit(err, "Could not add controller to manager")
			}

			if err := mgr.Start(ctx.Done()); err != nil {
				controllercmd.LogErrAndExit(err, "Error running manager")
			}
		},
	}

	aggOption.AddFlags(cmd.Flags())

	return cmd
}
