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

package gcp

import "path/filepath"

const (
	// TerraformerImageName is the name of the Terraformer image.
	TerraformerImageName = "terraformer"

	// TerraformerPurposeInfra is the terraformer infrastructure purpose.
	TerraformerPurposeInfra = "infra"

	// TerraformerOutputKeyVPCName is the name of the vpc_name terraform output variable.
	TerraformerOutputKeyVPCName = "vpc_name"
	// TerraformerOutputKeyServiceAccountEmail is the name of the service_account_email terraform output variable.
	TerraformerOutputKeyServiceAccountEmail = "service_account_email"
	// TerraformerOutputKeySubnetNodes is the name of the subnet_nodes terraform output variable.
	TerraformerOutputKeySubnetNodes = "subnet_nodes"
	// TerraformerOutputKeySubnetInternal is the name of the subnet_internal terraform output variable.
	TerraformerOutputKeySubnetInternal = "subnet_internal"
)

var (
	// ChartsPath is the path to the charts
	ChartsPath = filepath.Join("controllers", "provider-gcp", "charts")
	// InternalChartsPath is the path to the internal charts
	InternalChartsPath = filepath.Join(ChartsPath, "internal")
)
