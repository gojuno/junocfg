package junocfg

import (
	"encoding/json"
	"fmt"
)

func json2Items(data []byte) (ItemArray, error) {
	jsonMap := map[string]interface{}{}
	err := json.Unmarshal(data, &jsonMap)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}
	items, err := walk(jsonMap)
	return items, err
}

func Jsons2Items(data [][]byte) (ItemArray, error) {
	result := ItemArray{}
	for i, d := range data {
		items, err := json2Items(d)
		if err != nil {
			return nil, fmt.Errorf("unmarshal %d batch error: %v", i, err)
		}
		result = append(result, items...)
	}
	return result, nil
}

func Map2Json(data map[string]interface{}) ([]byte, error) {
	out, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %v", err)
	}
	return out, err
}
