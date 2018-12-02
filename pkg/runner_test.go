package pkg

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfigs(t *testing.T) {
	ctx := data.RenderContext{
		InputDir: []string{"../testdata/context/base3"},
	}

	ctx.ReadConfigs()

	assert.Equal(t, []string{"../testdata/context/base1", "../testdata/context/base2", "../testdata/context/base3"}, ctx.InputDir)
}
