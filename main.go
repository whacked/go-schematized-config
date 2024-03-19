package main

import (
	"fmt"
	"strconv"
	"strings"
)

// CoercePrimitiveValues converts string values in a map to their respective types based on a provided JSON schema.
func CoercePrimitiveValues(jsonSchema map[string]interface{}, data map[string]string) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	// copy first to preserve un-coerced values
	for k, v := range data {
		out[k] = v
	}

	properties, ok := jsonSchema["properties"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid json schema; `properties` not a map[string]interface{}")
	}

	for propertyName, propertySchema := range properties {
		propertySchemaMap, ok := propertySchema.(map[string]interface{})
		if !ok {
			continue // or return an error if strict validation is needed
		}
		propertyType, ok := propertySchemaMap["type"].(string)
		if !ok {
			continue // or return an error
		}
		propertyValue, exists := data[propertyName]
		if !exists {
			continue
		}

		var err error
		switch propertyType {
		case "integer":
			out[propertyName], err = strconv.Atoi(propertyValue)
			if err != nil {
				return nil, fmt.Errorf("error converting property %s to integer: %v", propertyName, err)
			}
		case "number":
			out[propertyName], err = strconv.ParseFloat(propertyValue, 64)
			if err != nil {
				return nil, fmt.Errorf("error converting property %s to float: %v", propertyName, err)
			}
		case "boolean":
			out[propertyName], err = parseBool(propertyValue)
			if err != nil {
				return nil, fmt.Errorf("error converting property %s to boolean: %v", propertyName, err)
			}
		}
	}

	return out, nil
}

// parseBool converts string to boolean.
func parseBool(value string) (bool, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "true", "yes", "y", "1", "on":
		return true, nil
	case "false", "no", "n", "0", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean value: %s", value)
	}
}

func main() {
	// Example usage
	jsonSchema := map[string]interface{}{
		"properties": map[string]interface{}{
			"age": map[string]interface{}{
				"type": "integer",
			},
			"active": map[string]interface{}{
				"type": "boolean",
			},
		},
	}
	data := map[string]string{
		"age":    "30",
		"active": "true",
	}

	coercedData, err := CoercePrimitiveValues(jsonSchema, data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(coercedData)
}
