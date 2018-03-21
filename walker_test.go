package junocfg

// https://golang.org/pkg/testing/

import (
	"encoding/json"
	"fmt"
	"testing"
)

var walkTests = []struct {
	in  []byte
	out string
}{
	{
		[]byte(`{"a": "aaa"}`),
		"[{[a] aaa}]",
	},
	{
		[]byte(`{"a":"a","b": "b","c": "c"}`),
		"[{[a] a} {[b] b} {[c] c}]",
	},
	{
		[]byte(`{"a":{"b":{"c":"d"}}}`),
		"[{[a b c] d}]",
	},
}

func Test_walk(t *testing.T) {
	for i, td := range walkTests {
		data := map[string]interface{}{}
		err := json.Unmarshal(td.in, &data)
		if err != nil {
			t.Errorf("For %d input data error %v", i, td.in)
		}
		d, err := walk(data)
		if err != nil {
			t.Errorf("For %d walk error %v", i, err)
		}
		out := fmt.Sprintf("%v", d)
		if out != td.out {
			t.Errorf("For %d error", i)
			t.Logf("\tin: %v", data)
			t.Logf("\texpected: %v", td.out)
			t.Logf("\tgot: %v", out)
		}
	}

}
