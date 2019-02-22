package processor

import (
	"github.com/elek/flekszible/pkg/data"
	"github.com/sirupsen/logrus"
	"regexp"
	"strconv"
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
		switch key.Value.(type) {
		case int:
			intval, err := strconv.Atoi(processor.Replacement)
			key.Value = intval
			if err != nil {
				logrus.Error("Invalid replacement value: number can be replaced only with number and not: " + processor.Replacement)
			}
		case string:
			key.Value = re.ReplaceAllString(key.Value.(string), processor.Replacement)

		default:
			logrus.Error("Invalid replacement only string or int can be replaced: " + result.Path.ToString())
		}

	}
}

func init() {
	prototype := Change{}
	ProcessorTypeRegistry.Add(&prototype)
}
