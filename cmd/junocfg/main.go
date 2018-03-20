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

func checkFatal(errStr string, err error) {
	if err != nil {
		log.Stderr.Printf("ARGS: checkTmpl: [%v] merge: [%v] input: [%v] output: [%v] tmpl: [%v]\n",
			checkTmpl, merge,
			input, output, tmpl,
		)
		log.Stderr.Printf(errStr, err)
		os.Exit(1)
	}
}

func main() {
	// junocfg.Test()
	// os.Exit(1)

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

	switch {
	case checkTmpl:
		// fmt.Fprintf(os.Stderr, "mode: check-tmpl\n")
		in, err := getInput(tmpl)
		log.Debug.Printf(in.dump())
		checkFatal("Error %v", in.err)
		_, err = junocfg.CheckTemplate(in.input[0])
		checkFatal("check tmpl error %v", err)
	case merge:
		// fmt.Fprintf(os.Stderr, "mode: merge\n")
		in, err := getInput(input)
		checkFatal("Error %v", err)
		out, err := junocfg.MergeYamls(in.input)
		// in.dump()
		checkFatal("Error %v", in.err)
		outResult(output, out)
	default:
		// fmt.Fprintf(os.Stderr, "mode: default\n")
		// default: generate config file from input + template
		if tmpl == "" {
			checkFatal("%v", fmt.Errorf("Field [template(-t|--tmpl)] required"))
		}

		tmplSrc, err := getInput(tmpl)
		checkFatal("Error %v", tmplSrc.err)

		template, err := junocfg.CheckTemplate(tmplSrc.input[0])
		checkFatal("check tmpl error %v", err)

		tmplSrc.dump()

		in, err := getInput(input)
		checkFatal("error %v", in.err)
		// in.dump()
		settingsData, err := junocfg.MergeYamls(in.input)
		checkFatal("error %v", err)

		settings, err := junocfg.UnmarshalYaml(settingsData)
		checkFatal("error %v", err)

		buffer := bytes.NewBuffer([]byte{})

		checkFatal("failed to render template [%s]\n", template.Execute(buffer, settings))

		buffer = junocfg.PreprocessYaml(buffer)
		// check yaml
		if err = junocfg.CheckYaml(buffer.Bytes()); err != nil {
			log.Stderr.Printf("Not valid output yaml: %s\n", err.Error())
			log.Stderr.Printf("%s\n", buffer.Bytes())
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
}
