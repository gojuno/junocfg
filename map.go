package junocfg

import (
	"fmt"
	"reflect"
)

func setValue(src map[string]interface{}, path []string, value interface{}) error {
	ptr := src
	for i, key := range path {
		//fmt.Printf("%d %v\n", i, key)
		v, ok := ptr[key]
		if !ok {
			ptr[key] = map[string]interface{}{}
			v, ok = ptr[key]
		}

		if i != len(path)-1 {
			value := reflect.ValueOf(v)
			if value.Kind() != reflect.Map {
				ptr[key] = map[string]interface{}{}
			}
			if ptr, ok = ptr[key].(map[string]interface{}); !ok {
				return fmt.Errorf("setValue error: %v", key)
			}
			continue
		}

		ptr[key] = value
		break
	}
	return nil
}

func Items2Map(items ItemArray) (map[string]interface{}, error) {
	dst := map[string]interface{}{}
	for _, i := range items {
		if err := setValue(dst, i.path, i.value); err != nil {
			return nil, err
		}
	}
	result := deinterface(dst)
	return result, nil
}
