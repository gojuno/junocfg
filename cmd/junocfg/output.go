package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	//	"github.com/gojuno/junocfg"
)

func outResult(filename string, out string) {
	outputBuffer := bufio.NewWriter(os.Stdout)
	if filename != "" {
		f, err := os.Create(filename)
		defer f.Close()
		if err != nil {
			log.Fatalf("Output file create error: [%v]\n", err)
		}
		outputBuffer = bufio.NewWriter(f)
	}

	if _, err := outputBuffer.WriteString(out); err != nil {
		fmt.Fprintf(os.Stderr, "Output write error: [%v]\n", err)
		os.Exit(1)
	}
	outputBuffer.Flush()
}
