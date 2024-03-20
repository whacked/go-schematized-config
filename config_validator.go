package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

// LoadJSON attempts to load JSON data from a string or map. If the input is a string, it treats it as a filepath.
func LoadJSON(jsonSource interface{}) (map[string]interface{}, error) {
	switch js := jsonSource.(type) {
	case string:
		data, err := os.ReadFile(js)
		if err != nil {
			return nil, err
		}
		var result map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, err
		}
		return result, nil
	case map[string]interface{}:
		return js, nil
	default:
		return nil, errors.New("unsupported type for jsonSource")
	}
}

// ValidateConfig validates the provided configuration map against the given JSON schema.
func ValidateConfig(schemaMap map[string]interface{}, config map[string]interface{}) error {
	schemaBytes, err := json.Marshal(schemaMap)
	if err != nil {
		return err
	}
	schema, err := jsonschema.CompileString("schema", string(schemaBytes))
	if err != nil {
		return err
	}
	if err := schema.Validate(config); err != nil {
		return err
	}
	return nil
}

func convertMapInterfaceToString(data map[string]interface{}) (map[string]string, error) {
	result := make(map[string]string)
	for k, v := range data {
		switch value := v.(type) {
		case string:
			result[k] = value
		case int, int32, int64, float32, float64, bool:
			// Use fmt.Sprintf to convert basic types to their string representations
			result[k] = fmt.Sprintf("%v", value)
		default:
			return nil, fmt.Errorf("unsupported type for key %s: %T", k, v)
		}
	}
	return result, nil
}

// LoadValidatedConfig loads and validates a configuration map against a JSON schema.
// It extracts declared items, coerces config values, and validates the config.
func LoadValidatedConfig(schema map[string]interface{}, config map[string]interface{}) (map[string]interface{}, error) {
	// Assume ExtractDeclaredItems and CoercePrimitiveValues are implemented and work similarly to their Python counterparts
	extractedConfig, err := ExtractDeclaredItems(schema, config)
	if err != nil {
		return nil, fmt.Errorf("error extracting declared items: %w", err)
	}

	extractStringMap, err := convertMapInterfaceToString(extractedConfig)
	if err != nil {
		return nil, fmt.Errorf("error converting map to string: %w", err)
	}
	coercedConfig, err := CoercePrimitiveValues(schema, extractStringMap)
	if err != nil {
		return nil, fmt.Errorf("error coercing config values: %w", err)
	}

	// Validate the coerced configuration against the JSON schema
	err = ValidateConfig(schema, coercedConfig)
	if err != nil {
		return nil, err
	}

	return coercedConfig, nil
}
