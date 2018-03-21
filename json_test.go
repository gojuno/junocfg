package junocfg

// https://golang.org/pkg/testing/

import (
	"fmt"
	"testing"
)

var json2itemsTests = []struct {
	in  []byte
	out string
}{
	{
		[]byte(`{"a": "aaa"}`),
		"[{[a] aaa}]",
	},
	{
		[]byte(`{"a":"aaa","b":"aaa"}`),
		"[{[a] aaa} {[b] aaa}]",
	},
	{
		[]byte(`{"a":"aaa","b":{"bb":"aaa"}}`),
		"[{[a] aaa} {[b bb] aaa}]",
	},
	{
		[]byte(`{"a":"aaa","b":{"b1":"1111","b2":"1111"}}`),
		"[{[a] aaa} {[b b1] 1111} {[b b2] 1111}]",
	},
}

func Test_json2Items(t *testing.T) {
	for i, tst := range json2itemsTests {
		ym, err := json2Items(tst.in)
		if err != nil {
			t.Errorf("For %d got unexpected error %v", i, err)
		}
		out := fmt.Sprintf("%s", ym)
		if out != tst.out {
			t.Errorf("For %d expected %v got %v", i, tst.out, out)
		}
	}

}
