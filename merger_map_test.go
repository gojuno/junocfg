package junocfg

// https://golang.org/pkg/testing/

import (
	//"bytes"
	"testing"
)

var mapTests = []struct {
	in  []map[string]interface{}
	out map[string]interface{}
}{
	{
		[]map[string]interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		},
		map[string]interface{}{},
	},
}

/*
func Test() {
	data := []struct {
		src   map[string]interface{}
		key   []string
		value interface{}
	}{
		{
			map[string]interface{}{},
			[]string{"a", "b", "c"}, "!!!!",
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
		},
	}
	for i, d := range data {
		fmt.Printf("\n=== TEST %d\n", i)
		fmt.Printf("src %v\n", d.src)
		err := setValue(d.src, d.key, d.value)
		fmt.Printf("dst %v\n", d.src)
		if err != nil {
			fmt.Printf("ERR %v\n", err)
		}
		d, err := walk(d.src)
		fmt.Printf("walker: %v\n%v\n", d, err)
	}
}

*/
func TestSetValue(t *testing.T) {
}

func TestMergeMaps(t *testing.T) {
	src := map[string]interface{}{
		"a": "a",
		"b": []string{"aa", "aa", "aa", "aa"},
		"c": map[string]interface{}{
			"a": "a",
			"b": []string{"aa", "aa", "aa", "aa"},
		},
		"d": "a",
	}
	t.Logf("src %v", src)
	dst, err := setValue(src, []string{"a", "c"}, "!!!!")
	t.Logf("dst %v", dst)
	t.Error("...")
	if err != nil {
		t.Error(err)
	}
}
