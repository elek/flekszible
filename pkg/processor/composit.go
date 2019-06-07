package processor

import (
	"errors"
	"fmt"
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
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

func (c *Composit) Before(ctx *RenderContext, resources []*data.Resource) {}
func (c *Composit) After(ctx *RenderContext, resources []*data.Resource)  {}

func (c *Composit) BeforeResource(resource *data.Resource) {
	for _, p := range c.Processors {
		p.BeforeResource(resource)
	}
}

func (c *Composit) AfterResource(resource *data.Resource) {
	for _, p := range c.Processors {
		p.AfterResource(resource)
	}
}
func parseTransformationParameters(config *yaml.MapSlice) map[string]interface{} {
	result := make(map[string]interface{})
	for _, item := range *config {
		result[item.Key.(string)] = item.Value
	}
	return result

}
func compositFactory(config *yaml.MapSlice, templateBytes []byte) (Processor, error) {
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
	fmt.Println(output.String())
	if err != nil {
		panic("The composit factory can't be parsed" + err.Error())
	}
	return &Composit{
		Processors: processors,
	}, nil
}

func addDefaultParameters(parameters map[string]interface{}) {
	kubeConfig := data.CreateKubeConfig();
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
			return compositFactory(config, body);
		},
	})
	return nil
}

//
//
////pase definition yaml file and register definitions to the global registry.
//func parseDefintion(path string) error {
//	content, err := ioutil.ReadManifestFile(path)
//	if err != nil {
//		return err
//	}
//	mapSlice := yaml.MapSlice{}
//	err = yaml.Unmarshal(content, &mapSlice)
//	if err != nil {
//		return err
//	}
//	if deftype, found := mapSlice.Get("type"); found {
//		transformations, hasDefintions := mapSlice.Get("transformations")
//		if !hasDefintions {
//			return errors.New(fmt.Sprintf("'transformations' key is missing from definition file %s. Please define transformations under the tranformations key.", path))
//
//		}
//		rawData, err := yaml.Marshal(transformations)
//		if err != nil {
//			return errors.New(fmt.Sprintf("Internal error during reread the definitions part of file %s: %s", path, err.Error()))
//		}
//		processors, err := ReadProcessorDefinition(rawData)
//		composit := Composit{
//			Processors: processors,
//		}
//		factory := func() Processor {
//			return &composit
//		}
//		ProcessorTypeRegistry.AddComposit(deftype.(string), factory)
//		return nil
//	} else {
//		return errors.New(fmt.Sprintf("'type' key is missing from definition file %s. Please define a unique identifier with type.", path))
//	}
//}
