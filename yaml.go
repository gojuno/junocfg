package junocfg

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

func Yamls2Maps(data [][]byte) ([]map[string]interface{}, error) {
	result := []map[string]interface{}{}

	for i, d := range data {
		yamlmap, err := Yaml2Map(d)
		if err != nil {
			return nil, fmt.Errorf("unmarshal %d batch error: %v", i, err)
		}
		result = append(result, yamlmap)
	}
	return result, nil
}

func Yaml2Map(data []byte) (map[string]interface{}, error) {
	rawYamlMap := map[string]interface{}{}
	err := yaml.Unmarshal(data, &rawYamlMap)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	// convert map[interface{}]interface{} -> map[string]interface{}
	tmpMap := map[string]interface{}{}
	if err := catMaps(rawYamlMap, tmpMap); err != nil {
		return nil, fmt.Errorf("merge map error: %v", err)
	}
	return tmpMap, nil
}

func Map2Yaml(data map[string]interface{}) ([]byte, error) {
	out, err := yaml.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %v", err)
	}
	return out, err
}

func CheckYaml(data []byte) error {
	y := map[string]interface{}{}
	return yaml.Unmarshal(data, &y)
}

func UnmarshalYaml(data []byte) (map[string]interface{}, error) {
	y := map[string]interface{}{}
	err := yaml.Unmarshal(data, &y)
	return y, err
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

func writeStr(buffer *bytes.Buffer, str string) {
	if strings.TrimSpace(str) != "" {
		buffer.WriteString(str)
		buffer.WriteString("\n")
	}
}
