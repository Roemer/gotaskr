package goext

import "os"

// EnvExists checks if the given environment variable exists or not.
func EnvExists(key string) bool {
	if _, ok := os.LookupEnv(key); ok {
		return true
	}
	return false
}

// GetEnvOrDefault returns the value if the environment variable exists or the default otherwise.
func GetEnvOrDefault(key string, defaultValue string) (string, bool) {
	if _, ok := os.LookupEnv(key); ok {
		return key, true
	}
	return defaultValue, false
}
