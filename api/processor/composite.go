package processor

import (
	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/yaml"
	"github.com/pkg/errors"
	"io/ioutil"
	gopath "path"
	"path/filepath"
	"strings"
	"text/template"
)

type Composite struct {
	DefaultProcessor
	ProcessorMetadata
	AdditionalResourcesDir string
	Processors             []Processor
	File                   string
	Template               string
	Parameters             map[string]string
	Trigger                Trigger
}

func (processor *Composite) ToString() string {
	return processor.ProcessorMetadata.Name + " (composite)"
}

func (c *Composite) OnKey(node *data.KeyNode) {
	for _, p := range c.Processors {
		p.OnKey(node)
	}
}
func (c *Composite) BeforeMap(node *data.MapNode) {
	for _, p := range c.Processors {
		p.BeforeMap(node)
	}
}
func (c *Composite) AfterMap(node *data.MapNode)                                   {}
func (c *Composite) BeforeMapItem(node *data.MapNode, key string, index int)       {}
func (c *Composite) AfterMapItem(node *data.MapNode, key string, index int)        {}
func (c *Composite) BeforeList(node *data.ListNode)                                {}
func (c *Composite) AfterList(node *data.ListNode)                                 {}
func (c *Composite) BeforeListItem(node *data.ListNode, item data.Node, index int) {}
func (c *Composite) AfterListItem(node *data.ListNode, item data.Node, index int)  {}

func (c *Composite) Before(ctx *RenderContext, node *ResourceNode) error { return nil }
func (c *Composite) After(ctx *RenderContext, node *ResourceNode) error  { return nil }

func (c *Composite) RegisterResources(ctx *RenderContext, node *ResourceNode) error {
	if c.AdditionalResourcesDir != "" {
		resources := data.ReadResourcesFromDir(c.AdditionalResourcesDir)
		node.Resources = append(node.Resources, resources...)
	}
	return nil
}

func (c *Composite) BeforeResource(resource *data.Resource, location *ResourceNode) error {
	if !c.Trigger.active(resource) {
		return nil
	}
	for _, p := range c.Processors {
		err := p.BeforeResource(resource, location)
		if err != nil {
			return errors.Wrap(err, "Resource transformation is failed "+p.ToString()+" on "+resource.Name())
		}
	}
	return nil
}

func (c *Composite) AfterResource(resource *data.Resource, location *ResourceNode) error {
	if !c.Trigger.active(resource) {
		return nil
	}
	for _, p := range c.Processors {
		err := p.AfterResource(resource, location)
		if err != nil {
			return errors.Wrap(err, "One of the child processors of the composite resource is failed")
		}

	}
	return nil
}
func parseTransformationParameters(metadata *ProcessorMetadata, config *yaml.MapSlice) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	for _, paramDef := range metadata.Parameters {
		if len(paramDef.Default) > 0 {
			result[paramDef.Name] = paramDef.Default
		}
	}

outer:
	for _, item := range *config {
		parameterName := item.Key.(string)
		if parameterName == "type" || parameterName == "scope" || parameterName == "trigger" {
			continue
		}
		validParamNames := make([]string, 0)
		for _, paramDef := range metadata.Parameters {
			if parameterName == paramDef.Name {
				result[parameterName] = item.Value
				continue outer
			} else {
				validParamNames = append(validParamNames, paramDef.Name)
			}
		}
		return result, errors.New("Unknown parameter '" + parameterName + "' used for composite transformations " + metadata.Name + ", valid parameters: " + strings.Join(validParamNames, ","))

	}
	for _, paramDef := range metadata.Parameters {
		if paramDef.Required {
			if _, found := result[paramDef.Name]; !found {
				return result, errors.New("Parameters " + paramDef.Name + " is required for composite transformation " + metadata.Name)
			}
		}
	}
	return result, nil

}

//instantiate composite transformation based on instance config, generic definition metadata and template
func compositeFactory(path string, metadata *ProcessorMetadata, config *yaml.MapSlice, templateBytes []byte) (*Composite, error) {
	funcmap := template.FuncMap{
		"Iterate": func(count int) []int {
			var i int
			var Items []int
			for i = 0; i < count; i++ {
				Items = append(Items, i)
			}
			return Items
		},
	}

	tpl, err := template.New("definition").Funcs(funcmap).Parse(string(templateBytes))
	if err != nil {
		return nil, errors.New("The definition template is invalid: " + err.Error())
	}
	output := strings.Builder{}
	parameters, err := parseTransformationParameters(metadata, config)
	if err != nil {
		return nil, errors.Wrap(err, "Couldn't set parameter for composite transformations")
	}
	addDefaultParameters(parameters)
	err = tpl.Execute(&output, parameters)
	if err != nil {
		return nil, errors.New("The render was failed: " + err.Error())
	}
	processors, err := ReadProcessorDefinition([]byte(output.String()))
	if err != nil {
		panic("The definition can't be parsed" + err.Error())
	}
	resourcesDir := metadata.Resources
	if resourcesDir != "" && !filepath.IsAbs(resourcesDir) {
		resourcesDir = filepath.Clean(gopath.Join(gopath.Dir(path), resourcesDir))
	}
	return &Composite{
		ProcessorMetadata:      *metadata,
		Processors:             processors,
		AdditionalResourcesDir: resourcesDir,
	}, nil
}

func addDefaultParameters(parameters map[string]interface{}) {
	kubeConfig := data.CreateKubeConfig()
	ns, err := kubeConfig.ReadCurrentNamespace()
	if err != nil {
		ns = "default"
	}
	parameters["namespace"] = ns
}

//pase definition yaml file and register definitions to the global registry.
func parseDefintion(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrap(err, "Can't load transformation definition from "+path)
	}
	head, body := splitDefinitionFile(content)
	metadata := ProcessorMetadata{}
	err = yaml.Unmarshal(head, &metadata)
	if err != nil {
		return "", errors.Wrap(err, "Can't parse transformation metadata from "+path)
	}
	metadata.Doc = metadata.Doc + "\n### Definition: \n" + string(content)
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: metadata,
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			comp, err := compositeFactory(path, &metadata, config, body)
			if err != nil {
				return nil, err
			}
			trigger, found := config.Get("trigger")
			if found {
				node, err := data.ConvertToNode(trigger, data.NewPath())
				if err != nil {
					return nil, err
				}
				comp.Trigger = Trigger{Definition: node}
			}
			return comp, nil
		},
	})
	return metadata.Name, nil
}
