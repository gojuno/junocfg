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

/*
func Test() {
	data := []struct {
	}{
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
	for i, d := range setTests {
		err := setValue(d.src, d.key, d.value)
		d, err := walk(d.src)
		t.Logf("walker: %v\n%v\n", d, err)
		if err != nil {
			t.Errorf("for %d error detected", i)
		}
	}
}

func TestMergeMaps(t *testing.T) {
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
