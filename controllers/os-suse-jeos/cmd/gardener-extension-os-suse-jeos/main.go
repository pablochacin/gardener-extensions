package main

import (
	"fmt"
	"github.com/gardener/gardener-extensions/controllers/os-common/pkg/app"
	"github.com/gardener/gardener-extensions/controllers/os-suse-jeos/pkg/generator"
	extcontroller "github.com/gardener/gardener-extensions/pkg/controller"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

// Default options
var options = &controller.Options{}

func main() {
	log.SetLogger(log.ZapLogger(false))

	g, err := generator.NewCloudInitGenerator()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cmd := app.NewControllerCommand(extcontroller.SetupSignalHandlerContext(), "suse-jeos", g, options)

	if err = cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
