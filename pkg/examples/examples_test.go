package examples

import (
	"github.com/elek/flekszible/public/processor"
	"testing"
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
