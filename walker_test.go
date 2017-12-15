package junocfg

// https://golang.org/pkg/testing/

import (
	"testing"
)

/*
func Test() {
	f := func(path []string, value interface{}) error {
		fmt.Printf("%v -> %v\n", path, value)
		return nil
	}

	fmt.Printf("\n=== TEST 1\n")
	src := map[string]interface{}{}
	fmt.Printf("src %v\n", src)
	dst, err := setValue(src, []string{"a", "b", "c"}, "!!!!")
	fmt.Printf("dst %v\n", dst)
	if err != nil {
		fmt.Printf("ERR %v\n", err)
	}
	walker(dst, []string{}, f)

	fmt.Printf("\n=== TEST 2\n")
	src = map[string]interface{}{
		"a": "a",
		"b": []string{"aa", "aa", "aa", "aa"},
		"c": map[string]interface{}{
			"a": "a",
			"b": []string{"aa", "aa", "aa", "aa"},
		},
		"d": "a",
	}
	fmt.Printf("src %v\n", src)
	dst, err = setValue(src, []string{"a", "c"}, "!!!!")
	fmt.Printf("dst %v\n", dst)
	if err != nil {
		fmt.Printf("ERR %v\n", err)
	}
	walker(dst, []string{}, f)

	fmt.Printf("\n=== TEST 3\n")
	src = map[string]interface{}{
		"a": "a",
		"b": []string{"aa", "aa", "aa", "aa"},
		"c": map[string]interface{}{
			"a": "a",
			"b": []string{"aa", "aa", "aa", "aa"},
		},
		"d": "a",
	}
	fmt.Printf("src %v\n", src)
	dst, err = setValue(src, []string{"c", "a"}, "!!!!")
	fmt.Printf("dst %v\n", dst)
	if err != nil {
		fmt.Printf("ERR %v\n", err)
	}
	walker(dst, []string{}, f)
}
*/
func TestWalker(t *testing.T) {

}
