package junocfg

import (
	"fmt"
	"reflect"
)

// convert `map[interface{}]interface{}` to `map[string]interface{}`
func deinterfacer(data interface{}) (interface{}, error) {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Map:
		result := make(map[string]interface{})
		_, ok := data.(map[string]interface{})
		if !ok {
			_, ok := data.(map[interface{}]interface{})
			if !ok {
				//
			} else {
				for key, value := range data.(map[interface{}]interface{}) {
					if v, err := deinterfacer(value); err != nil {
						//
					} else {
						result[fmt.Sprintf("%s", key)] = v
					}
				}
			}
		} else {
			for key, value := range data.(map[string]interface{}) {
				if v, err := deinterfacer(value); err != nil {
					//
				} else {
					result[fmt.Sprintf("%s", key)] = v
				}
			}
		}
		return result, nil
	case reflect.Slice:
		s := reflect.ValueOf(data)
		result := make([]interface{}, s.Len())
		for i := 0; i < s.Len(); i++ {
			value := s.Index(i).Interface()
			if v, err := deinterfacer(value); err != nil {
				//
			} else {
				result[i] = v
			}
		}
		return result, nil
	default:
		return data, nil
	}
	return nil, nil
}

func deinterface(data map[string]interface{}) (map[string]interface{}, error) {
	result, err := deinterfacer(data)
	return result.(map[string]interface{}), err
}
