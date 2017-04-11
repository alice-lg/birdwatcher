package bird

import (
	"reflect"
	"testing"
)

func TestParseBgpRoutes(t *testing.T) {

	inputs := []string{
		"1 imported, 0 filtered, 2 exported, 1 preferred",
		"0 imported, 2846 exported", // Bird 1.4.x
	}

	expected := []Parsed{
		Parsed{
			"imported":  int64(1),
			"filtered":  int64(0),
			"exported":  int64(2),
			"preferred": int64(1),
		},
		Parsed{
			"imported": int64(0),
			"exported": int64(2846),
		},
	}

	for i, in := range inputs {
		routes := parseBgpRoutes(in)
		if !reflect.DeepEqual(routes, expected[i]) {
			t.Error("Parse bgpRoutes:", routes, "expected:", expected[i])
		}
	}

}
