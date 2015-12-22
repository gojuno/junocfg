package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
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
	buffer := bytes.NewBuffer([]byte{})

	if filename == "<STDIN>" {
		info, _ := os.Stdin.Stat()
		if (info.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
			return nil, errors.New(fmt.Sprintf("The command is intended to work with pipes\n"))
		} else if info.Size() > 0 {
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
			fmt.Fprintf(os.Stderr, "Config file load error: [%v]\n", err)
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
	cfgFile := flag.String("c", "<STDIN>", "config template")
	output := flag.String("o", "<STDOUT>", "output")
	flag.Parse()

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
		strOut := strings.Split(buffer.String(), "\n")

		for posInFile, str := range strOut {
			if i := strings.Index(str, "<no value>"); i != -1 {
				fmt.Fprintf(os.Stderr, "<no value> at %s#%d:%s\n", *output, posInFile, str)
			}
		}

		outYaml := map[string]interface{}{}
		err = yaml.Unmarshal(buffer.Bytes(), &outYaml)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Not valid output yaml: %s", err.Error())
		}
	}

	outResult(*output, buffer)
}
