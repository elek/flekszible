package processor

import (
	"github.com/elek/flekszible/api/data"
	"github.com/elek/flekszible/api/yaml"
	"github.com/pkg/errors"
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
	Name   string      //the original import name used in FLekszible files
	Dir    string      //absolut path of the current dir
	Origin data.Source //source used to load the resource

	Destination              string
	Resources                []*data.Resource
	Children                 []*ResourceNode
	PreImportTransformations []byte
	Source                   []data.Source
	ProcessorRepository      *TransformationRepository
	ResourcesDir             string
	Definitions              []string //definitions loaded from this dir
}

type ResourceLocation struct {
	Name   string      //the name to import the location (like 'ozone')
	Source data.Source //the source used to load the resource
	Dir    string      //absolute path of the location
}

func CreateRenderContext(mode string, inputDir string, outputDir string) *RenderContext {
	return &RenderContext{
		OutputDir:    outputDir,
		Mode:         mode,
		RootResource: CreateResourceNode(inputDir, "", &data.LocalSource{Dir: inputDir}),
	}
}

func (context *RenderContext) ListResourceNodes() []*ResourceNode {
	return listResourceNodesInt(context.RootResource)
}

func (context *RenderContext) AddAdHocTransformations(transformations []string) error {

	proc, err := createTransformation(transformations)
	if err != nil {
		return err
	}
	context.RootResource.ProcessorRepository.AppendAll(proc)
	return nil
}

func parseTransformation(trafoDef string) (Processor, error) {
	parts := strings.SplitN(trafoDef, ":", 2)
	name := parts[0]
	parameterMap := make(map[string]string)
	if len(parts) > 1 {
		transofmationsString := strings.ReplaceAll(parts[1], "\\,", "__NON_SEPARATOR_COMA__")
		for _, rawParam := range strings.Split(transofmationsString, ",") {
			parameter := strings.ReplaceAll(rawParam, "__NON_SEPARATOR_COMA__", ",")
			paramParts := strings.SplitN(parameter, "=", 2)
			if len(paramParts) < 2 {
				return nil, errors.New("Parameters should be defined in the form key=value and not like " + parameter)
			}
			parameterMap[paramParts[0]] = paramParts[1]
		}
	}
	proc, err := ProcessorTypeRegistry.Create(name, parameterMap)
	if err != nil {
		return nil, errors.Wrap(err, "Can't create transformation based on the string "+trafoDef)
	}
	return proc, nil
}

//parse one-liner transformation definition
func createTransformation(transformationsDefinitions []string) ([]Processor, error) {
	result := make([]Processor, 0)
	for _, trafoDef := range strings.Split(os.Getenv("FLEKSZIBLE_TRANSFORMATION"), ";") {
		if len(strings.TrimSpace(trafoDef)) > 0 {
			transformation, err := parseTransformation(trafoDef)
			if err != nil {
				return result, errors.Wrap(err, "Can't parse transformation defined by FLEKSZIBLE_TRANSFORMATION: "+trafoDef)
			}
			result = append(result, transformation)
		}
	}
	for _, transformationsDefinition := range transformationsDefinitions {
		transformation, err := parseTransformation(transformationsDefinition)
		if err != nil {
			return result, errors.Wrap(err, "Can't parse transformation defined by cli arg: "+transformationsDefinition)
		}
		result = append(result, transformation)
	}
	return result, nil
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
	data.Generators = append(data.Generators, &data.OutputGenerator{})
	data.Generators = append(data.Generators, &data.ConfigGenerator{})
	data.Generators = append(data.Generators, &data.SecretGenerator{})
	cacheManager := data.NewSourceCacheManager(context.RootResource.Dir)
	return context.RootResource.LoadResourceConfig(&cacheManager, context.OutputDir)
}

//List all the resources from the resource tree.
func (context *RenderContext) Resources() []*data.Resource {
	result := make([]*data.Resource, 0)
	result = append(result, context.RootResource.AllResources()...)
	return result
}

