package processor

import (
	"github.com/elek/flekszible/api/data"
	"github.com/stretchr/testify/assert"
	"testing"
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
