package processor

import (
	"testing"
)

func TestInitializer(t *testing.T) {
	ExecuteProcessorAndCompare(t, "initializer", "ss")
}
