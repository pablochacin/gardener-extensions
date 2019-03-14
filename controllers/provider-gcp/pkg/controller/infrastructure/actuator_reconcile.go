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

package infrastructure

import (
	"context"
	"fmt"
	gcpv1alpha1 "github.com/gardener/gardener-extensions/controllers/provider-gcp/pkg/apis/gcp/v1alpha1"
	"github.com/gardener/gardener-extensions/controllers/provider-gcp/pkg/gcp"
	"github.com/gardener/gardener-extensions/pkg/controller"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/gardener/gardener/pkg/operation/terraformer"
	"path/filepath"
)

// Reconcile implements infrastructure.Actuator.
func (a *actuator) Reconcile(ctx context.Context, infra *extensionsv1alpha1.Infrastructure, cluster *controller.Cluster) error {
	config := &gcpv1alpha1.InfrastructureConfig{}
	if _, _, err := a.decoder.Decode(infra.Spec.ProviderConfig.Raw, nil, config); err != nil {
		return err
	}

	serviceAccount, err := gcp.GetServiceAccount(ctx, a.client, infra.Spec.SecretRef.Namespace, infra.Spec.SecretRef.Name)
	if err != nil {
		return err
	}

	projectID, err := gcp.ExtractServiceAccountProjectID(serviceAccount)
	if err != nil {
		return err
	}

	variablesEnvironment, err := ComputeTerraformerVariablesEnvironment(serviceAccount)
	if err != nil {
		return err
	}

	tfChartValues := ComputeTerraformerChartValues(infra, projectID, config, cluster)

	release, err := a.chartRenderer.Render(filepath.Join(gcp.InternalChartsPath, "gcp-infra"), "gcp-infra", infra.Namespace, tfChartValues)
	if err != nil {
		return err
	}

	tf, err := a.newTerraformer(gcp.TerraformerPurposeInfra, infra.Namespace, infra.Name)
	if err != nil {
		return err
	}

	err = tf.
		SetVariablesEnvironment(variablesEnvironment).
		InitializeWith(terraformer.DefaultInitializer(
			a.client,
			release.FileContent("main.tf"),
			release.FileContent("variables.tf"),
			[]byte(release.FileContent("terraform.tfvars"))),
		).
		Apply()
	if err != nil {
		return fmt.Errorf("failed to update the provider: %v", err)
	}

	return a.updateProviderStatus(ctx, tf, infra, config)
}
