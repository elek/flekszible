package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"path"
)

type HelmWriter struct {
	K8sWriter
	OutputDir string
}

func (writer *HelmWriter) Before(ctx *data.RenderContext) {
	writer.resourceOutputDir = path.Join(ctx.OutputDir, "templates")
}

func init() {
	prototype := HelmWriter{}
	ProcessorTypeRegistry.Add(&prototype)
}