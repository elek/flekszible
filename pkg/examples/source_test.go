package examples

import (
	"github.com/elek/flekszible/api/processor"
	"testing"
)

func TestSource(t *testing.T) {
	processor.TestFromDir(t, "source")
}
