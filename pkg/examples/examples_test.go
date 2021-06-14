package examples

import (
	"testing"

	"github.com/elek/flekszible/api/v2/processor"
)

func TestGettingStarted(t *testing.T) {
	processor.TestExample(t, "gettingstarted")
}

func TestGettingEnvs(t *testing.T) {
	processor.TestExample(t, "envs/dev")
}

func TestInstantiate(t *testing.T) {
	processor.TestExample(t, "instantiate")
}
