package utils

import "os"

func GetEnv(envVar, defaultVar string) string {
	v := os.Getenv(envVar)
	if v == "" {
		if defaultVar != "" {
			return defaultVar
		}
	}
	return v
}
