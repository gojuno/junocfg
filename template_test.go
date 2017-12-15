package junocfg

// https://golang.org/pkg/testing/

import (
	"testing"
)

var tmplTests = []struct {
	name     string
	tmpl     []byte
	positive bool
}{
	{
		"incorrect braces",
		[]byte(`
a: {{.aa}
b: bbb
`),
		false,
	},
	{
		"incorrect braces",
		[]byte(`{"a": "aaa"}`),
		true,
	},
	{
		"incorrect braces",
		[]byte(`{"a": "aaa"}`),
		true,
	},
}

func TestCheckTmpl(t *testing.T) {
	for i, tst := range tmplTests {
		_, err := CheckTemplate(tst.tmpl)
		if tst.positive && err != nil {
			t.Errorf(
				"For %d '%s' got unexpected error %v in positive case", i, tst.name, err,
			)
		}
		if !tst.positive && err == nil {
			t.Errorf(
				"For %d '%s' lost expected error for negative case", i, tst.name,
			)
		}

	}

}
