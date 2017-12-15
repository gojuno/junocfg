package junocfg

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"strings"
)

func MergeYamls(data [][]byte) ([]byte, error) {
	result := map[string]interface{}{}

	for i, d := range data {
		tmpMap := map[string]interface{}{}
		err := yaml.Unmarshal(d, &tmpMap)
		if err != nil {
			return nil, fmt.Errorf("unmarshal %d batch error: %v", i, err)
		}
		if err := mergeMaps(tmpMap, result); err != nil {
			return nil, fmt.Errorf("merge map error: %v", err)
		}
	}

	out, err := yaml.Marshal(result)
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
