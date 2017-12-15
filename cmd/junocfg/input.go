package main

import (
	"bufio"
	"bytes"
	// "errors"
	"fmt"
	//"io"
	"io/ioutil"
	//"log"
	"os"
	"strings"
	//	"github.com/gojuno/junocfg"
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
	if fi, err := os.Stdin.Stat(); err != nil || fi.Mode()&os.ModeNamedPipe == 0 {
		panic(fmt.Errorf("cannt check stdin stat [%v] or incorrect mode: [%v]\n", err, fi.Mode()))
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		buffer := bytes.NewBuffer([]byte{})

		for scanner.Scan() { // internally, it advances token based on sperator
			line := scanner.Text()
			buffer.Write([]byte(fmt.Sprintf("%s\n", line)))
		}
		i.input = append(i.input, buffer.Bytes())
	}
}

func initInput() *inputData {
	return &inputData{
		input: [][]byte{},
	}
}

func getInput(inputString string) (*inputData, error) {
	in := initInput()
	if input != "" {
		in.readFiles(inputString)
	} else {
		in.readStdin()
	}
	return in, in.err
}
