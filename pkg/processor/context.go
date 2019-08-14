package processor

import (
	"github.com/pkg/errors"
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strings"
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
	Resources                []*data.Resource
	Children                 []*ResourceNode
	PreImportTransformations []byte
	Origin                   data.Source
	Source                   []data.Source
	ProcessorRepository      *TransformationRepository
	Standalone               bool
}

func CreateRenderContext(mode string, inputDir string, outputDir string) *RenderContext {
	logrus.Infof("Input dir: %s, output dir: %s", inputDir, outputDir)
	return &RenderContext{
		OutputDir:    outputDir,
		Mode:         mode,
		RootResource: CreateResourceNode(inputDir, "", &data.CurrentDir{CurrentDir: inputDir}),
	}
}

func (context *RenderContext) ListResourceNodes() []*ResourceNode {
	return listResourceNodesInt(context.RootResource)
}

func listResourceNodesInt(node *ResourceNode) []*ResourceNode {
	result := make([]*ResourceNode, 0)
	result = append(result, node)
	for _, child := range node.Children {
		result = append(result, listResourceNodesInt(child)...)
	}
	return result
}

func (context *RenderContext) LoadResourceTree() error {
	cacheManager := data.NewSourceCacheManager(context.RootResource.Dir)
	return context.RootResource.LoadResourceConfig(&cacheManager)
}

//List all the resources from the resource tree.
func (context *RenderContext) Resources() []*data.Resource {
	result := make([]*data.Resource, 0)
	result = append(result, context.RootResource.AllResources()...)
	return result
}

func (node *ResourceNode) AllResources() []*data.Resource {
	result := make([]*data.Resource, 0)
	result = append(result, node.Resources...)
	for _, child := range node.Children {
		result = append(result, child.AllResources()...)
	}
	return result
}

func CreateResourceNode(dir string, destination string, source data.Source) *ResourceNode {
	node := ResourceNode{
		Dir:                 dir,
		ProcessorRepository: CreateProcessorRepository(),
		Children:            make([]*ResourceNode, 0),
		Destination:         destination,
		Origin:              source,
	}
	return &node
}
func (context *RenderContext) AddResources(resources ...*data.Resource) {
	context.RootResource.Resources = append(context.RootResource.Resources, resources...)
}

func (context *RenderContext) Init() error {
	//load flekszible.yaml (recursive)
	err := context.LoadResourceTree()
	if err != nil {
		return err
	}
	//load 'definitions' subdirs
	context.LoadDefinitions()
	//load 'transformations' subdirs
	return context.InitializeTransformations()
}

func (context *RenderContext) InitializeTransformations() error {
	return context.RootResource.InitializeTransformations(context)
}

func (node *ResourceNode) InitializeTransformations(context *RenderContext) error {
	if node.PreImportTransformations != nil {
		processors, err := ReadProcessorDefinition(node.PreImportTransformations)
		if err != nil {
			return err
		}
		node.ProcessorRepository.AppendAll(processors)
	}
	if !node.Standalone {
		processors, e := ParseTransformations(node.Dir)
		if e != nil {
			logrus.Error("Can't read transformations from " + node.Dir + " " + e.Error())
		} else {
			for _, processor := range processors {
				if processor.GetScope() == "global" {
					context.RootResource.ProcessorRepository.Append(processor)
				} else {
					node.ProcessorRepository.InsertToBeginning(processor)
				}
			}
		}
	}
	for _, child := range node.Children {
		err := child.InitializeTransformations(context)
		if err != nil {
			return err
		}
	}
	return nil
}

func (context *RenderContext) AppendCustomProcessor(name string, parameters map[string]string) error {
	if definition, ok := ProcessorTypeRegistry.TypeMap[strings.ToLower(name)]; ok {
		config := yaml.MapSlice{}
		for k, v := range parameters {
			config.Put(k, v)
		}
		processor, err := definition.Factory(&config)
		if err != nil {
			return err
		}
		context.RootResource.ProcessorRepository.Append(processor)
		return nil
	} else {
		return errors.New("No such processor definition")
	}
}

func (context *RenderContext) AppendProcessor(processor Processor) {
	root := context.RootResource
	repo := root.ProcessorRepository
	repo.Append(processor)
}

type execute func(transformation Processor, context *RenderContext, resources []*data.Resource) error

