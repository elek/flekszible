package processor

import (
	"errors"
	"fmt"
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
	"io/ioutil"
)

type Composit struct {
	DefaultProcessor
	Processors []Processor
	File       string
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

func (c *Composit) Before(ctx *RenderContext, resources []data.Resource) {}
func (c *Composit) After(ctx *RenderContext, resources []data.Resource)  {}

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


//pase definition yaml file and register definitions to the global registry.
func parseDefintion(path string) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	mapSlice := yaml.MapSlice{}
	err = yaml.Unmarshal(content, &mapSlice)
	if err != nil {
		return err
	}
	if deftype, found := mapSlice.Get("type"); found {
		transformations, hasDefintions := mapSlice.Get("transformations")
		if !hasDefintions {
			return errors.New(fmt.Sprintf("'transformations' key is missing from definition file %s. Please define transformations under the tranformations key.", path))

		}
		rawData, err := yaml.Marshal(transformations)
		if err != nil {
			return errors.New(fmt.Sprintf("Internal error during reread the definitions part of file %s: %s", path, err.Error()))
		}
		processors, err := ReadProcessorDefinition(rawData)
		composit := Composit{
			Processors: processors,
		}
		factory := func() Processor {
			return &composit
		}
		ProcessorTypeRegistry.AddComposit(deftype.(string), factory)
		return nil
	} else {
		return errors.New(fmt.Sprintf("'type' key is missing from definition file %s. Please define a unique identifier with type.", path))
	}
}

