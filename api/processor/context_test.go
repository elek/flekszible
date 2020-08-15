package processor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInit(t *testing.T) {
	context := CreateRenderContext("k8s", "../../testdata/readconfigs", "../../testdata/readconfigs/out")
	err := context.Init()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(context.RootResource.Resources))
	assert.Equal(t, 2, len(context.RootResource.Children))
	assert.Equal(t, 2, len(context.Resources()))
}

func TestCreateTransformation(t *testing.T) {

	proc, err := parseTransformation("image:image=asd")
	assert.Nil(t, err)

	image := proc.(*Image)

	assert.Equal(t, "asd", image.Image)
}

func TestCreateTransformationImageWithColon(t *testing.T) {

	proc, err := parseTransformation("image:image=localhost:5000/image/name:tag")
	assert.Nil(t, err)

	image := proc.(*Image)

	assert.Equal(t, "localhost:5000/image/name:tag", image.Image)
}

func TestCreateTransformationImageWithQuotedComa(t *testing.T) {

	proc, err := parseTransformation("image:image=test\\,sg")
	assert.Nil(t, err)

	image := proc.(*Image)

	assert.Equal(t, "test,sg", image.Image)
}
