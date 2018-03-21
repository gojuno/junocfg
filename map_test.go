package junocfg

// https://golang.org/pkg/testing/

import (
	//"bytes"
	"testing"
)

var setTests = []struct {
	src   map[string]interface{}
	key   []string
	value interface{}
	out   string
}{
	{
		map[string]interface{}{},
		[]string{"a", "b", "c"},
		"!!!!",
		"",
	},
	{
		map[string]interface{}{
			"a": "a",
			"b": []string{"aa", "aa", "aa", "aa"},
			"c": map[string]interface{}{
				"a": "a",
				"b": []string{"aa", "aa", "aa", "aa"},
			},
			"d": "a",
		},
		[]string{"a", "c"}, "!!!!",
		"",
	},
	{
		map[string]interface{}{
			"a": "a",
			"b": []string{"aa", "aa", "aa", "aa"},
			"c": map[string]interface{}{
				"a": "a",
				"b": []string{"aa", "aa", "aa", "aa"},
			},
			"d": "a",
		},
		[]string{"c", "a"}, "!!!!",
		"",
	},
}

func Test_setValue(t *testing.T) {
	for i, d := range setTests {
		err := setValue(d.src, d.key, d.value)
		d, err := walk(d.src)
		t.Logf("walker: %v\n%v\n", d, err)
		if err != nil {
			t.Errorf("for %d error detected", i)
		}
	}
}

func Test_Items2Map(t *testing.T) {
	// src := map[string]interface{}{
	// 	"a": "a",
	// 	"b": []string{"aa", "aa", "aa", "aa"},
	// 	"c": map[string]interface{}{
	// 		"a": "a",
	// 		"b": []string{"aa", "aa", "aa", "aa"},
	// 	},
	// 	"d": "a",
	// }
	// t.Logf("src %v", src)
	// err := setValue(src, []string{"a", "c"}, "!!!!")
	// t.Logf("dst %v", src)
	// t.Error("...")
	// if err != nil {
	// 	t.Error(err)
	// }
}
