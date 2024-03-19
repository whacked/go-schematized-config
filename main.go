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

// ExtractDeclaredItems filters the input data map based on the keys declared in the JSON schema.
// It removes keys not declared in the schema and adds keys with their default values if they are
// declared in the schema but not present in the input data.
func ExtractDeclaredItems(jsonSchema map[string]interface{}, data map[string]interface{}) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	properties, ok := jsonSchema["properties"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected 'properties' to be a map[string]interface{}")
	}

	// Keep only the keys that are declared in the schema
	for key, value := range data {
		if _, exists := properties[key]; exists {
			out[key] = value
		}
	}

	// Add defaults for any declared keys missing from the data
	for key, propInterface := range properties {
		prop, ok := propInterface.(map[string]interface{})
		if !ok {
			continue // or return an error if strict validation is needed
		}
		if _, exists := out[key]; !exists {
			if defaultValue, hasDefault := prop["default"]; hasDefault {
				out[key] = defaultValue
			}
		}
	}

	return out, nil
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
