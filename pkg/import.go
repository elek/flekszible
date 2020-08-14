package pkg

import (
	"github.com/elek/flekszible/api/data"
	"github.com/elek/flekszible/api/processor"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

func Import(resourceFile string, transformations []string, outputDir string) error {
	context := &processor.RenderContext{
		OutputDir: outputDir,
		Mode:      "k8s",
	}

	root := &processor.ResourceNode{
		Dir:                 "stdout",
		Destination:         ".",
		Resources:           make([]*data.Resource, 0),
		Children:            make([]*processor.ResourceNode, 0),
		ProcessorRepository: processor.CreateProcessorRepository(),
	}

	context.RootResource = root

	err := context.AddAdHocTransformations(transformations)
	if err != nil {
		return err
	}
	var bytesOfResources []byte
	if resourceFile != "" {
		bytesOfResources, err = ioutil.ReadFile(resourceFile)
		if err != nil {
			return errors.Wrap(err, "Can't open the resource file defined by the argument "+resourceFile)
		}
	} else {
		bytesOfResources, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return errors.Wrap(err, "Stdin can't be read")
		}
	}

	resources, err := data.LoadResourceFromByte(bytesOfResources)
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
