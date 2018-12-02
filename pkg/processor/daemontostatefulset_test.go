package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDaemonToStatefulSet(t *testing.T) {
	ctx := ExecuteProcessorAndCompare(t, "daemontostateful", "datanode")
	assert.Equal(t, 2, len(ctx.Resources))

	expected, err := data.LoadFrom("../../testdata/daemontostateful", "service_expected.yaml")
	assert.Nil(t, err)

	assert.EqualValues(t, ToSimpleYaml(&expected[0]), ToSimpleYaml(&ctx.Resources[1]))


}
