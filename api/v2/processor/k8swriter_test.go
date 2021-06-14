package processor

import (
	"testing"

	"github.com/elek/flekszible/api/v2/data"
	"github.com/stretchr/testify/assert"
)

func TestK8sWriter(t *testing.T) {
	node, err := data.ReadManifestFile("../../testdata/writer/ss.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, node)

	resource := data.Resource{
		Content: node,
	}

	writer := CreateStdK8sWriter()

	err = writer.BeforeResource(&resource)
	assert.Nil(t, err)
	node.Accept(writer)

}
