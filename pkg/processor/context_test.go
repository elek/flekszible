package processor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfigs(t *testing.T) {
	context := CreateRenderContext("k8s", "../../testdata/readconfigs", "../../testdata/readconfigs")
	res, err := context.ReadConfigs()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(res))

	_, ok := res["../../testdata/readconfigs/i1"]
	assert.True(t, ok)

	_, ok = res["../../testdata/readconfigs/x1"]
	assert.True(t, ok)

}

func TestReadResources(t *testing.T) {
	ctx := RenderContext{
		InputDir: []string{"../../testdata/readresources"},
	}

	_, err := ctx.ReadConfigs()
	assert.Nil(t, err)
	ctx.ReadResources()
	assert.Equal(t, 2, len(ctx.Resources))
	assert.Equal(t, "s3g", ctx.Resources[0].Name())
}
