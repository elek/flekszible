package processor

import "testing"

func TestNameFilter(t *testing.T) {
	TestFromDir(t, "namefilter-include")
}
