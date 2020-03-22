package processor

import (
	"github.com/elek/flekszible/api/data"
	"github.com/elek/flekszible/api/yaml"
	"github.com/pkg/errors"
	"io/ioutil"
	"strings"
	"text/template"
)

type Composit struct {
	DefaultProcessor
	ProcessorMetadata
	Processors []Processor
	File       string
	Template   string
	Parameters map[string]string
	Trigger    Trigger

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

func (c *Composit) Before(ctx *RenderContext, resources []*data.Resource) error { return nil }
func (c *Composit) After(ctx *RenderContext, resources []*data.Resource) error  { return nil }

func (c *Composit) BeforeResource(resource *data.Resource) error {
	if !c.Trigger.active(resource) {
		return nil
	}
	for _, p := range c.Processors {
		err := p.BeforeResource(resource)
		if err != nil {
			return errors.Wrap(err, "One of the child processors of the composite resource is failed")
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
func compositFactory(config *yaml.MapSlice, templateBytes []byte) (*Composit, error) {
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
		panic("The composit factory can't be parsed" + err.Error())
	}
	return &Composit{
		Processors: processors,
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
		return err
	}
	head, body := splitDefinitionFile(content)
	metadata := ProcessorMetadata{}
	err = yaml.Unmarshal(head, &metadata)
	if err != nil {
		return err
	}
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: metadata,
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			comp, err := compositFactory(config, body)
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
