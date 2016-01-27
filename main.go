package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"

	"github.com/juno-lab/argparse"
)

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

func getConfig(filenames string) (map[string]interface{}, error) {
	buffers := []*bytes.Buffer{}

	if filenames == "<STDIN>" {
		info, _ := os.Stdin.Stat()
		// outMode(info.Mode())
		if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
			return nil, errors.New(fmt.Sprintf("The command is intended to work with pipes\n"))
		} else {
			buffer := bytes.NewBuffer([]byte{})
			reader := bufio.NewReader(os.Stdin)
			for {
				input, err := reader.ReadString('\n')
				if err != nil && err == io.EOF {
					break
				}
				buffer.WriteString(input)
			}
			buffers = append(buffers, buffer)
		}
	} else {
		for _, filename := range strings.Split(filenames, ",") {
			buffer := bytes.NewBuffer([]byte{})
			content, err := loadFile(filename)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Config file load error: [%v]\n", err))
			}
			buffer.Write(content)
			buffers = append(buffers, buffer)
		}
	}

	config := map[string]interface{}{}
	for _, buffer := range buffers {
		cfg := map[string]interface{}{}
		if err := yaml.Unmarshal(buffer.Bytes(), &cfg); err != nil {
			return nil, errors.New(fmt.Sprintf("Could not parse YAML file: %s\n", err))
		}
		config = mergeMaps(config, cfg)
	}

	return config, nil
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
	parser.AddStringOption("template", "t", "tmpl").Default("--")
	parser.AddFlagOption("check", "", "check").Default("false").Action(argparse.SET_TRUE)
	parser.AddFlagOption("merge", "", "merge").Default("false").Action(argparse.SET_TRUE)

	args := parser.ParseArgs()

	// if *tmplFile == "" {
	// 	fmt.Fprintf(os.Stderr, "template (-t) file required\n")
	// 	os.Exit(1)
	// }

	var cfg map[string]interface{}
	var err error

	cfg, err = getConfig(args.AsString("input"))
	if err != nil {
		log.Fatal(err.Error())
	}

	correct := true
	buffer := bytes.NewBuffer([]byte{})

	if args.AsFlag("merge") {
		d, err := yaml.Marshal(cfg)
		if err != nil {
			log.Fatal(fmt.Sprintf("Could not create YAML file: %s\n", err))
		}
		buffer.Write(d)
	} else {
		var tmpl *template.Template

		if args.AsString("template") == "--" {
			log.Fatal("ParseArgs: Field [template(-t|--tmpl)] required")
		}
		tmpl, err = getTemplate(args.AsString("template"))
		if err != nil {
			log.Fatal(err.Error())
		}

		if err := tmpl.Execute(buffer, cfg); err != nil {
			log.Fatalf("failed to render template [%s]\n[%s]\n", err, cfg)
		}

		if args.AsFlag("check") {
			strOut := strings.Split(buffer.String(), "\n")

			for posInFile, str := range strOut {
				if i := strings.Index(str, "<no value>"); i != -1 {
					fmt.Fprintf(os.Stderr, "<no value> at %s#%d:%s\n", args.AsString("output"), posInFile, str)
					correct = correct && false
				}
			}

			outYaml := map[string]interface{}{}
			err = yaml.Unmarshal(buffer.Bytes(), &outYaml)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Not valid output yaml: %s", err.Error())
				correct = correct && false
			}
			fmt.Fprintf(os.Stderr, "\n")
		}
	}

	outResult(args.AsString("output"), buffer)

	if !correct {
		os.Exit(1)
	}
}