func (node *ResourceNode) Execute(context *RenderContext, functionCall execute) ([]*data.Resource, error) {
	resources := make([]*data.Resource, 0)
	for _, child := range node.Children {
		childResults, err := child.Execute(context, functionCall)
		if err != nil {
			return nil, errors.Wrap(err, "The function can't be executed on one of the childResources")
		}
		resources = append(resources, childResults...)
	}
	resources = append(resources, node.Resources...)
	for _, transformation := range node.ProcessorRepository.Processors {
		err := functionCall(transformation, context, resources)
		if err != nil {
			return nil, errors.Wrap(err, "The function can't be executed on all the resources. Transformation: ")
		}
	}
	return resources, nil

}
func (ctx *RenderContext) Render() error {
	before := func(processor Processor, context *RenderContext, resources []*data.Resource) error {
		err := processor.Before(context, resources)
		if err != nil {
			return errors.Wrap(err, "Before execution phase is failed");
		}
		return nil
	}
	after := func(processor Processor, context *RenderContext, resources []*data.Resource) error {
		err := processor.After(context, resources)
		if err != nil {
			return errors.Wrap(err, "After execution phase is failed");
		}
		return nil
	}
	process := func(processor Processor, context *RenderContext, resources []*data.Resource) error {
		for _, resource := range resources {
			err := processor.BeforeResource(resource)
			if err != nil {
				return errors.Wrap(err, "Applyinig transformation BeforeResource "+resource.Filename+" is failed")
			}
			resource.Content.Accept(processor)
			err = processor.AfterResource(resource)
			if err != nil {
				return errors.Wrap(err, "Applyinig transformation AfterResource "+resource.Filename+" is failed")
			}

		}
		return nil
	}
	_, err := ctx.RootResource.Execute(ctx, before)
	if err != nil {
		return err
	}
	_, err = ctx.RootResource.Execute(ctx, process)
	if err != nil {
		return err
	}
	_, err = ctx.RootResource.Execute(ctx, after)
	if err != nil {
		return err
	}
	return nil
}

//parse the directory structure and the flekszible configs from the dirs
func (node *ResourceNode) LoadResourceConfig(sourceCache *data.SourceCacheManager) error {
	conf, configFilePath, err := data.ReadConfiguration(node.Dir)
	if err != nil {
		return errors.Wrap(err, "Can't parse flekszible.yaml/Flekszible descriptor from  "+node.Dir)
	}
	if path.Base(configFilePath) != "Flekszible" {
		node.Standalone = false
	} else {
		node.Standalone = true
	}
	if !node.Standalone {
		node.Resources = data.ReadResourcesFromDir(node.Dir)
	}
	node.Source = make([]data.Source, 0)
	for _, definedSource := range conf.Source {
		if definedSource.Url != "" {
			node.Source = append(node.Source, &data.GoGetter{Url: definedSource.Url})
		} else if definedSource.Path != "" {
			node.Source = append(node.Source, &data.LocalSource{BaseDir: node.Dir, RelativeDir: definedSource.Path})
		}
	}
	//update destinations of the direct k8s resources
	for ix, _ := range node.Resources {
		node.Resources[ix].Destination = node.Destination
	}
	if len(conf.Transformations) > 0 {
		bytes, err := yaml.Marshal(conf.Transformations)
		if err != nil {
			return err
		}
		node.PreImportTransformations = bytes
	}
	for _, importDefinition := range conf.Import {
		source, err := locate(node.Dir, importDefinition.Path, node.Source, sourceCache)
		if err != nil {
			return err
		}
		dir, _ := source.GetPath(sourceCache, importDefinition.Path)
		childNode := CreateResourceNode(dir, importDefinition.Destination, source)
		if importDefinition.Destination == "" {
			childNode.Destination = node.Destination
		}
		childNode.Origin = source
		if node.Dir == childNode.Dir {
			panic("Recursive directory parser " + node.Dir + " loads" + childNode.Dir)
		}
		err = childNode.LoadResourceConfig(sourceCache)
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
	if !node.Standalone {
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
	}
	for _, child := range node.Children {
		child.LoadDefinitions()
	}

}

//try to find the first Source which contains the file
func locate(basedir string, dir string, sources []data.Source, cacheManager *data.SourceCacheManager) (data.Source, error) {
	allSources := make([]data.Source, 0)
	allSources = append(allSources, &data.EnvSource{})
	allSources = append(allSources, &data.CurrentDir{CurrentDir: basedir})
	allSources = append(allSources, sources...)

	for _, source := range allSources {
		resourcePath, err := source.GetPath(cacheManager, dir)
		if err != nil {
			tpe, value := source.ToString()
			logrus.Error("Can't check dir from the source " + tpe + "/" + value + err.Error())
		} else if resourcePath != "" {
			if _, err := os.Stat(resourcePath); !os.IsNotExist(err) {
				return source, nil
			}
		}
	}
	return nil, errors.New("Couldn't find dir: " + dir)
}
