package schema

import (
	"strings"
	"testing"

	"github.com/sdinsure/agent/pkg/logger"

	"github.com/stretchr/testify/assert"
)

func TestTraverseMap(t *testing.T) {

	traversedMap := make(map[string]interface{})
	TraverseMap(logger.NewLogger(), m1, []string{}, func(prefixes []string, current string, val any) error {
		path := strings.Join(append(prefixes, current), ".")
		traversedMap[path] = val
		return nil
	})

	expectedMap := map[string]interface{}{
		"foo":                           "bar",
		"fooInt":                        64,
		"foostruct.foostructslice":      []string{"val1", "val2", "val3"},
		"foostruct.foostructslicefloat": []float64{6.14, 12.245},
		"foostruct.foostructstring":     "bar",
		"foostruct.foostructstruct.foostructstructstring": "bar",
	}
	assert.EqualValues(t, expectedMap, traversedMap)

}

var (
	m1 = map[string]interface{}{
		"foo":    "bar",
		"fooInt": 64,
		"foostruct": map[string]interface{}{
			"foostructstring": "bar",
			"foostructslice": []string{
				"val1",
				"val2",
				"val3",
			},
			"foostructslicefloat": []float64{
				6.14,
				12.245,
			},
			"foostructstruct": map[string]interface{}{
				"foostructstructstring": "bar",
			},
		},
	}
)
