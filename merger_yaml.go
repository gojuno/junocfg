package junocfg

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

func MergeYamls(data [][]byte) ([]byte, error) {
	m := map[string]interface{}{}

	for i, d := range data {
		err := yaml.Unmarshal(d, &m)
		if err != nil {
			return nil, fmt.Errorf("unmarshal %d batch error: %v", i, err)
		}
	}
	out, err := yaml.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %v", err)
	}

	return out, err
}

func CheckYaml(data []byte) error {
	y := map[string]interface{}{}
	///fmt.Printf("[[\n%s]]", string(data))
	return yaml.Unmarshal(data, &y)
}

func UnmarshalYaml(data []byte) (map[string]interface{}, error) {
	y := map[string]interface{}{}
	err := yaml.Unmarshal(data, &y)
	return y, err
}
