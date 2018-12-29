package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
)

type K8sWriter struct {
	DefaultProcessor
	arrayParent       bool
	mapIndex          int
	started           bool
	resourceOutputDir string
	output            io.Writer
	file              *os.File
}

func (writer *K8sWriter) Before(ctx *RenderContext, resources []data.Resource) {
	writer.resourceOutputDir = ctx.OutputDir
}
func (writer *K8sWriter) createOutputPath(outputDir, name, kind string) string {
	fileName := name + "-" + strings.ToLower(kind) + ".yaml"
	return path.Join(outputDir, fileName)
}

func (writer *K8sWriter) BeforeResource(resource *data.Resource) {
	writer.started = false
	outputDir := writer.resourceOutputDir
	if outputDir == "-" {
		writer.output = os.Stderr
	} else {

		outputFile := writer.createOutputPath(outputDir, resource.Name(), resource.Kind())
		os.MkdirAll(outputDir, os.ModePerm)
		output, err := os.Create(outputFile)
		if err != nil {
			panic(err)
		}
		writer.output = output
		writer.file = output
	}
}

func (writer *K8sWriter) AfterResource(*data.Resource) {
	if writer.file != nil {
		logrus.Infof("Result is written to %s", writer.file.Name())
		writer.file.Close()
	} else {
		write(writer.output, "\n---\n")
	}
}

func (writer *K8sWriter) OnKey(node *data.KeyNode) {
	printKey(writer.output, node.Path, node.Value)
	writer.arrayParent = false
}
func (writer *K8sWriter) BeforeMap(node *data.MapNode) {
	if node.Len() == 0 {
		write(writer.output, " {}")
	}
	if !writer.arrayParent && writer.started {
		write(writer.output, "\n")
	}
	writer.started = true
}

func (writer *K8sWriter) BeforeMapItem(node *data.MapNode, key string, index int) {
	if writer.arrayParent && index == 0 {
		identedPrint(writer.output, 0, " "+key+":")
		writer.arrayParent = false
	} else {
		identedPrint(writer.output, node.Path.Length(), key+":")
	}
}

func (writer *K8sWriter) BeforeList(node *data.ListNode) {
	if node.Len() == 0 {
		write(writer.output, " []\n")
	} else {
		write(writer.output, "\n")
		writer.arrayParent = true
	}
}

func (writer *K8sWriter) BeforeListItem(node *data.ListNode, item data.Node, index int) {
	identedPrint(writer.output, node.Path.Length(), "-")
	writer.arrayParent = true
}

func CreateStdK8sWriter() *K8sWriter {
	writer := K8sWriter{
		resourceOutputDir: "-",
	}
	return &writer
}

func identedPrint(output io.Writer, ident int, s string) {
	for a := 0; a < ident*2; a++ {
		write(output, " ")
	}
	write(output, s)
}


func init() {
	prototype := K8sWriter{}
	ProcessorTypeRegistry.Add(&prototype)
}

func printKey(output io.Writer, path data.Path, rawValue interface{}) {
	value := ""
	switch converted := rawValue.(type) {
	case string:
		value = converted
	case bool:
		value = strconv.FormatBool(converted)
	case int:
		value = strconv.Itoa(converted)
	}
	if strings.Contains(value, "\n") {
		write(output, " |-\n")
		for _, line := range strings.Split(value, "\n") {
			identedPrint(output, path.Length()+1, line+"\n")
		}
	} else if len(value) == 0 {
		write(output, " ")
		write(output, "\"\"")
		write(output, "\n")
	} else if path.Length() > 2 && path.Segment(-2) == "annotations" || path.Segment(-3) == "env" || path.Segment(-2) == "data" {
		write(output, " ")
		//escaped := strings.Replace(value, "\"", "\\\"", -1)
		write(output, "\""+value+"\"")
		write(output, "\n")
	} else {
		write(output, " ")
		write(output, value)
		write(output, "\n")
	}


}

func write(w io.Writer, content string) (int, error) {
	return w.Write([]byte(content))
}