func (node *ResourceNode) AllResources() []*data.Resource {
	result := make([]*data.Resource, 0)
	for _, child := range node.Children {
		result = append(result, child.AllResources()...)
	}
	result = append(result, node.Resources...)
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
			return errors.Wrap(err, "Couldn't parse transformations from the dir "+node.Dir)
		}
		node.ProcessorRepository.AppendAll(processors)
	}

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

type execute func(transformation Processor, context *RenderContext, node *ResourceNode) error

func (node *ResourceNode) Execute(context *RenderContext, functionCall execute) error {
	//execute on child
	for _, child := range node.Children {
		err := child.Execute(context, functionCall)
		if err != nil {
			return errors.Wrap(err, "The function can't be executed on one of the childResources")
		}
	}

	//execute on this resource
	for _, transformation := range node.ProcessorRepository.Processors {
		err := functionCall(transformation, context, node)
		if err != nil {
			return errors.Wrap(err, "The function can't be executed on all the resources. Transformation: ")
		}
	}
	return nil

}

//execute all the transformations in the rendercontext.
func (ctx *RenderContext) Render() error {

	before := func(processor Processor, context *RenderContext, node *ResourceNode) error {
		err := processor.Before(context, node)
		if err != nil {
			return errors.Wrap(err, "Before execution phase is failed")
		}
		return nil
	}
	after := func(processor Processor, context *RenderContext, node *ResourceNode) error {
		err := processor.After(context, node)
		if err != nil {
			return errors.Wrap(err, "After execution phase is failed")
		}
		return nil
	}
	process := func(processor Processor, context *RenderContext, node *ResourceNode) error {
		for _, resource := range node.AllResources() {
			err := processor.BeforeResource(resource)
			if err != nil {
				return errors.Wrap(err, "Applying transformation BeforeResource "+resource.Filename+" is failed")
			}
			resource.Content.Accept(processor)
			err = processor.AfterResource(resource)
			if err != nil {
				return errors.Wrap(err, "Applying transformation AfterResource "+resource.Filename+" is failed")
			}

		}
		return nil
	}
	err := ctx.RootResource.Execute(ctx, func(processor Processor, context *RenderContext, node *ResourceNode) error {
		err := processor.RegisterResources(context, node)
		if err != nil {
			return errors.Wrap(err, "RegisterResource execution phase is failed")
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = ctx.RootResource.Execute(ctx, before)
	if err != nil {
		return err
	}
	err = ctx.RootResource.Execute(ctx, process)
	if err != nil {
		return err
	}
	err = ctx.RootResource.Execute(ctx, after)
	if err != nil {
		return err
	}
	return nil
}

//parse the directory structure and the flekszible configs from the dirs
func (node *ResourceNode) LoadResourceConfig(sourceCache *data.SourceCacheManager, outputDir string) error {
	conf, _, err := data.ReadConfiguration(node.Dir)
	if err != nil {
		return errors.Wrap(err, "Can't parse flekszible.yaml/Flekszible descriptor from  "+node.Dir)
	}

	node.Resources = data.ReadResourcesFromDir(path.Join(node.Dir, conf.ResourcesDir))

	for _, generator := range data.Generators {
		dirs, err := ioutil.ReadDir(node.Dir)
		if err == nil {
			for _, dir := range dirs {
				managedDir := path.Join(node.Dir, dir.Name())
				if dir.IsDir() && generator.IsManagedDir(managedDir) {
					resources, err := generator.Generate(managedDir, outputDir)
					if err != nil {
						return errors.Wrap(err, "Can't generate resources from the dir "+managedDir)
					}
					node.Resources = append(node.Resources, resources...)

				}
			}
		}
	}

	node.Source = make([]data.Source, 0)
	for _, definedSource := range conf.Source {
		if definedSource.Url != "" {
			node.Source = append(node.Source, &data.RemoteSource{Url: definedSource.Url})
		} else if definedSource.Path != "" {
			node.Source = append(node.Source, &data.LocalSource{Dir: path.Join(node.Dir, definedSource.Path)})
		}
	}
	//update destinations of the direct k8s resources
	for ix, _ := range node.Resources {
		node.Resources[ix].Destination = node.Destination
	}

	//read inline transformations
	if len(conf.Transformations) > 0 {
		bytes, err := yaml.Marshal(conf.Transformations)
		if err != nil {
			return err
		}
		node.PreImportTransformations = bytes
	}

	//import _global directories
	for _, source := range node.Source {
		sourceDir, err := source.GetPath(sourceCache)
		if err != nil {
			errors.Wrap(err, "Source definition defined in "+node.Dir+" couldn't be loaded")
		}
		globalDirs := []string{path.Join(sourceDir, "flekszible", "_global"), path.Join(sourceDir, "_global")}
		for _, globalDir := range globalDirs {
			if stat, err := os.Stat(globalDir); err == nil && stat.IsDir() {
				childNode := CreateResourceNode(globalDir, node.Destination, source)
				childNode.Name = "_global"
				childNode.Origin = source
				err = childNode.LoadResourceConfig(sourceCache, outputDir)
				if err != nil {
					return errors.Wrap(err, "Couldn't load _global dir from "+globalDir)
				}
				node.Children = append(node.Children, childNode)
				break
			}
		}
	}

	//read imported directories
	for _, importDefinition := range conf.Import {
		source, err := locate(node.Dir, importDefinition.Path, node.Source, sourceCache)
		if err != nil {
			return err
		}
		if source == nil {
			return errors.New("Directory dependency `" + importDefinition.Path + "` defined in " + node.Dir + " can't be found")
		}
		sourceDir, _ := source.GetPath(sourceCache)
		childNode := CreateResourceNode(checkPath(sourceDir, importDefinition.Path), importDefinition.Destination, source)
		childNode.Name = importDefinition.Path
		if importDefinition.Destination == "" {
			childNode.Destination = node.Destination
		}
		childNode.Origin = source
		if node.Dir == childNode.Dir {
			panic("Recursive directory parser " + node.Dir + " loads" + childNode.Dir)
		}
		err = childNode.LoadResourceConfig(sourceCache, outputDir)
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
				name, err := parseDefintion(definitionFile)
				if err != nil {
					logrus.Errorf("Can't parse the definition file %s: %s", definitionFile, err.Error())
				}
				if node.Definitions == nil {
					node.Definitions = make([]string, 0)
				}
				node.Definitions = append(node.Definitions, name)
			}
		}

	for _, child := range node.Children {
		child.LoadDefinitions()
	}

}

//try to find the first Source which contains the dir
func locate(basedir string, dir string, sources []data.Source, cacheManager *data.SourceCacheManager) (data.Source, error) {
	allSources := make([]data.Source, 0)
	allSources = append(allSources, data.LocalSourcesFromEnv()...)
	allSources = append(allSources, &data.LocalSource{Dir: basedir})
	allSources = append(allSources, sources...)

	for _, source := range allSources {
		resourcePath, err := source.GetPath(cacheManager)
		if err != nil {
			value := source.ToString()
			logrus.Error("Can't check dir from the source " + value + err.Error())
		} else if resourcePath != "" {
			path := checkPath(resourcePath, dir)
			if path != "" {
				return source, nil
			}
		}
	}
	return nil, errors.New("Couldn't find dir: " + dir)
}

func checkPath(baseDir string, subdir string) string {
	optionalSubDir := path.Join(baseDir, "flekszible")
	if stat, err := os.Stat(optionalSubDir); !os.IsNotExist(err) && stat.IsDir() {
		result := path.Join(baseDir, "flekszible", subdir)
		if _, err := os.Stat(result); !os.IsNotExist(err) {
			return result
		}
	}
	result := path.Join(baseDir, subdir)
	if _, err := os.Stat(result); !os.IsNotExist(err) {
		return result
	}
	return ""
}