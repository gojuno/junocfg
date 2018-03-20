package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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
	yaml2json bool

	verbose bool

	input  string
	output string
	tmpl   string
)

func init() {
	flag.BoolVar(&checkTmpl, "check-tmpl", false, "check tmpl")
	flag.BoolVar(&merge, "merge", false, "merge")
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.BoolVar(&yaml2json, "yaml2json", false, "yaml2json")

	flag.StringVar(&input, "i", "", "input")
	flag.StringVar(&input, "input", "", "input")
	flag.StringVar(&output, "o", "", "output")
	flag.StringVar(&output, "output", "", "output")
	flag.StringVar(&tmpl, "t", "", "template")
	flag.StringVar(&tmpl, "template", "", "template")
}

func checkFatal(errStr string, err error) {
	if err != nil {
		log.Stderr.Printf("ARGS: checkTmpl: [%v] merge: [%v] yaml2json: [%v] verbose: [%v] input: [%v] output: [%v] tmpl: [%v]\n",
			checkTmpl, merge, yaml2json, verbose,
			input, output, tmpl,
		)
		log.Stderr.Printf(errStr, err)
		os.Exit(1)
	}
}

func initLogger() {
	logger := &log.Logger{
		os.Stderr, // ioutil.Discard,
		os.Stderr, // ioutil.Discard,
		os.Stderr,
		os.Stderr,
		os.Stderr,
	}
	log.InitLoggers(logger)
	log.Stdout = stllog.New(os.Stdout, "", 0)
	log.Stderr = stllog.New(os.Stderr, "", 0)

	if !verbose {
		log.Debug = stllog.New(ioutil.Discard, "", 0)
	}
}

func main() {
	flag.Parse()
	initLogger()

	var settingsInput *inputData
	var tmplInput *inputData
	var err error

	if !checkTmpl {
		settingsInput, err = getInput(input)
		settingsInput.dump()
		checkFatal("Error %v", err)
	}
	if !yaml2json && !merge {
		tmplInput, err = getInput(tmpl)
		tmplInput.dump()
		checkFatal("Error %v", err)
	}

	switch {
	case checkTmpl:
		_, err = junocfg.CheckTemplate(tmplInput.input[0])
		checkFatal("check tmpl error %v", err)
	case yaml2json:
		maps, err := junocfg.Yamls2Maps(settingsInput.input)
		checkFatal("Yamls2Maps error %v", err)

		resultMap, err := junocfg.MergeMaps(maps)
		checkFatal("MergeMaps error %v", err)

		out, err := junocfg.Map2Json(resultMap)
		checkFatal("Map2Json error %v", err)

		outResult(output, string(out))
	case merge:
		maps, err := junocfg.Yamls2Maps(settingsInput.input)
		checkFatal("Yamls2Maps error %v", err)

		resultMap, err := junocfg.MergeMaps(maps)
		checkFatal("MergeMaps error %v", err)

		out, err := junocfg.Map2Yaml(resultMap)
		checkFatal("Map2Yaml error %v", err)

		outResult(output, string(out))
	default:
		// fmt.Fprintf(os.Stderr, "mode: default\n")
		// default: generate config file from input + template
		if tmpl == "" {
			checkFatal("%v", fmt.Errorf("Field [template(-t|--tmpl)] required"))
		}

		maps, err := junocfg.Yamls2Maps(settingsInput.input)
		checkFatal("Yamls2Maps error %v", err)

		settingsMap, err := junocfg.MergeMaps(maps)
		checkFatal("MergeMaps error %v", err)

		// settings, err := junocfg.Map2Yaml(resultMap)
		// checkFatal("Map2Yaml error %v", err)

		out, err := junocfg.RenderAndCheckTemplate(tmplInput.input[0], settingsMap)
		checkFatal("RenderAndCheckTemplate error %v", err)

		// check variables
		strOut := strings.Split(out, "\n")
		for posInFile, str := range strOut {
			if i := strings.Index(str, "<no value>"); i != -1 {
				log.Stderr.Printf("<no value> at %s#%d:%s\n", output, posInFile, str)
			}
		}
		outResult(output, out)
	}
}
