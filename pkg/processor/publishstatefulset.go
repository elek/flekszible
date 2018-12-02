package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"strings"
)

type PublishStatefulSet struct {
	DefaultProcessor
	Trigger Trigger
}

func (processor *PublishStatefulSet) Before(ctx *data.RenderContext) {
	newResources := make([]data.Resource, 0)
	for _, resource := range ctx.Resources {
		if processor.Trigger.active(&resource) && resource.Kind() == "Service" && hasNoneClusterIp(resource.Content) {
			newContent := DeepCopy(resource.Content)

			metadata := newContent.Get("metadata").(*data.MapNode)
			metadata.PutValue("name", metadata.GetStringValue("name")+"-public")
			spec := newContent.Get("spec").(*data.MapNode)
			spec.Remove("clusterIP")
			spec.PutValue("type", "NodePort")
			r := data.Resource{
				Content: newContent,
			}
			newResources = append(newResources, r)
		}

	}
	ctx.Resources = append(ctx.Resources, newResources...)
}
func hasNoneClusterIp(slice *data.MapNode) bool {
	return true
}

func init() {
	prototype := PublishStatefulSet{}
	ProcessorTypeRegistry.Add(&prototype)
}

func DeepCopy(src data.Node) *data.MapNode {

	buffer := strings.Builder{}
	writer := K8sWriter{
		output: &buffer,
	}
	src.Accept(&writer)

	mapNode, err := data.ReadString([]byte(buffer.String()))
	if err != nil {
		panic(err)
	}
	return mapNode
}
