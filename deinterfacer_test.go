package junocfg

// https://golang.org/pkg/testing/

import (
	"encoding/json"
	"reflect"
	//"bytes"

	"testing"
)

var deinterfacerTests = []struct {
	in  map[string]interface{}
	out map[string]interface{}
}{
	{
		in: map[string]interface{}{
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
			"f": []int{1, 2, 3, 4},
		},
		out: map[string]interface{}{
			"a": "a",
			"b": []interface{}{"b1", "b2", "b3", "b4"},
			"c": map[string]interface{}{
				"ca": "ca",
				"cb": []interface{}{
					"cb1",
					"cb2",
					map[string]interface{}{
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
			"f": []interface{}{1, 2, 3, 4},
		},
	},
	{in: nil, out: map[string]interface{}{}},
	{in: map[string]interface{}{}, out: map[string]interface{}{}},
	{
		in: map[string]interface{}{
			"key_0": map[int]int{
				1:  1,
				2:  3,
				5:  8,
				13: 21,
			},
		},
		out: map[string]interface{}{
			"key_0": map[string]interface{}{
				"1":  1,
				"2":  3,
				"5":  8,
				"13": 21,
			},
		},
	},
}

func Test_deinterfacer(t *testing.T) {
	for i, tst := range deinterfacerTests {
		out := deinterface(tst.in)
		if !reflect.DeepEqual(tst.out, out) {
			t.Errorf("For %d \nexpected %#v \ngot %#v", i, tst.out, out)
		}
	}
}

func Test_deinterfacereMarshallable(t *testing.T) {
	for i, tst := range deinterfacerTests {
		out := deinterface(tst.in)
		_, err := json.Marshal(out)
		if err != nil {
			t.Errorf("For %d got error while marshalling %s", i, err)
		}
	}
}
