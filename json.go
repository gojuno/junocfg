package junocfg

import (
	"encoding/json"
	"fmt"
)

func mergeJsons(data [][]byte) ([]byte, error) {
	m := map[string]interface{}{}

	for i, d := range data {
		err := json.Unmarshal(d, &m)
		if err != nil {
			return nil, fmt.Errorf("unmarshal %d batch error: %v", i, err)
		}
	}
	fmt.Printf("%+v", m)
	out, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %v", err)
	}

	return out, err
}

func Jsons2Maps(data [][]byte) ([]map[string]interface{}, error) {
	result := []map[string]interface{}{}

	for i, d := range data {
		jsonmap, err := Json2Map(d)
		if err != nil {
			return nil, fmt.Errorf("unmarshal %d batch error: %v", i, err)
		}
		result = append(result, jsonmap)
	}
	return result, nil
}

func Json2Map(data []byte) (map[string]interface{}, error) {
	rawJsonMap := map[string]interface{}{}
	err := json.Unmarshal(data, &rawJsonMap)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	// convert map[interface{}]interface{} -> map[string]interface{}
	tmpMap := map[string]interface{}{}
	if err := catMaps(rawJsonMap, tmpMap); err != nil {
		return nil, fmt.Errorf("merge map error: %v", err)
	}
	return tmpMap, nil
}

func Map2Json(data map[string]interface{}) ([]byte, error) {
	out, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %v", err)
	}
	return out, err
}
