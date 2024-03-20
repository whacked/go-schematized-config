package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// GetEnvironAsMap returns the environment variables as a map[string]interface{}.
// Note: Environment variables are inherently strings, so the values are all strings,
// but they are put in an interface{} map for flexibility.
func GetEnvironAsMap() map[string]interface{} {
	envMap := make(map[string]interface{})
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 {
			envMap[pair[0]] = pair[1]
		}
	}
	return envMap
}

func main() {
	args := os.Args[1:] // Ignore the program name

	if len(args) < 1 || len(args) > 2 {
		fmt.Println("Usage: myapp <path to schema> [path to .env file]")
		fmt.Println(" - If no .env path is provided, the program will attempt to load a .env file from the current directory.")
		return
	}

	schemaPath := args[0]
	var err error

	if len(args) == 2 {
		// If a .env path is provided, load it
		err = godotenv.Load(args[1])
	} else {
		// Attempt to load .env from the current directory
		err = godotenv.Load()
	}

	if err != nil {
		log.Printf("[WARN] Could not load .env: %v", err)
	} else {
		log.Printf("[WARN] Loaded .env")
	}

	envMap := GetEnvironAsMap()
	jsonSchema, err := LoadJsonnetFile(schemaPath) // Use the schemaPath variable
	if err != nil {
		log.Fatalf("Failed to load jsonnet: %v", err)
	}
	validatedConfig, err := LoadValidatedConfig(jsonSchema, envMap)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// pretty print the validated config
	validatedConfigBytes, err := json.MarshalIndent(validatedConfig, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal validated config: %v", err)
	}
	fmt.Println(string(validatedConfigBytes))
}
