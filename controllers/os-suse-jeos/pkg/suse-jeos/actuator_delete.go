package jeos

import (
	"context"

	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
)

func (a *actuator) delete(ctx context.Context, config *extensionsv1alpha1.OperatingSystemConfig) error {
	return nil
}
