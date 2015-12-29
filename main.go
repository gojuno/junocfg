package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"

	"github.com/juno-lab/argparse"
)

func loadFile(filename string) ([]byte, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte(""), err
	}
	return content, nil
}

func getTemplate(filename string) (*template.Template, error) {
	content, err := loadFile(filename)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Template file load error: [%v]\n", err))
	}
	tmpl, err := template.New("template").Parse(string(content))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable parse data (%v) %q as a template: [%v]\n", filename, string(content), err))
	}
	return tmpl, nil
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

func getConfig(filename string) (map[string]interface{}, error) {
	buffer := bytes.NewBuffer([]byte{})

	if filename == "<STDIN>" {
		info, _ := os.Stdin.Stat()
		// outMode(info.Mode())
		if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
			return nil, errors.New(fmt.Sprintf("The command is intended to work with pipes\n"))
		} else {
			reader := bufio.NewReader(os.Stdin)
			for {
				input, err := reader.ReadString('\n')
				if err != nil && err == io.EOF {
					break
				}
				buffer.WriteString(input)
			}
		}
	} else {
		content, err := loadFile(filename)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Config file load error: [%v]\n", err))
		}
		buffer.Write(content)
	}

	cfg := map[string]interface{}{}
	if err := yaml.Unmarshal(buffer.Bytes(), &cfg); err != nil {
		return nil, errors.New(fmt.Sprintf("Could not parse YAML file: %s\n", err))
	}
	return cfg, nil
}

func outResult(filename string, buffer *bytes.Buffer) {
	outputBuffer := bufio.NewWriter(os.Stdout)
	if filename != "<STDOUT>" {
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

func main() {
	parser, _ := argparse.ArgumentParser()
	parser.AddStringOption("input", "i", "input").Default("<STDIN>")
	parser.AddStringOption("output", "o", "output").Default("<STDOUT>")
	parser.AddStringOption("template", "t", "tmpl")
	parser.AddFlagOption("check", "", "check").Default("false").Action(argparse.SET_TRUE)

	args := parser.ParseArgs()

	// if *tmplFile == "" {
	// 	fmt.Fprintf(os.Stderr, "template (-t) file required\n")
	// 	os.Exit(1)
	// }

	var tmpl *template.Template
	var cfg map[string]interface{}
	var err error

	tmpl, err = getTemplate(args.AsString("template"))
	if err != nil {
		log.Fatal(err.Error())
	}

	cfg, err = getConfig(args.AsString("input"))
	if err != nil {
		log.Fatal(err.Error())
	}

	buffer := bytes.NewBuffer([]byte{})
	if err := tmpl.Execute(buffer, cfg); err != nil {
		log.Fatalf("failed to render template [%s]\n[%s]\n", err, cfg)
	}

	if args.AsFlag("check") {
		strOut := strings.Split(buffer.String(), "\n")

		for posInFile, str := range strOut {
			if i := strings.Index(str, "<no value>"); i != -1 {
				fmt.Fprintf(os.Stderr, "<no value> at %s#%d:%s\n", args.AsString("output"), posInFile, str)
			}
		}

		outYaml := map[string]interface{}{}
		err = yaml.Unmarshal(buffer.Bytes(), &outYaml)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Not valid output yaml: %s", err.Error())
		}
		// TODO! check output
		// find <no value> substring
	} else {
		outResult(args.AsString("output"), buffer)
	}
}
