package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"regexp"
)

type Change struct {
	DefaultProcessor
	Trigger     Trigger
	Path        data.Path
	Pattern     string
	Replacement string
}

func (processor *Change) BeforeResource(resource *data.Resource) {
	var re = regexp.MustCompile(processor.Pattern)
	if !processor.Trigger.active(resource) {
		return
	}
	getter := data.SmartGetAll{Path: processor.Path}
	resource.Content.Accept(&getter)
	for _, result := range getter.Result {
		key := result.Value.(*data.KeyNode)
		key.Value = re.ReplaceAllString(key.Value.(string), processor.Replacement)
	}
}

func init() {
	prototype := Change{}
	ProcessorTypeRegistry.Add(&prototype)
}
