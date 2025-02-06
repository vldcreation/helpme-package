package configurator

import (
	"os"
	"strconv"
)

// GetEnv retrieves the value of the environment variable named by the key.
// It returns the value, which will be defaultVal if the variable is not present.
func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// GetEnvAsInt retrieves the value of the environment variable named by the key and converts it to an integer.
// It returns the integer value, which will be defaultVal if the variable is not present or cannot be converted to an integer.
func GetEnvAsInt(key string, defaultVal int) int {
	valueStr := GetEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// GetEnvAsBool retrieves the value of the environment variable named by the key and converts it to a boolean.
// It returns the boolean value, which will be defaultVal if the variable is not present or cannot be converted to a boolean.
func GetEnvAsBool(key string, defaultVal bool) bool {
	valueStr := GetEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// GetEnvAsFloat retrieves the value of the environment variable named by the key and converts it to a float64.
// It returns the float64 value, which will be defaultVal if the variable is not present or cannot be converted to a float64.
func GetEnvAsFloat(key string, defaultVal float64) float64 {
	valueStr := GetEnv(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return defaultVal
}
