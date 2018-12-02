package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/sirupsen/logrus"
)

func Generate(processorRepository ProcessorRepository, ctx *data.RenderContext) {
	logrus.Info("Active processors:")
	for _, proc := range processorRepository.Processors {
		logrus.Infof("%T %s", proc, proc)
	}
	ctx.ReadResources()

	for _, processor := range processorRepository.Processors {
		processor.Before(ctx)
	}

	for _, resource := range ctx.Resources {
		logrus.Infof("Processing %s (%s) from %s", resource.Name(), resource.Kind(), resource.Filename)
		for _, processor := range processorRepository.Processors {
			processor.BeforeResource(&resource)
			resource.Content.Accept(processor)
			processor.AfterResource(&resource)
		}
	}

	for _, processor := range processorRepository.Processors {
		processor.After(ctx)
	}

}


