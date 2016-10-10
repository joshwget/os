package config

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func convertMapKeysToStrings(existingMap map[string]interface{}) map[string]interface{} {
	newMap := make(map[string]interface{})

	for k, v := range existingMap {
		newMap[k] = convertKeysToStrings(v)
	}

	return newMap
}

func convertKeysToStrings(item interface{}) interface{} {
	switch typedDatas := item.(type) {

	case map[interface{}]interface{}:
		newMap := make(map[string]interface{})

		for key, value := range typedDatas {
			stringKey := key.(string)
			newMap[stringKey] = convertKeysToStrings(value)
		}
		return newMap

	case []interface{}:
		// newArray := make([]interface{}, 0) will cause golint to complain
		var newArray []interface{}
		newArray = make([]interface{}, 0)

		for _, value := range typedDatas {
			newArray = append(newArray, convertKeysToStrings(value))
		}
		return newArray

	default:
		return item
	}
}

func Validate(rawCfg map[string]interface{}) (*gojsonschema.Result, error) {
	data := convertMapKeysToStrings(rawCfg)
	loader := gojsonschema.NewGoLoader(data)
	schemaLoader := gojsonschema.NewStringLoader(schema)
	_ = loader
	_ = schemaLoader
	fmt.Println(data)
	return nil, nil
	//return gojsonschema.Validate(schemaLoader, loader)
}
