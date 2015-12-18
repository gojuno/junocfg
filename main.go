package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
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
		fmt.Fprintf(os.Stderr, "Template file load error: [%v]\n", err)
		return nil, errors.New(fmt.Sprintf("Template file load error: [%v]\n", err))
	}
	tmpl, err := template.New("template").Parse(string(content))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable parse data (%v) %q as a template: [%v]\n", filename, string(content), err))
	}
	return tmpl, nil
}

func getConfig(filename string) (map[string]interface{}, error) {
	content, err := loadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config file load error: [%v]\n", err)
		return nil, errors.New(fmt.Sprintf("Config file load error: [%v]\n", err))
	}

	cfg := map[string]interface{}{}

	if err := yaml.Unmarshal(content, &cfg); err != nil {
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
			fmt.Fprintf(os.Stderr, "Output file create error: [%v]\n", err)
			os.Exit(1)
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
	check := flag.Bool("check", false, "check mode")
	tmplFile := flag.String("t", "", "config template")
	cfgFile := flag.String("c", "settings.dev.yaml", "config template")
	output := flag.String("o", "<STDOUT>", "output")
	flag.Parse()

	if *check {
		fmt.Fprintf(os.Stderr, "Data check...\n")
	}

	if *tmplFile == "" {
		fmt.Fprintf(os.Stderr, "template (-t) file required\n")
		os.Exit(1)
	}

	var tmpl *template.Template
	var cfg map[string]interface{}
	var err error

	tmpl, err = getTemplate(*tmplFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	cfg, err = getConfig(*cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	buffer := bytes.NewBuffer([]byte{})
	if err := tmpl.Execute(buffer, cfg); err != nil {
		panic(fmt.Sprintf("failed to render template [%s]\n[%s]\n", err, cfg))
	}

	if *check {
		// TODO! check output
		// find <no value> substring
	} else {
		outResult(*output, buffer)
	}
}
