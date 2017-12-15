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

func outMode(mode os.FileMode) {
	flags := map[os.FileMode]string{
		os.ModeDir:        "os.ModeDir",
		os.ModeAppend:     "os.ModeAppend",
		os.ModeExclusive:  "os.ModeExclusive",
		os.ModeTemporary:  "os.ModeTemporary",
		os.ModeSymlink:    "os.ModeSymlink",
		os.ModeDevice:     "os.ModeDevice",
		os.ModeNamedPipe:  "os.ModeNamedPipe",
		os.ModeSocket:     "os.ModeSocket",
		os.ModeSetuid:     "os.ModeSetuid",
		os.ModeSetgid:     "os.ModeSetgid",
		os.ModeCharDevice: "os.ModeCharDevice",
		os.ModeSticky:     "os.ModeSticky",
	}

	log.Printf("info: %032b", mode)
	for flag, name := range flags {
		if (mode & flag) == flag {
			log.Printf("%s\n", name)
		}
	}
}
