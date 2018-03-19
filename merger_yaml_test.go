package junocfg

// https://golang.org/pkg/testing/

import (
	"bytes"
	"testing"
)

var yamlTests = []struct {
	in  [][]byte
	out []byte
}{
	{
		[][]byte{
			[]byte(`{"a": "aaa"}`),
			[]byte(`{"b": "bbb"}`),
		},
		[]byte("a: aaa\nb: bbb\n"),
	},
	{
		[][]byte{
			[]byte(`{"a":"aaa","b":"aaa"}`),
			[]byte(`{"b":"bbb","c": "ccc"}`),
		},
		[]byte("a: aaa\nb: bbb\nc: ccc\n"),
	},
	{
		[][]byte{
			[]byte(`{"a":"aaa","b":{"bb":"aaa"}}`),
			[]byte(`{"b":"bbb","c": "ccc"}`),
		},
		[]byte("a: aaa\nb: bbb\nc: ccc\n"),
	},
	{
		[][]byte{
			[]byte(`{"a":"aaa","b":{"b1":"1111","b2":"1111"}}`),
			[]byte(`{"b":{"b2":"2222","b3":"3333"},"c": "ccc"}`),
		},
		[]byte("a: aaa\nb:\n  b1: \"1111\"\n  b2: \"2222\"\n  b3: \"3333\"\nc: ccc\n"),
	},
}

func TestMergeYamls(t *testing.T) {
	for i, tst := range yamlTests {
		if out, err := MergeYamls(tst.in); err != nil {
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
