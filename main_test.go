package main

import (
	"reflect"
	"testing"
)

func TestCoercePrimitiveValues(t *testing.T) {
	tests := []struct {
		name        string
		jsonSchema  map[string]interface{}
		data        map[string]string
		expected    map[string]interface{}
		expectError bool
	}{
		{
			name: "everything is correct",
			jsonSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"STRING":  map[string]interface{}{"type": "string"},
					"NUMBER":  map[string]interface{}{"type": "number"},
					"INTEGER": map[string]interface{}{"type": "number"},
					"BOOLEAN": map[string]interface{}{"type": "boolean"},
				},
			},
			data: map[string]string{
				"STRING":  "asdf",
				"NUMBER":  "1232529.56",
				"INTEGER": "98758585858232",
				"BOOLEAN": "TRUE",
			},
			expected: map[string]interface{}{
				"STRING":  "asdf",
				"NUMBER":  1232529.56,
				"INTEGER": float64(98758585858232), // Go's JSON unmarshal converts numbers to float64 by default
				"BOOLEAN": true,
			},
			expectError: false,
		},
		{
			name: "undeclared types do not get converted",
			jsonSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"STRING":         map[string]interface{}{"type": "string"},
					"SOMETHING_ELSE": map[string]interface{}{},
				},
			},
			data: map[string]string{
				"NUMBER":         "1232529.56",
				"INTEGER":        "98758585858232",
				"BOOLEAN":        "TRUE",
				"SOMETHING_ELSE": "value", // This test case needs adjustment, as the original Python version cannot be directly translated due to type differences
			},
			expected: map[string]interface{}{
				"NUMBER":         "1232529.56",
				"INTEGER":        "98758585858232",
				"BOOLEAN":        "TRUE",
				"SOMETHING_ELSE": "value",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CoercePrimitiveValues(tt.jsonSchema, tt.data)
			if (err != nil) != tt.expectError {
				t.Errorf("CoercePrimitiveValues() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("CoercePrimitiveValues() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractDeclaredItems(t *testing.T) {
	jsonSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"STRING":         map[string]string{"type": "string"},
			"SOMETHING_ELSE": map[string]interface{}{},
			"HAS_DEFAULT":    map[string]interface{}{"type": "boolean", "default": "NO COERCION!"},
		},
	}
	data := map[string]interface{}{
		"NUMBER":         "1232529.56",
		"INTEGER":        "98758585858232",
		"BOOLEAN":        "TRUE",
		"SOMETHING_ELSE": map[string]interface{}{"a": 1, "b": "C"},
	}
	expected := map[string]interface{}{
		"HAS_DEFAULT":    "NO COERCION!",
		"SOMETHING_ELSE": map[string]interface{}{"a": 1, "b": "C"},
	}

	result, err := ExtractDeclaredItems(jsonSchema, data)
	if err != nil {
		t.Fatalf("ExtractDeclaredItems() error: %v", err)
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ExtractDeclaredItems() got = %v, want %v", result, expected)
	}
}
