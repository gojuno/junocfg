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
	json2yaml bool

	verbose bool

	input  string
	output string
	tmpl   string
)

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose")

	flag.BoolVar(&checkTmpl, "check-tmpl", false, "check tmpl")

	flag.BoolVar(&merge, "merge", false, "*.yaml -> yaml")
	flag.BoolVar(&yaml2json, "yaml2json", false, "*.yaml -> json")
	flag.BoolVar(&json2yaml, "json2yaml", false, "*.json -> yaml")

	flag.StringVar(&input, "i", "", "input")
	flag.StringVar(&input, "input", "", "input")
	flag.StringVar(&output, "o", "", "output")
	flag.StringVar(&output, "output", "", "output")
	flag.StringVar(&tmpl, "t", "", "template")
	flag.StringVar(&tmpl, "template", "", "template")
}

func checkFatal(errStr string, err error) {
	if err != nil {
		log.Debug.Printf("ARGS:\n\tverbose: [%v]\n\tcheckTmpl: [%v]\n\tmerge: [%v] yaml2json: [%v] json2yaml: [%v]\n\tinput: [%v] output: [%v] tmpl: [%v]\n",
			verbose,
			checkTmpl, merge, yaml2json, json2yaml,
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
	var items junocfg.ItemArray
	var settingsMap map[string]interface{}
	var err error

	// read input
	if !checkTmpl {
		settingsInput, err = getInput(input)
		settingsInput.dump()
		checkFatal("Error %v", err)
	}
	if !yaml2json && !json2yaml && !merge {
		tmplInput, err = getInput(tmpl)
		tmplInput.dump()
		checkFatal("Error %v", err)
	}

	// prepare input - merge multiply files to single map
	if checkTmpl {

	} else if json2yaml {
		items, err = junocfg.Jsons2Items(settingsInput.input)
		checkFatal("Jsons2Items error %v", err)
	} else { // yaml2json, merge, default
		items, err = junocfg.Yamls2Items(settingsInput.input)
		checkFatal("Yamls2Items error %v", err)
	}

	if !checkTmpl {
		settingsMap, err = junocfg.Items2Map(items)
		checkFatal("Items2Map error %v", err)
	}

	// process
	switch {
	case checkTmpl:
		_, err = junocfg.CheckTemplate(tmplInput.input[0])
		checkFatal("check tmpl error %v", err)
	case yaml2json:
		out, err := junocfg.Map2Json(settingsMap)
		checkFatal("Map2Json error %v", err)

		outResult(output, string(out))
	case json2yaml:
		out, err := junocfg.Map2Yaml(settingsMap)
		checkFatal("Map2Yaml error %v", err)

		outResult(output, string(out))
	case merge:
		out, err := junocfg.Map2Yaml(settingsMap)
		checkFatal("Map2Yaml error %v", err)

		outResult(output, string(out))
	default:
		// fmt.Fprintf(os.Stderr, "mode: default\n")
		// default: generate config file from input + template
		if tmpl == "" {
			checkFatal("%v", fmt.Errorf("Field [template(-t|--tmpl)] required"))
		}

		out, err := junocfg.RenderAndCheckTemplate(tmplInput.input[0], settingsMap)
		checkFatal("RenderAndCheckTemplate error %v", err)

		// check variables
		errorsCount := 0
		strOut := strings.Split(out, "\n")
		for posInFile, str := range strOut {
			if i := strings.Index(str, "<no value>"); i != -1 {
				log.Stderr.Printf("<no value> at %s#%d:%s\n", output, posInFile, str)
				errorsCount++
			}
		}
		outResult(output, out)
		if errorsCount > 0 {
			checkFatal("Empty variables %v", fmt.Errorf("%d", errorsCount))
		}
	}
}
