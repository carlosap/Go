package environment

import (
	"log"
	"os"
)

// GetSetEnv Gets and potentially sets environment variable to the fallback value.
// Returns environment-set value if present, fallback otherwise
func GetSet(envVar, fallback string) string {
	envString := os.Getenv(envVar)
	if envString == "" {
		envString = fallback
		if err := os.Setenv(envVar, envString); err != nil {
			log.Printf("Unable to set env %s = %s: %v", envVar, envString, err)
		}
	}

	return envString
}