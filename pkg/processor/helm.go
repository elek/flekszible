package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/helm"
	"github.com/elek/flekszible/pkg/yaml"
	"io/ioutil"
	"path"
	"strings"
)

type HelmDecorator struct {
	DefaultProcessor
	Values    helm.Values
	ChartName string
	HostNames []string
}

func (processor *HelmDecorator) After(ctx *RenderContext, resources []data.Resource) {
	data, err := yaml.Marshal(processor.Values)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path.Join(ctx.OutputDir, "values.yaml"), data, 0644)
	if err != nil {
		panic(err)
	}
}

func (p *HelmDecorator) Before(ctx *RenderContext, resources []data.Resource) {
	p.HostNames = make([]string, 0)
	for _, resource := range resources {
		kind := resource.Kind()
		if kind == "StatefulSet" {
			name := resource.Name()

			serviceNameGet := data.Get{Path: data.NewPath("spec", "serviceName")}
			resource.Content.Accept(&serviceNameGet)
			if serviceNameGet.Found {
				hostname := name + "-0." + serviceNameGet.ValueAsString()
				p.HostNames = append(p.HostNames, hostname)
			}
		}
	}
}

func (p *HelmDecorator) BeforeResource(resource *data.Resource) {

	content := resource.Content
	prefix := "{{ template \"ozone.fullname\" . }}-"
	changeImage := func(original interface{}) interface{} {
		imageString := original.(string)
		if !strings.Contains(imageString, ":") {
			imageString = imageString + ":latest"
		}
		imageAndTag := strings.Split(imageString, ":")
		if p.Values.Image.Repository == "" {
			p.Values.Image.Repository = imageAndTag[0]
			p.Values.Image.Tag = imageAndTag[1]
			p.Values.Image.PullPolicy = "Always"
		}
		return "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
	}

	prefixName := func(original interface{}) interface{} {
		return prefix + original.(string)
	}

	prefixHostName := func(original interface{}) interface{} {
		result := original.(string)
		for _, hostName := range p.HostNames {
			splitted := strings.Split(hostName, ".")
			result = strings.Replace(result, hostName, prefix+splitted[0]+"."+prefix+splitted[1], -1)
		}
		return result
	}

	content.Accept(&data.Apply{Path: data.NewPath("metadata", "name"), Function: prefixName})
	content.Accept(&data.Apply{Path: data.NewPath("spec", "serviceName"), Function: prefixName})
	content.Accept(&data.Apply{Path: data.NewPath("spec", "template", "spec", ".*ontainers", ".*", "env", ".*", "value"), Function: prefixHostName})

	content.Accept(&data.Apply{Path: data.NewPath("data", ".*"), Function: prefixHostName})
	content.Accept(&data.Apply{Path: data.NewPath("spec", "template", "spec", ".*ontainers", ".*", "image"), Function: changeImage})
	content.Accept(&data.Apply{Path: data.NewPath("spec", "template", "spec", ".*ontainers", ".*", "envFrom", ".*", "configMapRef", "name", ), Function: prefixName})


	labelsGetter := data.Get{Path: data.NewPath("metadata", "labels")}
	content.Accept(&labelsGetter)
	labelsGetter.ReturnValue.(*data.MapNode).PutValue("apps.kubernetes.io/name", "{{ include \"ozone.name\" . }}")
	labelsGetter.ReturnValue.(*data.MapNode).PutValue("helm.sh/chart", "{{ include \"ozone.chart\" . }}")
	labelsGetter.ReturnValue.(*data.MapNode).PutValue("app.kubernetes.io/instance", "{{ .Release.Name }}")
	labelsGetter.ReturnValue.(*data.MapNode).PutValue("app.kubernetes.io/managed-by", "{{ .Release.Service }}")
}


func init() {
	prototype := HelmDecorator{}
	ProcessorTypeRegistry.Add(&prototype)
}