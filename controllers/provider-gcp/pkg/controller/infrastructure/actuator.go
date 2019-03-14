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
	"bytes"
	"context"
	"encoding/json"
	gcpv1alpha1 "github.com/gardener/gardener-extensions/controllers/provider-gcp/pkg/apis/gcp/v1alpha1"
	"github.com/gardener/gardener-extensions/controllers/provider-gcp/pkg/gcp"
	"github.com/gardener/gardener-extensions/pkg/controller"
	"github.com/gardener/gardener-extensions/pkg/controller/infrastructure"
	gardencorev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/gardener/gardener/pkg/chartrenderer"
	glogger "github.com/gardener/gardener/pkg/logger"
	"github.com/gardener/gardener/pkg/operation/terraformer"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// DefaultVPCName is the default VPC terraform name.
	DefaultVPCName = "${google_compute_network.network.name}"
	// TerraformVarServiceAccount is the name of the terraform service account environment variable.
	TerraformVarServiceAccount = "TF_VAR_SERVICEACCOUNT"
)

type actuator struct {
	client           client.Client
	decoder          runtime.Decoder
	scheme           *runtime.Scheme
	restConfig       *rest.Config
	chartRenderer    chartrenderer.ChartRenderer
	terraformerImage string
}

// NewActuator creates a new infrastructure.Actuator that instantiates Terraformers with the given
// terraformerImage.
func NewActuator(terraformerImage string) infrastructure.Actuator {
	return &actuator{
		terraformerImage: terraformerImage,
	}
}

// InjectScheme implements inject.Scheme.
func (a *actuator) InjectScheme(scheme *runtime.Scheme) error {
	a.scheme = scheme
	a.decoder = serializer.NewCodecFactory(a.scheme).UniversalDecoder()
	return nil
}

// InjectClient implements inject.Client.
func (a *actuator) InjectClient(client client.Client) error {
	a.client = client
	return nil
}

// InjectConfig implements inject.Config.
func (a *actuator) InjectConfig(config *rest.Config) error {
	a.restConfig = config

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	chartRenderer, err := chartrenderer.New(kubeClient)
	if err != nil {
		return err
	}
	a.chartRenderer = chartRenderer

	return nil
}

func (a *actuator) newTerraformer(purpose, namespace, name string) (*terraformer.Terraformer, error) {
	return terraformer.NewForConfig(glogger.NewLogger("info"), a.restConfig, purpose, namespace, name, a.terraformerImage)
}

func (a *actuator) updateProviderStatus(
	ctx context.Context,
	tf *terraformer.Terraformer,
	infra *extensionsv1alpha1.Infrastructure,
	config *gcpv1alpha1.InfrastructureConfig,
) error {
	outputKeys := []string{
		gcp.TerraformerOutputKeyVPCName,
		gcp.TerraformerOutputKeySubnetNodes,
		gcp.TerraformerOutputKeyServiceAccountEmail,
	}
	hasInternal := config.Networks.Internal != nil
	if hasInternal {
		outputKeys = append(outputKeys, gcp.TerraformerOutputKeySubnetInternal)
	}

	vars, err := tf.GetStateOutputVariables(outputKeys...)
	if err != nil {
		return err
	}

	infra.Status.ProviderStatus = &runtime.RawExtension{Object: ComputeInfrastructureStatusFromTerraformVars(vars, hasInternal)}
	return a.client.Update(ctx, infra)
}

// ComputeInfrastructureStatusFromTerraformVars computes an InfrastructureStatus from the given
// Terraform variables.
func ComputeInfrastructureStatusFromTerraformVars(vars map[string]string, hasInternal bool) *gcpv1alpha1.InfrastructureStatus {
	var (
		serviceAccountEmail = vars[gcp.TerraformerOutputKeyServiceAccountEmail]
		status              = &gcpv1alpha1.InfrastructureStatus{
			Networks: &gcpv1alpha1.NetworkStatus{
				VPC: gcpv1alpha1.VPC{
					Name: vars[gcp.TerraformerOutputKeyVPCName],
				},
				Subnets: []gcpv1alpha1.Subnet{
					{
						Purpose: gcpv1alpha1.PurposeNodes,
						Name:    vars[gcp.TerraformerOutputKeySubnetNodes],
					},
				},
			},
			ServiceAccountEmail: &serviceAccountEmail,
		}
	)

	if hasInternal {
		status.Networks.Subnets = append(status.Networks.Subnets, gcpv1alpha1.Subnet{
			Purpose: gcpv1alpha1.PurposeInternal,
			Name:    vars[gcp.TerraformerOutputKeySubnetInternal],
		})
	}
	return status
}

// getK8SNetworks gets the K8SNetworks from the given controller.Cluster.
func getK8SNetworks(cluster *controller.Cluster) *gardencorev1alpha1.K8SNetworks {
	return &cluster.Shoot.Spec.Cloud.GCP.Networks.K8SNetworks
}

// ComputeTerraformerVariablesEnvironment computes the Terraformer variables environment from the
// given ServiceAccount.
func ComputeTerraformerVariablesEnvironment(serviceAccountJSON []byte) (map[string]string, error) {
	var buf bytes.Buffer
	if err := json.Compact(&buf, serviceAccountJSON); err != nil {
		return nil, err
	}

	return map[string]string{
		TerraformVarServiceAccount: buf.String(),
	}, nil
}

// ComputeTerraformerChartValues computes the values for the GCP Terraformer chart.
func ComputeTerraformerChartValues(
	infra *extensionsv1alpha1.Infrastructure,
	project string,
	config *gcpv1alpha1.InfrastructureConfig,
	cluster *controller.Cluster,
) map[string]interface{} {
	var (
		vpcName   = DefaultVPCName
		createVPC = true
	)

	networks := getK8SNetworks(cluster)

	if config.Networks.VPC != nil {
		createVPC = false
		vpcName = config.Networks.VPC.Name
	}

	return map[string]interface{}{
		"google": map[string]interface{}{
			"region":  infra.Spec.Region,
			"project": project,
		},
		"create": map[string]interface{}{
			"vpc": createVPC,
		},
		"vpc": map[string]interface{}{
			"name": vpcName,
		},
		"clusterName": cluster.Shoot.Name,
		"networks": map[string]interface{}{
			"pods":     networks.Pods,
			"services": networks.Services,
			"worker":   config.Networks.Worker,
			"internal": config.Networks.Internal,
		},
	}
}
