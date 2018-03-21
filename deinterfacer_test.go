package junocfg

// https://golang.org/pkg/testing/

import (
	//"bytes"
	"encoding/json"
	"testing"
)

var deinterfacerTests = []struct {
	in  map[string]interface{}
	out string
}{
	{
		map[string]interface{}{
			"a": "a",
			"b": []string{"b1", "b2", "b3", "b4"},
			"c": map[interface{}]interface{}{
				"ca": "ca",
				"cb": []interface{}{
					"cb1",
					"cb2",
					map[interface{}]interface{}{
						"cb3a": "cb3aa",
					},
					"cb4",
					map[string]interface{}{
						"cb5a": "cb5aa",
					},
				},
			},
			"d": "a",
			"e": map[string]interface{}{
				"ea": "ea",
				"eb": "eb",
				"ec": "ec",
			},
		},
		`{"a":"a","b":["b1","b2","b3","b4"],"c":{"ca":"ca","cb":["cb1","cb2",{"cb3a":"cb3aa"},"cb4",{"cb5a":"cb5aa"}]},"d":"a","e":{"ea":"ea","eb":"eb","ec":"ec"}}`,
	},
}

func Test_deinterfacer(t *testing.T) {
	for i, tst := range deinterfacerTests {
		out, err := deinterface(tst.in)
		if err != nil {
			t.Errorf("For %d deinterface error %v", i, err)
		}
		outStr, err := json.Marshal(out)
		if string(outStr) != tst.out {
			t.Errorf("For %d expected %v got %v", i, tst.out, string(outStr))
		}
	}
}
