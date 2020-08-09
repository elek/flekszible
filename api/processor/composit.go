package processor

import (
	"github.com/elek/flekszible/api/data"
	"github.com/elek/flekszible/api/yaml"
	"github.com/pkg/errors"
	"io/ioutil"
	gopath "path"
	"path/filepath"
	"strings"
	"text/template"
)

type Composit struct {
	DefaultProcessor
	ProcessorMetadata
	AdditionalResourcesDir string
	Processors             []Processor
	File                   string
	Template               string
	Parameters             map[string]string
	Trigger                Trigger
}

func (c *Composit) OnKey(node *data.KeyNode) {
	for _, p := range c.Processors {
		p.OnKey(node)
	}
}
func (c *Composit) BeforeMap(node *data.MapNode) {
	for _, p := range c.Processors {
		p.BeforeMap(node)
	}
}
func (c *Composit) AfterMap(node *data.MapNode)                                   {}
func (c *Composit) BeforeMapItem(node *data.MapNode, key string, index int)       {}
func (c *Composit) AfterMapItem(node *data.MapNode, key string, index int)        {}
func (c *Composit) BeforeList(node *data.ListNode)                                {}
func (c *Composit) AfterList(node *data.ListNode)                                 {}
func (c *Composit) BeforeListItem(node *data.ListNode, item data.Node, index int) {}
func (c *Composit) AfterListItem(node *data.ListNode, item data.Node, index int)  {}

func (c *Composit) Before(ctx *RenderContext, node *ResourceNode) error { return nil }
func (c *Composit) After(ctx *RenderContext, node *ResourceNode) error  { return nil }

func (c *Composit) RegisterResources(ctx *RenderContext, node *ResourceNode) error {
	if c.AdditionalResourcesDir != "" {
		resources := data.ReadResourcesFromDir(c.AdditionalResourcesDir)
		node.Resources = append(node.Resources, resources...)
	}
	return nil
}

func (c *Composit) BeforeResource(resource *data.Resource) error {
	if !c.Trigger.active(resource) {
		return nil
	}
	for _, p := range c.Processors {
		err := p.BeforeResource(resource)
		if err != nil {
			return errors.Wrap(err, "Resource transformation is failed " + p.GetType() +" on " + resource.Name())
		}
	}
	return nil
}

func (c *Composit) AfterResource(resource *data.Resource) error {
	if !c.Trigger.active(resource) {
		return nil
	}
	for _, p := range c.Processors {
		err := p.AfterResource(resource)
		if err != nil {
			return errors.Wrap(err, "One of the child processors of the composite resource is failed")
		}

	}
	return nil
}
func parseTransformationParameters(config *yaml.MapSlice) map[string]interface{} {
	result := make(map[string]interface{})
	for _, item := range *config {
		result[item.Key.(string)] = item.Value
	}
	return result

}

//instantiate composite transformation based on instance config, generic definition metadata and template
func compositFactory(path string, metadata *ProcessorMetadata, config *yaml.MapSlice, templateBytes []byte) (*Composit, error) {
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
	parameters := parseTransformationParameters(config)
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
	return &Composit{
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
func parseDefintion(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "Can't load transformation definition from "+path)
	}
	head, body := splitDefinitionFile(content)
	metadata := ProcessorMetadata{}
	err = yaml.Unmarshal(head, &metadata)
	if err != nil {
		return errors.Wrap(err, "Can't parse transformation metadata from "+path)
	}
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: metadata,
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			comp, err := compositFactory(path, &metadata, config, body)
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
	return nil
}
