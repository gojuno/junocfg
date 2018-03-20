package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mguzelevich/go.log"
)

type inputData struct {
	input [][]byte
	err   error
}

func (i *inputData) dump() string {
	buffer := bytes.NewBuffer([]byte{})
	for idx, d := range i.input {
		buffer.WriteString(fmt.Sprintf("=== %d ===\n%v\n", idx, string(d)))
	}
	buffer.WriteString(fmt.Sprintf("err: %v", i.err))
	return buffer.String()
}

func (i *inputData) readFiles(filenames string) {
	for _, filename := range strings.Split(filenames, ",") {
		i.readFile(filename)
	}
}

func (i *inputData) readFile(filename string) {
	if i.err != nil {
		return
	}
	if data, err := ioutil.ReadFile(filename); err != nil {
		i.err = fmt.Errorf("Config file load error: [%v]\n", err)
	} else {
		i.input = append(i.input, data)
	}
}

func (i *inputData) readStdin() {
	info, err := os.Stdin.Stat()
	if err != nil || (info.Mode()&os.ModeCharDevice) == os.ModeCharDevice {
		outMode(info.Mode())
		panic(fmt.Errorf("cannt check stdin stat [%v] or incorrect mode: [%v]\n", err, info.Mode()))
	} else if info.Size() > 0 {
		scanner := bufio.NewScanner(os.Stdin)
		buffer := bytes.NewBuffer([]byte{})

		for scanner.Scan() { // internally, it advances token based on sperator
			line := scanner.Text()
			buffer.Write([]byte(fmt.Sprintf("%s\n", line)))
		}
		i.input = append(i.input, buffer.Bytes())
	} else {
		panic(fmt.Errorf("empty STDIN"))
	}
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

	log.Stderr.Printf("info: %032b", mode)
	for flag, name := range flags {
		if (mode & flag) == flag {
			log.Stderr.Printf("%s\n", name)
		}
	}
}

func initInput() *inputData {
	return &inputData{
		input: [][]byte{},
	}
}

func getInput(inputString string) (*inputData, error) {
	in := initInput()
	if inputString != "" {
		in.readFiles(inputString)
	} else {
		in.readStdin()
	}
	return in, in.err
}
