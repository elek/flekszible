package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPublishStatefulset(t *testing.T) {
	ctx := ExecuteProcessorAndCompare(t, "publishstatefulset", "datanode")
	assert.Equal(t, 2, len(ctx.Resources))

	println("-----")
	expected, err := data.LoadFrom("../../testdata/publishstatefulset", "datanode_generated.yaml")
	assert.Nil(t, err)

	assert.EqualValues(t, ToSimpleYaml(&expected[0]), ToSimpleYaml(&ctx.Resources[1]))

}
