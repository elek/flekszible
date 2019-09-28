package processor

import (
	"github.com/elek/flekszible/public/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestK8sWriter(t *testing.T) {
	node, err := data.ReadManifestFile("../../testdata/writer/ss.yaml")
	assert.Nil(t, err)
	assert.NotNil(t, node)

	resource := data.Resource{}

	writer := CreateStdK8sWriter()

	writer.BeforeResource(&resource)
	node.Accept(writer)

}
