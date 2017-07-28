package main

import (
	"log"
	"os"
	"strings"
)

func setConfigValue(parameter string, defaultValue string) string {
	value := os.Getenv(parameter)
	if value == "" {
		log.Println(strings.Join([]string{"$", parameter, " not set. Using program defaults."}, ""))
		return defaultValue
	}
	return value
}

func loadEnvironmentVars() map[string]string {
	result := make(map[string]string)
	result["web_host"] = setConfigValue("WOLF_WEBHOST", "0.0.0.0")
	result["web_port"] = setConfigValue("WOLF_WEBPORT", "8080")
	result["msx_port"] = setConfigValue("WOLF_MSXPORT", "9300")
	return result
}
