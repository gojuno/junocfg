package junocfg

import (
	"fmt"
	"reflect"
	"sort"
)

func walker(data map[string]interface{}, path []string, handler func([]string, interface{}) error) error {
	for key := range data {
		ptr := append(path, key)
		value := reflect.ValueOf(data[key])
		if value.Kind() == reflect.Map {
			v := data[key]
			vi, ok := v.(map[string]interface{})
			if !ok {
				vi = make(map[string]interface{})
				for k, vv := range v.(map[interface{}]interface{}) {
					vi[fmt.Sprintf("%s", k)] = vv
				}
			}
			walker(vi, ptr, handler)
		} else {
			handler(ptr, data[key])
		}
	}
	return nil
}

func walk(data map[string]interface{}) ([]item, error) {
	result := itemArray{}
	f := func(path []string, value interface{}) error {
		result = append(result, item{path, value})
		return nil
	}
	walker(data, []string{}, f)
	sort.Sort(result)
	return result, nil
}
