package processor

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/elek/flekszible/api/v2/data"
	"github.com/elek/flekszible/api/v2/yaml"
	"github.com/pkg/errors"
)

type K8sWriter struct {
	DefaultProcessor
	arrayParent       bool
	mapIndex          int
	started           bool
	resourceOutputDir string
	output            io.Writer
	file              *os.File
	header            string
}

func (writer *K8sWriter) ToString() string {
	return "k8swriter"
}

func (writer *K8sWriter) Before(ctx *RenderContext, node *ResourceNode) error {
	writer.resourceOutputDir = ctx.OutputDir
	return nil
}

func CreateOutputFileName(name string, kind string) string {
	return strings.ToLower(name) + "-" + strings.ToLower(kind) + ".yaml"
}
func (writer *K8sWriter) createOutputPath(outputDir, name, kind string, destination string, destinationFile string) string {
	if destinationFile != "" {
		return path.Join(outputDir, destination, destinationFile)
	} else {
		return path.Join(outputDir, destination, CreateOutputFileName(name, kind))
	}
}

func (writer *K8sWriter) BeforeResource(resource *data.Resource) error {
	if exclude, ok := resource.Metadata["exclude"]; ok {
		if exclude == "true" {
			return nil
		}
	}
	if resource.Kind() == "Kustomization" {
		return nil
	}
	writer.started = false
	outputDir := writer.resourceOutputDir
	if outputDir == "-" {
		content, err := resource.Content.ToString()
		if err != nil {
			return errors.Wrap(err, "Can't render the content of a transformed resource file")
		}
		fmt.Println(content)
		fmt.Println("---")
	} else {
		licenceHeader := ""
		licenceHeaderFile := path.Join(outputDir, "LICENSE.header")
		if _, err := os.Stat(licenceHeader); os.IsNotExist(err) {
			content, _ := os.ReadFile(licenceHeaderFile)
			licenceHeader = string(content) + "\n"
		}
		outputFile := writer.createOutputPath(outputDir, resource.Name(), resource.Kind(), resource.Destination, resource.DestinationFileName)
		err := os.MkdirAll(path.Dir(outputFile), os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "Can't create the output directory "+path.Dir(outputFile))
		}

		content, err := resource.Content.ToString()
		if err != nil {
			return errors.Wrap(err, "Can't render the content of a transformed resource file")
		}

		err = os.WriteFile(outputFile, []byte(licenceHeader+content), 0655)
		if err != nil {
			return errors.Wrap(err, "Can't write the k8s file out "+outputFile)
		}
	}
	return nil
}

func CreateStdK8sWriter() *K8sWriter {
	writer := K8sWriter{
		resourceOutputDir: "-",
	}
	return &writer
}

func ActivateK8sWriter(registry *ProcessorTypes) {
	registry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "K8sWriter",
			Description: "Internal transformation to print out k8s resources as yaml",
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&K8sWriter{}, config)
		},
	})
}
