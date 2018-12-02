package processor

import "testing"

func TestImageSet(t *testing.T) {
	ExecuteProcessorAndCompare(t, "imageset", "ss")
}
