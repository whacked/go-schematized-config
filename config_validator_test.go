package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadAndValidateConfig(t *testing.T) {
	schemaJsonnetPath := filepath.Join("testdata", "test-schema-1.schema.jsonnet")
	schema, err := LoadJsonnetFile(schemaJsonnetPath)
	if err != nil {
		t.Fatalf("Failed to load jsonnet: %v", err)
	}

	// Define your test cases
	tests := []struct {
		name     string
		config   map[string]interface{}
		expect   map[string]interface{}
		expectOK bool
	}{
		{
			name:     "Empty config should fail",
			config:   map[string]interface{}{},
			expect:   nil,
			expectOK: false,
		},
		{
			name: "Valid config with defaults",
			config: map[string]interface{}{
				"string_value_with_enum": "these",
				"MY_INTEGER_VALUE":       "1122334",
				"A_NUMERIC_VALUE":        "24.89",
			},
			expect: map[string]interface{}{
				"string_value_with_enum":                "these",
				"MY_INTEGER_VALUE":                      1122334,
				"A_NUMERIC_VALUE":                       24.89,
				"_____A_STRING_VALUE____with_default__": "underscores_and spaces",
			},
			expectOK: true,
		},
		{
			name: "Valid config with scientific notation and negative integer",
			config: map[string]interface{}{
				"string_value_with_enum": "these",
				"MY_INTEGER_VALUE":       "-85",
				"A_NUMERIC_VALUE":        "1.23e4",
			},
			expect: map[string]interface{}{
				"string_value_with_enum":                "these",
				"MY_INTEGER_VALUE":                      -85,
				"A_NUMERIC_VALUE":                       12300.0,
				"_____A_STRING_VALUE____with_default__": "underscores_and spaces",
			},
			expectOK: true,
		},
	}

	// Iterate over test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Here you would call your validation function
			// For demonstration, let's assume it returns an error if validation fails
			validatedConfig, err := LoadValidatedConfig(schema, tc.config)

			// Check if the validation result matches the expected outcome
			if (err == nil) != tc.expectOK {
				t.Errorf("Expected validation result %v, got error: %v", tc.expectOK, err)
			}

			// If the validation result is expected to be successful, check the result
			if tc.expectOK {
				if !reflect.DeepEqual(validatedConfig, tc.expect) {
					t.Errorf("LoadValidatedConfig() got = %v, want %v", validatedConfig, tc.expect)
				}
			}

			// (visual check) pretty print result
			if !false {
				prettyConfig, err := json.MarshalIndent(validatedConfig, "", "  ")
				if err != nil {
					t.Fatalf("Failed to pretty print config: %v", err)
				}
				fmt.Printf("Validated config: %s\n", prettyConfig)
			}
		})
	}
}
