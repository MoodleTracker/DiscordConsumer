package main

import (
	"log"
	"os"
)

// GetEnv returns the value of the environment variable or, if unset, the value of the file specified by the environment variable with the same name and "_FILE" appended.
func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		file := os.Getenv(key + "_FILE")
		if file == "" {
			log.Fatalf("Environment variable %s is empty!", key)
		}
		content, err := os.ReadFile(file)

		if err != nil {
			log.Fatal("Error trying to load environment variable '", key, "' from file '", file, "': ", err)
		}

		newlineCount := 0
		for _, char := range content {
			if char == '\n' {
				newlineCount++
				if newlineCount > 1 {
					log.Fatalf("Environment variable %s contains multiple lines!", key)
				}
			}
		}

		value = string(content)
		if value == "" {
			log.Fatalf("Environment variable %s is empty!", key)
		}
		value = removeTrailingNewlines(value)
	}
	return value
}

func removeTrailingNewlines(value string) string {
	if value[len(value)-1] == '\n' || value[len(value)-1] == '\r' {
		value = value[:len(value)-1]
	}
	if value[len(value)-1] == '\n' || value[len(value)-1] == '\r' {
		value = value[:len(value)-1]
	}
	return value
}

func GetEnvOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
