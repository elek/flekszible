package processor

import (
	"fmt"
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type RenderContext struct {
	OutputDir     string
	Conf          data.Configuration
	Mode          string
	ImageOverride string
	Namespace     string
	RootResource  *ResourceNode
}

type ResourceNode struct {
	Dir                      string
	Destination              string
	Resources                []data.Resource
	Children                 []*ResourceNode
	PreImportTransformations []byte
	ProcessorRepository      *ProcessorRepository
}

func CreateRenderContext(mode string, inputDir string, outputDir string) *RenderContext {
	return &RenderContext{
		OutputDir:    outputDir,
		Mode:         mode,
		RootResource: CreateResourceNode(inputDir, ""),
	}
}

func (context *RenderContext) LoadResourceTree() error {
	return context.RootResource.LoadResourceConfig()
}

//List all the resources from the resource tree.
func (context *RenderContext) Resources() []data.Resource {
	result := make([]data.Resource, 0)
	result = append(result, context.RootResource.AllResources()...)
	return result
}

func (node *ResourceNode) AllResources() []data.Resource {
	result := make([]data.Resource, 0)
	result = append(result, node.Resources...)
	for _, child := range node.Children {
		result = append(result, child.AllResources()...)
	}
	return result
}

func CreateResourceNode(dir string, destination string) *ResourceNode {
	node := ResourceNode{
		Dir:                 dir,
		ProcessorRepository: CreateProcessorRepository(),
		Children:            make([]*ResourceNode, 0),
		Destination:         destination,
	}
	return &node
}
func (context *RenderContext) AddResources(resources ...data.Resource) {
	context.RootResource.Resources = append(context.RootResource.Resources, resources...)
}

func (context *RenderContext) Init() error {
	err := context.LoadResourceTree()
	if err != nil {
		return err
	}
	context.LoadDefinitions()
	return context.InitializeTransformations()
}

func (context *RenderContext) InitializeTransformations() error {
	return context.RootResource.InitializeTransformations()
}

func (node *ResourceNode) InitializeTransformations() error {
	if node.PreImportTransformations != nil {
		processors, err := ReadProcessorDefinition(node.PreImportTransformations)
		if err != nil {
			return err
		}
		node.ProcessorRepository.AppendAll(processors)
	}
	node.ProcessorRepository.ParseProcessors(node.Dir)
	for _, child := range node.Children {
		err := child.InitializeTransformations()
		if err != nil {
			return err
		}
	}
	return nil
}

func (context *RenderContext) AppendProcessor(processor Processor) {
	root := context.RootResource
	repo := root.ProcessorRepository
	repo.Append(processor)
	fmt.Println(repo)
	fmt.Println(root)
}

type execute func(transformation Processor, context *RenderContext, resources []data.Resource)

func (node *ResourceNode) Execute(context *RenderContext, functionCall execute) []data.Resource {
	resources := make([]data.Resource, 0)
	for _, child := range node.Children {
		resources = append(resources, child.Execute(context, functionCall)...)
	}
	resources = append(resources, node.Resources...)
	for _, transformation := range node.ProcessorRepository.Processors {
		functionCall(transformation, context, resources)
	}
	return resources

}
func (ctx *RenderContext) Render() {
	before := func(processor Processor, context *RenderContext, resources []data.Resource) {
		processor.Before(context, resources)
	}
	after := func(processor Processor, context *RenderContext, resources []data.Resource) {
		processor.After(context, resources)
	}
	process := func(processor Processor, context *RenderContext, resources []data.Resource) {
		for _, resource := range resources {
			processor.BeforeResource(&resource)
			resource.Content.Accept(processor)
			processor.AfterResource(&resource)

		}
		processor.After(context, resources)
	}
	ctx.RootResource.Execute(ctx, before)
	ctx.RootResource.Execute(ctx, process)
	ctx.RootResource.Execute(ctx, after)
}

//parse the directory structure and the flekszible configs from the dirs
func (node *ResourceNode) LoadResourceConfig() error {
	configFile := path.Join(node.Dir, "flekszible.yaml")
	conf, err := data.ReadConfiguration(configFile)
	if err != nil {
		return err
	}
	node.Resources = data.ReadResourcesFromDir(node.Dir)
	for ix, _ := range node.Resources {
		node.Resources[ix].Destination = node.Destination
	}
	for _, importDefinition := range conf.Import {
		var importedDir string
		if !filepath.IsAbs(importDefinition.Path) {
			importedDir = path.Join(node.Dir, importDefinition.Path)
		} else {
			importedDir = importDefinition.Path
		}
		childNode := CreateResourceNode(importedDir, importDefinition.Destination)
		err := childNode.LoadResourceConfig()
		if err != nil {
			return err
		}
		if len(importDefinition.Transformations) > 0 {
			bytes, err := yaml.Marshal(importDefinition.Transformations)
			if err != nil {
				return err
			}
			childNode.PreImportTransformations = bytes
		}
		node.Children = append(node.Children, childNode)
	}
	return nil
}

//load transformation definitions from ./definitions dir (all dir)
func (ctx *RenderContext) LoadDefinitions() {
	ctx.RootResource.LoadDefinitions()
}

func (node *ResourceNode) LoadDefinitions() {
	defDir := path.Join(node.Dir, "definitions")
	if _, err := os.Stat(defDir); !os.IsNotExist(err) {
		files, err := ioutil.ReadDir(defDir)
		if err != nil {
			logrus.Warningf("Can't load definition directory %s: %s", defDir, err.Error())
		}
		for _, file := range files {
			definitionFile := path.Join(defDir, file.Name())
			err := parseDefintion(definitionFile)
			if err != nil {
				logrus.Errorf("Can't parse the definition file %s: %s", definitionFile, err.Error())
			}
		}
	}
	for _, child := range node.Children {
		child.LoadDefinitions()
	}

}
