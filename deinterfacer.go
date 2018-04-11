package junocfg

import (
	"fmt"
	"reflect"
)

// deinterfacer converts `map[interface{}]interface{}` to `map[string]interface{}` recursively.
func deinterfacer(data interface{}) interface{} {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Map:
		result := make(map[string]interface{})
		for _, k := range value.MapKeys() {
			result[fmt.Sprintf("%v", k)] = deinterfacer(value.MapIndex(k).Interface())
		}

		return result
	case reflect.Slice:
		result := make([]interface{}, value.Len())
		for i := 0; i < value.Len(); i++ {
			result[i] = deinterfacer(value.Index(i).Interface())
		}
		return result
	default:
		return data
	}
}

func deinterface(data map[string]interface{}) map[string]interface{} {
	return deinterfacer(data).(map[string]interface{})
}
