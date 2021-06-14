package processor

import (
	"testing"
)

func TestKustomize(t *testing.T) {
	TestFromDir(t, "kustomize")
}
