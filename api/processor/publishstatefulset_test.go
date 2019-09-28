package processor

import (
	"testing"
)

func TestPublishStatefulset(t *testing.T) {
	TestFromDir(t, "publishstatefulset")
}

func TestPublishStatefulsetWithPort(t *testing.T) {
	TestFromDir(t, "publishstatefulset-withport")
}
