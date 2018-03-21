package junocfg

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

func yaml2Items(data []byte) (ItemArray, error) {
	yamlMap := map[string]interface{}{}
	err := yaml.Unmarshal(data, &yamlMap)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}
	items, err := walk(yamlMap)
	return items, err
}

func Yamls2Items(data [][]byte) (ItemArray, error) {
	result := ItemArray{}
	for i, d := range data {
		items, err := yaml2Items(d)
		if err != nil {
			return nil, fmt.Errorf("unmarshal %d batch error: %v", i, err)
		}
		result = append(result, items...)
	}
	return result, nil
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
