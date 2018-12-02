package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadResources(t *testing.T) {
	ctx := data.RenderContext{
		InputDir: []string{"../../testdata/context/base3"},
	}

	ctx.ReadConfigs()
	ctx.ReadResources()
	assert.Equal(t, 2, len(ctx.Resources))
	assert.Equal(t, "datanode", ctx.Resources[0].Name())
}
