package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/elek/flekszible/pkg/helm"
	"github.com/elek/flekszible/pkg/yaml"
	"io/ioutil"
	"path"
)

type HelmDecorator struct {
	DefaultProcessor
	Values helm.Values
}

func (processor *HelmDecorator) After(ctx *data.RenderContext) {
	data, err := yaml.Marshal(processor.Values)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path.Join(ctx.OutputDir, "values.yaml"), data, 0644)
	if err != nil {
		panic(err)
	}
}

func (p *HelmDecorator) BeforeResource(resource *data.Resource) {

	content := resource.Content

	content.Accept(&data.Apply{Path: data.NewPath("metadata", "name"), Function: chageName})
	content.Accept(&data.Apply{Path: data.NewPath("spec", "template", "spec", "containers", "*", "image"), Function: changeImage})
	content.Accept(&data.Apply{Path: data.NewPath("spec", "template", "spec", "initContainers", "*", "image"), Function: changeImage})

	labelsGetter := data.Get{Path: data.NewPath("metadata", "labels")}
	content.Accept(&labelsGetter)
	labelsGetter.ReturnValue.(*data.MapNode).PutValue("apps.kubernetes.io/name", "{{ include \"ozone.name\" . }}")
	labelsGetter.ReturnValue.(*data.MapNode).PutValue("helm.sh/chart", "{{ include \"ozone.chart\" . }}")
	labelsGetter.ReturnValue.(*data.MapNode).PutValue("app.kubernetes.io/instance", "{{ .Release.Name }}")
	labelsGetter.ReturnValue.(*data.MapNode).PutValue("app.kubernetes.io/managed-by", "{{ .Release.Service }}")
}

func chageName(name interface{}) interface{} {
	return "{{ template \"fullname\" . }}-" + name.(string)
}

func changeImage(name interface{}) interface{} {
	return "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
}
//
//func (processor *HelmDecorator) ProcessKey(path data.Path, value interface{}) interface{} {
//
//	if path.MatchSegments("spec", "template", "spec", "containers", "*", "envFrom", "*", "configMapRef", "name") {
//		return "\"{{ template \"fullname\" . }}-" + value.(string) + "\""
//	}
//	if path.MatchSegments("spec", "template", "spec", "containers", "*", "image") ||
//		path.MatchSegments("spec", "template", "spec", "initContainers", "*", "image") {
//		image := value.(string)
//		parts := strings.Split(image, ":")
//		processor.Values.Image.Repository = parts[0]
//		processor.Values.Image.PullPolicy = "IfNotPresent"
//		if len(parts) > 1 {
//			processor.Values.Image.Tag = parts[1]
//		} else {
//			processor.Values.Image.Tag = "latest"
//		}
//		return "\"{{ .Values.image.repository }}:{{ .Values.image.tag }}\""
//
//	}
//	return value
//}
//
//func (*HelmDecorator) BeforeMap(path data.Path, object yaml.MapSlice) yaml.MapSlice {
//	if path.MatchSegments("metadata", "labels") {
//		object = object.Put("app", "{{ template \"fullname\" . }}")
//	}
//	if path.MatchSegments("spec", "template", "spec", "containers", "*") ||
//		path.MatchSegments("spec", "template", "spec", "initContainers", "*") {
//		if _, ok := object.Get("imagePullPolicy"); !ok {
//			object = object.Put("imagePullPolicy", "{{ .Values.image.pullPolicy}}")
//		}
//	}
//	if path.MatchSegments("spec", "template", "spec", "volumes", "*") {
//		if config, ok := object.Get("configMap"); ok {
//			switch configMap := config.(type) {
//			case yaml.MapSlice:
//				if name, ok := configMap.Get("name"); ok {
//					configMap = configMap.Put("name", "\"{{ template \"fullname\" . }}-"+name.(string)+"\"")
//					object = object.Put("configMap", configMap)
//				}
//
//			default:
//				panic("Configmap is not a map")
//			}
//
//		}
//	}
//
//	return object
//}

func init() {
	prototype := HelmDecorator{}
	ProcessorTypeRegistry.Add(&prototype)
}