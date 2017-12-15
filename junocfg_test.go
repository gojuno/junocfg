package junocfg

// https://golang.org/pkg/testing/

import (
	// "bytes"
	// "os"
	"fmt"
	"testing"
)

func TestRepot(t *testing.T) {
	result := true
	err := &result
	err = nil
	if err != nil {
		test := "Test"
		expected := 1
		out := "Err"
		t.Error(
			"For", test,
			"expected", expected,
			"got", err,
			"output\n", out,
		)
	}
}

func BenchmarkRepot(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmt.Sprintf("hello")
	}
}

func Example_suffix() {

}

func ExampleF_suffix() {

}

func ExampleT_suffix() {

}

func ExampleT_M_suffix() {

}
