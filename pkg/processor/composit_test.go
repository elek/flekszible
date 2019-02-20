package processor

import (
	"testing"
)

func TestComposit(t *testing.T) {
	TestFromDir(t, "composit")
}

func TestCompositInherited(t *testing.T) {
	TestFromDir(t, "composit-inherited")
}
