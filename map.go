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

func catMaps(src map[string]interface{}, dst map[string]interface{}) error {
	if items, err := walk(src); err != nil {
		return err
	} else {
		for _, i := range items {
			if err := setValue(dst, i.path, i.value); err != nil {
				return err
			}
		}
	}
	return nil
}

func MergeMaps(data []map[string]interface{}) (map[string]interface{}, error) {
	dst := map[string]interface{}{}
	for _, d := range data {
		if err := catMaps(d, dst); err != nil {
			return nil, err
		}
	}
	return dst, nil
}

func mutator(data map[string]interface{}, result map[string]interface{}) error {
	for key := range data {
		value := reflect.ValueOf(data[key])
		switch value.Kind() {
		case reflect.Map:
			result[key] = make(map[string]interface{})
			mutator(data[key].(map[string]interface{}), result[key].(map[string]interface{}))
		case reflect.Slice:
			result[key] = []interface{}{}
			for idx, element := range data[key].([]interface{}) {
				v := reflect.ValueOf(element)
				switch v.Kind() {
				case reflect.Map:
					e := make(map[string]interface{})
					mutator(data[key].(map[string]interface{}), result[key].(map[string]interface{}))
				case reflect.Slice:
					result[key] = []interface{}{}
					for idx, v := range data[key].([]interface{}) {

					}
				default:
					result[key] = data[key]
				}
			}
		default:
			result[key] = data[key]
		}
	}
	return nil
}

func mutate(data map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	mutator(data, result)
	return result, nil
}
