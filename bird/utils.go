package bird

import (
	"regexp"
)

/*
 Template Replace:
 See https://golang.org/pkg/regexp/#Regexp.Expand
 for a template reference.
*/
func TemplateExpand(expr, template, input string) string {
	reg := regexp.MustCompile(expr)
	match := reg.FindStringSubmatchIndex(input)

	dst := []byte{}
	res := reg.ExpandString(dst, template, input, match)

	return string(res)
}
