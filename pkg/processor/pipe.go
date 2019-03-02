package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/yaml"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

type Pipe struct {
	DefaultProcessor
	Command string
	Trigger Trigger
	Args    []string
}

func (p *Pipe) BeforeResource(resource *data.Resource) {
	if !p.Trigger.active(resource) {
		return
	}
	converted := data.ConvertToYaml(resource.Content)
	bytes, err := yaml.Marshal(converted)
	if err != nil {
		logrus.Error("Can't parse yaml file of " + resource.Path + err.Error())
		return
	}
	str := string(bytes)

	cmd := exec.Command(p.Command, p.Args...)

	builder := strings.Builder{}
	cmd.Stdin = strings.NewReader(str)
	cmd.Stdout = &builder
	cmd.Stderr = os.Stdout
	err = cmd.Run()
	if err != nil {
		logrus.Error("Can't execute the command " + p.Command + " " + err.Error())
		return
	}
	output := builder.String()
	logrus.Infof("Command is executed with the following output %s", output)
	node, err := data.ReadString([]byte(output ))
	if err != nil {
		logrus.Error("Can't parse the result of the piped command for " + resource.Path + " " + err.Error())
		return
	}
	resource.Content = node
}

func init() {
	ProcessorTypeRegistry.Add(ProcessorDefinition{
		Metadata: ProcessorMetadata{
			Name:        "Pipe",
			Description: "Transform content with external shell command.",
			Parameter: []ProcessorParameter{
				{
					Name:        "command",
					Description: "External program which transforms standard input to output",
				},
				{
					Name:        "args",
					Description: "List of the arguments of the command",
				},
				TriggerParameter,
			},
			Doc: `
Pipe executes a specific command to transform a k8s resources to a new resources.

The original manifest will be sent to the stdin of the process and the stdout will be processed as a the stdout of the file.

Example:

'''
- type: Pipe
  command: sed
  args: ["s/nginx/qwe/g"]
'''
` + TriggerDoc,
		},
		Factory: func(config *yaml.MapSlice) (Processor, error) {
			return configureProcessorFromYamlFragment(&Pipe{}, config)
		},
	})
}
