package pkg

import (
	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/processor"
	"github.com/elek/flekszible/api/v2/yaml"
)

type ResourceList struct {
	Items  []interface{}           `yaml:"items"`
	Config data.FlekszibleResource `yaml:"functionConfig"`
}

// Kustomize provides kustomize integration: reads everything from stdinput and generates output
func Kustomize(raw []byte) error {
	context := processor.CreateRenderContext("k8s", "-", "-")

	var resources []*data.Resource
	rl := ResourceList{}
	err := yaml.Unmarshal(raw, &rl)
	if err != nil {
		return err
	}

	for _, r := range rl.Items {
		rawResource, err := yaml.Marshal(r)
		if err != nil {
			return err
		}
		parsedFragment, err := data.ReadManifestString(rawResource)
		if err != nil {
			return err
		}

		r := data.NewResource()
		r.Content = parsedFragment
		resources = append(resources, &r)
	}

	s := data.NewSourceCacheManager(".")
	err = context.RootResource.InitFromConfig(rl.Config.Spec, &s, "-")
	if err != nil {
		return err
	}

	err = context.Init()
	if err != nil {
		return err
	}
	context.RootResource.Resources = resources
	AddInternalTransformations(context, false)
	err = context.Render()
	if err != nil {
		return err
	}
	return nil

}
