package utils

import (
	"fmt"
	"reflect"
)

func ConvertYAMLToJSON(yamlData interface{}) (interface{}, error) {
	switch reflect.TypeOf(yamlData).Kind() {
	case reflect.Map:
		jsonMap := make(map[string]interface{})
		m, ok := yamlData.(map[string]interface{})
		if ok {
			for key, value := range m {
				jsonValue, err := ConvertYAMLToJSON(value)
				if err != nil {
					return nil, err
				}
				jsonMap[key] = jsonValue
			}
		} else {
			for key, value := range yamlData.(map[interface{}]interface{}) {
				strKey, ok := key.(string)
				if !ok {
					return nil, fmt.Errorf("unsupported type: %T for map key", key)
				}
				jsonValue, err := ConvertYAMLToJSON(value)
				if err != nil {
					return nil, err
				}
				jsonMap[strKey] = jsonValue
			}
		}

		return jsonMap, nil
	default:
		return yamlData, nil
	}
}
