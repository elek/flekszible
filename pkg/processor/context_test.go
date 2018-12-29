package processor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	context := CreateRenderContext("k8s", "../../testdata/readconfigs", "../../testdata/readconfigs")
	err := context.Init()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(context.RootResource.Resources))
	assert.Equal(t, 2, len(context.RootResource.Children))
	assert.Equal(t, 2, len(context.Resources()))
}
