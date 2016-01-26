package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func loadFile(filename string) ([]byte, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte(""), err
	}
	return content, nil
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

func outDict(dict map[string]interface{}) {
	d, err := yaml.Marshal(dict)
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not create YAML file: %s\n", err))
	}

	fmt.Fprintf(os.Stderr, "data: <!--\n%s\n-->\n", d)
}

func mergeMaps(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		dest[k] = v
	}
	return dest
}
