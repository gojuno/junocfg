package junocfg

// https://golang.org/pkg/testing/

import (
	"bytes"
	"testing"
)

var jsonTests = []struct {
	in  [][]byte
	out []byte
}{
	{
		[][]byte{
			[]byte(`{"a": "aaa"}`),
			[]byte(`{"b": "bbb"}`),
		},
		[]byte(`{"a":"aaa","b":"bbb"}`),
	},
	{
		[][]byte{
			[]byte(`{"a":"aaa","b":"aaa"}`),
			[]byte(`{"b":"bbb","c": "ccc"}`),
		},
		[]byte(`{"a":"aaa","b":"bbb","c":"ccc"}`),
	},
}

func TestMergeJsons(t *testing.T) {
	for i, tst := range jsonTests {
		if out, err := mergeJsons(tst.in); err != nil {
			t.Error(
				"For", i,
				"got unexpected error", err,
			)
		} else if bytes.Compare(out, tst.out) != 0 {
			t.Error(
				"For", i,
				"expect", string(tst.out),
				"got", string(out),
			)
		}
	}
}
