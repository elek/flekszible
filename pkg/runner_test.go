package pkg

import (
	"github.com/elek/flekszible/api/processor"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateTransformation(t *testing.T) {

	proc, err := parseTransformation("image:image=asd")
	assert.Nil(t, err)

	image := proc[0].(*processor.Image)

	assert.Equal(t, "asd", image.Image)
}

func TestCreateTransformationImageWithColon(t *testing.T) {

	proc, err := parseTransformation("image:image=localhost:5000/image/name:tag")
	assert.Nil(t, err)

	image := proc[0].(*processor.Image)

	assert.Equal(t, "localhost:5000/image/name:tag", image.Image)
}

func TestCreateTransformationMultiple(t *testing.T) {

	proc, err := parseTransformation("image:image=localhost:5000/image/name:tag;mount:path=/tmp,nhostPath=/tmp/test")
	assert.Nil(t, err)

	image := proc[0].(*processor.Image)

	assert.Equal(t, "localhost:5000/image/name:tag", image.Image)

	mount := proc[1].(*processor.Mount)

	assert.Equal(t, "/tmp", mount.Path)
	assert.Equal(t, "/tmp/test", mount.HostPath)
}
