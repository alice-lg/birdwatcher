package bird

import (
	"testing"
)

func TestTemplateExpand(t *testing.T) {
	// Expect input to be something like
	//   TEST_FOO_42_src
	//
	// Template:
	//   RESULT_$1
	//
	// Result:
	//    RESULT_FOO_42
	//
	src := "TEST_FOO_42_src"
	expr := `TEST_(.*)_src`
	tmpl := "RESULT_${1}_${1}"

	res := TemplateExpand(expr, tmpl, src)

	if res != "RESULT_FOO_42_FOO_42" {
		t.Error("Unexpacted expansion:", res)
	}
}
