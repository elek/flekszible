package pkg

import (
	"github.com/elek/flekszible/api/processor"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateTransformation(t *testing.T) {

	proc, err := parseTransformation("image:image=asd")
	assert.Nil(t, err)

	image := proc.(*processor.Image)

	assert.Equal(t, "asd", image.Image)
}

func TestCreateTransformationImageWithColon(t *testing.T) {

	proc, err := parseTransformation("image:image=localhost:5000/image/name:tag")
	assert.Nil(t, err)

	image := proc.(*processor.Image)

	assert.Equal(t, "localhost:5000/image/name:tag", image.Image)
}

func TestCreateTransformationImageWithQuotedComa(t *testing.T) {

	proc, err := parseTransformation("image:image=test\\,sg")
	assert.Nil(t, err)

	image := proc.(*processor.Image)

	assert.Equal(t, "test,sg", image.Image)
}
