package processor

import (
	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/yaml"
	"sigs.k8s.io/kustomize/api/hasher"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/resource"
)

type Merge struct {
	DefaultProcessor
	merge *resource.Resource
}

func (processor *Merge) BeforeResource(res *data.Resource) error {
	str, err := res.Content.ToString()
	if err != nil {
		return err
	}

	rf := resource.NewFactory(&hasher.Hasher{})
	rmf := resmap.NewFactory(rf)
	k8sRes, err := rf.FromBytes([]byte(str))
	if err != nil {
		return err
	}
	m := rmf.FromResource(k8sRes)

	err = m.ApplySmPatch(resource.MakeIdSet([]*resource.Resource{k8sRes}), processor.merge)
	if err != nil {
		return err
	}

	transformed, err := m.GetById(k8sRes.CurId())
	if err != nil {
		return err
	}

	y, err := transformed.AsYAML()
	if err != nil {
		return err
	}

	mf, err := data.ReadManifestString(y)
	if err != nil {
		return err
	}
	res.Content = mf
	return nil
}

func ActivateMerge(registry *ProcessorTypes) {
	registry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "merge",
			Description: "Use kustomize style strategic merge",
			Doc:         addDoc,
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			get, _ := config.Get("merge")
			raw, err := yaml.Marshal(get)
			if err != nil {
				return nil, err
			}

			rf := resource.NewFactory(&hasher.Hasher{})
			mergeDef, err := rf.FromBytes(raw)
			if err != nil {
				return nil, err
			}

			return &Merge{
				merge: mergeDef,
			}, nil
		},
	})
}
