package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"

	"github.com/gojuno/junocfg"
)

var (
	check     bool
	checkTmpl bool
	merge     bool

	input  string
	output string
	tmpl   string
)

func execute() {
	filenames := strings.Split(input, ",")

	correct := true

	switch {
	case checkTmpl:
		// fmt.Fprintf(os.Stderr, "mode: check-tmpl\n")
		if rawTmpl, err := junocfg.ReadData(tmpl); err != nil {
			log.Fatal(err.Error())
			correct = false
		} else {
			if _, err := template.New("template").Parse(rawTmpl); err != nil {
				fmt.Fprintf(os.Stderr, "template error: [%v]\n", err)
				correct = false
			}
		}
	case merge:
		// fmt.Fprintf(os.Stderr, "mode: merge\n")
		if cfg, err := junocfg.GetConfigs(filenames); err == nil {
			d, err := yaml.Marshal(cfg)
			if err != nil {
				log.Fatal(fmt.Sprintf("Could not create YAML file: %s\n", err))
			}
			buffer := bytes.NewBuffer([]byte{})
			buffer.Write(d)
			junocfg.OutResult(output, buffer)
		} else {
			log.Fatalf("%s\n", err)
		}
	default:
		// fmt.Fprintf(os.Stderr, "mode: default\n")
		// default: generate config file from input + template
		buffer := bytes.NewBuffer([]byte{})

		if tmpl == "" {
			log.Fatal("Field [template(-t|--tmpl)] required")
		}
		tmpl, err := junocfg.GetTemplate(tmpl)
		if err != nil {
			log.Fatal(err.Error())
		}

		cfg := make(map[string]interface{})
		if config, err := junocfg.GetConfigs(filenames); err != nil {
			log.Fatalf("%s\n", err)
		} else {
			cfg = config
		}

		if err := tmpl.Execute(buffer, cfg); err != nil {
			log.Fatalf("failed to render template [%s]\n[%s]\n", err, cfg)
		}

		buffer = junocfg.PreprocessYaml(buffer)

		// check yaml
		outYaml := map[string]interface{}{}
		if err = yaml.Unmarshal(buffer.Bytes(), &outYaml); err != nil {
			fmt.Fprintf(os.Stderr, "Not valid output yaml: %s\n", err.Error())
			correct = correct && false
		}

		// check variables
		if check {
			strOut := strings.Split(buffer.String(), "\n")

			for posInFile, str := range strOut {
				if i := strings.Index(str, "<no value>"); i != -1 {
					fmt.Fprintf(os.Stderr, "<no value> at %s#%d:%s\n", output, posInFile, str)
					correct = correct && false
				}
			}
		}
		junocfg.OutResult(output, buffer)
	}

	if !correct {
		os.Exit(1)
	}
}

func init() {
	flag.BoolVar(&check, "check", false, "check")
	flag.BoolVar(&checkTmpl, "check-tmpl", false, "check tmpl")
	flag.BoolVar(&merge, "merge", false, "merge")

	flag.StringVar(&input, "i", "", "input")
	flag.StringVar(&input, "input", "", "input")
	flag.StringVar(&output, "o", "", "output")
	flag.StringVar(&output, "output", "", "output")
	flag.StringVar(&tmpl, "t", "", "template")
	flag.StringVar(&tmpl, "template", "", "template")
}

func main() {
	flag.Parse()
	execute()
}
