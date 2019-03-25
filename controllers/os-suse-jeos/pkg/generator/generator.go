package generator

import (
	template_gen "github.com/gardener/gardener-extensions/controllers/os-common/pkg/template"
	"github.com/gobuffalo/packr/v2"
	"text/template"
)

//go:generate packr2

func NewCloudInitGenerator() (*template_gen.CloudInitGenerator, error) {
	box := packr.New("templates", "./templates")
	cloudInitTemplateString, err := box.FindString("cloud-init.template")
	if err != nil {
		return nil, err
	}

	cloudInitTemplate, err := template.New("cloud-init").Parse(cloudInitTemplateString)
	if err != nil {
		return nil, err
	}
	generator := template_gen.NewCloudInitGenerator(cloudInitTemplate, template_gen.DefaultUnitsPath)
	return generator, nil
}
