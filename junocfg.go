package junocfg

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

func readDataFromScanner(scanner *bufio.Scanner) []string {
	result := []string{}
	for scanner.Scan() { // internally, it advances token based on sperator
		line := scanner.Text()
		result = append(result, line)
	}
	return result
}

func ReadData(filename string) (string, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	if fi.Mode()&os.ModeNamedPipe == 0 {
		if filename == "" {
			// log.Printf("no pipe, no file :(")
			return "", fmt.Errorf("empty input")
		} else {
			// log.Printf("read data from file")
			fd, err := os.Open(filename)
			if err != nil {
				return "", fmt.Errorf("open file: reading error")
			}
			defer fd.Close()

			scanner = bufio.NewScanner(fd)
		}
	} else {
		if filename == "" {
			// log.Printf("pipe!\n")
		} else {
			// log.Printf("pipe, file skipped\n")
		}
	}

	result := readDataFromScanner(scanner)
	if err != nil {
		return "", fmt.Errorf("readDataFromScanner: reading error")
	}
	return strings.Join(result, "\n"), nil
}

func GetConfigs(filenames []string) (map[string]interface{}, error) {
	var rawConfigs []string

	for _, filename := range filenames {
		if data, err := ReadData(filename); err != nil {
			return nil, fmt.Errorf("Could not read input: %s", err)
		} else {
			rawConfigs = append(rawConfigs, data)
		}
	}

	cfg := make(map[string]interface{})
	for _, buffer := range rawConfigs {
		tmpCfg := make(map[string]interface{})
		if err := yaml.Unmarshal([]byte(buffer), tmpCfg); err != nil {
			return nil, fmt.Errorf("Could not parse YAML file: %s\n", err)
		}
		cfg = MergeMaps(cfg, tmpCfg)
	}
	return cfg, nil
}

func GetTemplate(filename string) (*template.Template, error) {
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

func getYamlTemplate(filename string) (map[string]interface{}, error) {
	content, err := loadFile(filename)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Template file load error: [%v]\n", err))
	}

	tmplYaml := map[string]interface{}{}
	err = yaml.Unmarshal(content, &tmplYaml)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot unmarshal yaml.tmpl: %s\n", err.Error())
		fmt.Fprintf(os.Stderr, "yaml:\n%s", string(content))
		log.Fatal("...")
	}
	return tmplYaml, nil
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
		config = MergeMaps(config, cfg)
	}

	return config, nil
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
