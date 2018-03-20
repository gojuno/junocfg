package junocfg

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

type item struct {
	path  []string
	value interface{}
}

func (i *item) pathString() string {
	return strings.Join(i.path, " / ")
}

type itemArray []item

func (a itemArray) Len() int { return len(a) }
func (a itemArray) Less(i, j int) bool {
	return strings.Compare(a[i].pathString(), a[j].pathString()) < 1
}
func (a itemArray) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func writeStr(buffer *bytes.Buffer, str string) {
	if strings.TrimSpace(str) != "" {
		buffer.WriteString(str)
		buffer.WriteString("\n")
	}
}

func OutResult(filename string, buffer *bytes.Buffer) {
	outputBuffer := bufio.NewWriter(os.Stdout)
	if filename != "" {
		f, err := os.Create(filename)
		defer f.Close()
		if err != nil {
			log.Fatalf("Output file create error: [%v]\n", err)
		}
		outputBuffer = bufio.NewWriter(f)
	}

	if _, err := outputBuffer.WriteString(buffer.String()); err != nil {
		fmt.Fprintf(os.Stderr, "Output write error: [%v]\n", err)
		os.Exit(1)
	}
	outputBuffer.Flush()
}
