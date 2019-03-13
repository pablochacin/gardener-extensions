package main

import (
	"fmt"
	"github.com/gardener/gardener-extensions/controllers/os-suse-jeos/cmd/gardener-extension-os-suse-jeos/app"
	"github.com/gardener/gardener-extensions/pkg/controller"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func main() {
	log.SetLogger(log.ZapLogger(false))
	cmd := app.NewControllerCommand(controller.SetupSignalHandlerContext())

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
