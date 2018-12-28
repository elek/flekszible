package processor

import (
	"path"
)

type HelmWriter struct {
	K8sWriter
	OutputDir string
}

func (writer *HelmWriter) Before(ctx *RenderContext) {
	writer.resourceOutputDir = path.Join(ctx.OutputDir, "templates")
}

func init() {
	prototype := HelmWriter{}
	ProcessorTypeRegistry.Add(&prototype)
}