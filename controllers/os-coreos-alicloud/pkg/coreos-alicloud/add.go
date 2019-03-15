// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package coreos

import (
	extensionscontroller "github.com/gardener/gardener-extensions/pkg/controller"
	"github.com/gardener/gardener-extensions/pkg/controller/operatingsystemconfig"

	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Type is the type of operating system configs the CoreOS Alicloud controller monitors.
const Type = "coreos-alicloud"

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
		ControllerOptions: opts,
		Predicates:        operatingsystemconfig.DefaultPredicates(Type),
	})
}

// Add adds a controller with the default Options.
func Add(mgr manager.Manager) error {
	return AddWithOptions(mgr, Options)
}
