package processor

import (
	"testing"
)

func TestComposit(t *testing.T) {
	TestFromDir(t, "composite")
}

func TestCompositInherited(t *testing.T) {
	TestFromDir(t, "composite-inherited")
}

func TestCompositWithParam(t *testing.T) {
	TestFromDir(t, "composite-param")
}

func TestCompositWithAdditionalResources(t *testing.T) {
	TestFromDir(t, "composite-resources")
}
