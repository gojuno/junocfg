package junocfg

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

func getYamlTemplate(filename string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Template file load error: [%v]\n", err))
	}

	tmplYaml := map[string]interface{}{}
	err = yaml.Unmarshal(data, &tmplYaml)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot unmarshal yaml.tmpl: %s\n", err.Error())
		fmt.Fprintf(os.Stderr, "yaml:\n%s", string(data))
		log.Fatal("...")
	}
	return tmplYaml, nil
}

func writeStr(buffer *bytes.Buffer, str string) {
	if strings.TrimSpace(str) != "" {
		buffer.WriteString(str)
		buffer.WriteString("\n")
	}
}

func PreprocessYaml(input *bytes.Buffer) *bytes.Buffer {
	buffer := bytes.NewBuffer([]byte{})

	ident := ""
	for _, str := range strings.Split(input.String(), "\n") {
		if strings.HasSuffix(str, "|") {
			ident = "|"
		} else if ident == "|" {
			count := len(str) - len(strings.TrimLeft(str, " "))
			// fmt.Println(str, len(str), len(strings.TrimLeft(str, " ")), count)
			ident = strings.Repeat(" ", count)
			writeStr(buffer, str)
			continue
		} else if ident != "" && (strings.HasPrefix(str, " ") || str == "") {
			ident = ""
		} else {
		}
		if ident != "|" && ident != "" {
			buffer.WriteString(ident)
		}
		writeStr(buffer, str)
	}
	return buffer
}

func OutResult(filename string, buffer *bytes.Buffer) {
	outputBuffer := bufio.NewWriter(os.Stdout)
	if filename != "" {
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

type cfgItem struct {
	path  string
	value interface{}
}

func appendMap(dest map[string]interface{}, item cfgItem) map[string]interface{} {
	path := strings.Split(item.path, "\t")[1:]

	var cursor interface{}
	cursor = dest

	for _, p := range path[:len(path)-1] {
		if ms, ok := cursor.(map[string]interface{}); ok {
			if _, ok := ms[p]; !ok {
				ms[p] = map[string]interface{}{}
			}
			cursor = ms[p]
		} else {
			panic("ms !ok")
		}
	}
	if ms, ok := cursor.(map[string]interface{}); ok {
		ms[path[len(path)-1]] = item.value
	} else {
		panic("ms[p] !ok")
	}

	return dest
}

func map2list(src map[interface{}]interface{}, srcPath string, cfg *[]cfgItem) *[]cfgItem {
	for k, v := range src {
		key := k.(string)
		path := srcPath + "\t" + key

		var item *cfgItem
		switch t := v.(type) {
		case map[interface{}]interface{}:
			cfg = map2list(t, path, cfg)
		//feel free to add as many supported types if you want
		case []interface{}, string, int, int64, bool, float64, float32:
			item = &cfgItem{path: path, value: t}
		default:
			log.Fatalf("map2list: unexpected type of the key %q with value %v found in yaml file\n", strings.Replace(path, "\t", ".", -1), v)
		}

		if item != nil {
			*cfg = append(*cfg, *item)
		}
	}
	return cfg
}

func MergeMaps(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	tmp := map[interface{}]interface{}{}
	for k, v := range src {
		tmp[k] = v
	}

	cfg := map2list(tmp, "", new([]cfgItem))
	for _, item := range *cfg {
		dest = appendMap(dest, item)
	}

	return dest
}
