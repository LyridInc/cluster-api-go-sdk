package utils

import (
	"encoding/json"
	"fmt"
)

func ConvertYAMLToJSON(yamlData interface{}) ([]byte, error) {
	switch v := yamlData.(type) {
	case map[interface{}]interface{}:
		jsonMap := make(map[string]interface{})
		for key, value := range v {
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
		return json.Marshal(jsonMap)
	case []interface{}:
		jsonSlice := make([]interface{}, len(v))
		for i, item := range v {
			jsonValue, err := ConvertYAMLToJSON(item)
			if err != nil {
				return nil, err
			}
			jsonSlice[i] = jsonValue
		}
		return json.Marshal(jsonSlice)
	default:
		return json.Marshal(v)
	}
}
