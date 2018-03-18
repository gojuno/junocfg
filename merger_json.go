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
