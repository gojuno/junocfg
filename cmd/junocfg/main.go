package main

import (
	"bytes"
	"flag"
	"fmt"
	stllog "log"
	"os"
	"strings"
	// "errors"
	// "io"

	"github.com/mguzelevich/go.log"

	"github.com/gojuno/junocfg"
)

var (
	checkTmpl bool
	merge     bool

	input  string
	output string
	tmpl   string
)

func init() {
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

	log.InitLoggers(&log.Logger{
		os.Stderr, // ioutil.Discard,
		os.Stderr, // ioutil.Discard,
		os.Stderr,
		os.Stderr,
		os.Stderr,
	})
	log.Stdout = stllog.New(os.Stdout, "", 0)
	log.Stderr = stllog.New(os.Stderr, "", 0)

	success := true
	switch {
	case checkTmpl:
		// fmt.Fprintf(os.Stderr, "mode: check-tmpl\n")
		in, err := getInput(tmpl)
		log.Debug.Printf(in.dump())
		if err != nil {
			fmt.Printf("Error %v", in.err)
			success = false
		} else if _, err := junocfg.CheckTemplate(in.input[0]); err != nil {
			fmt.Printf("check tmpl error %v", err)
			success = false
		}
	case merge:
		// fmt.Fprintf(os.Stderr, "mode: merge\n")
		in, err := getInput(input)
		if err != nil {
			fmt.Printf("Error %v", in.err)
			success = false
		}
		out, err := junocfg.MergeYamls(in.input)
		// in.dump()
		if err != nil {
			fmt.Printf("Error %v", in.err)
			success = false
		}
		outResult(output, out)
	default:
		// fmt.Fprintf(os.Stderr, "mode: default\n")
		// default: generate config file from input + template
		if tmpl == "" {
			log.Error.Printf("Field [template(-t|--tmpl)] required")
			os.Exit(1)
		}

		tmplSrc, err := getInput(tmpl)
		if err != nil {
			log.Error.Printf("Error %v", tmplSrc.err)
			os.Exit(1)
		}
		template, err := junocfg.CheckTemplate(tmplSrc.input[0])
		if err != nil {
			log.Error.Printf("check tmpl error %v", err)
			os.Exit(1)
		}

		tmplSrc.dump()

		in, err := getInput(input)
		if err != nil {
			log.Error.Printf("error %v", in.err)
			os.Exit(1)
		}
		in.dump()
		settingsData, err := junocfg.MergeYamls(in.input)
		if err != nil {
			log.Error.Printf("error %v", err)
			os.Exit(1)
		}

		settings, err := junocfg.UnmarshalYaml(settingsData)
		if err != nil {
			log.Error.Printf("error %v", err)
			os.Exit(1)
		}

		buffer := bytes.NewBuffer([]byte{})

		if err := template.Execute(buffer, settings); err != nil {
			log.Error.Printf("failed to render template [%s]\n[%s]\n", err, settings)
			os.Exit(1)
		}

		buffer = junocfg.PreprocessYaml(buffer)
		// check yaml
		if err = junocfg.CheckYaml(buffer.Bytes()); err != nil {
			log.Stderr.Printf("Not valid output yaml: %s\n", err.Error())
			os.Exit(1)
		}

		// check variables
		strOut := strings.Split(buffer.String(), "\n")
		for posInFile, str := range strOut {
			if i := strings.Index(str, "<no value>"); i != -1 {
				log.Stderr.Printf("<no value> at %s#%d:%s\n", output, posInFile, str)
			}
		}
		junocfg.OutResult(output, buffer)
	}
	if !success {
		os.Exit(1)
	}
}
