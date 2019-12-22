package utils

import (
	"log"
	"os"
	"strconv"
)

// MustGet : returns Environment variable otherwise error
func MustGet(envVar string) string {
	v := os.Getenv(envVar)
	if v == "" {
		log.Panicf("Env variable missing: " + envVar)
	}
	return v
}

// MustGetBool : returns bool Environment variable otherwise error
func MustGetBool(envVar string) bool {
	v := os.Getenv(envVar)
	if v == "" {
		log.Panicf("Env variable missing: " + envVar)
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		log.Panicf("Error parsing env variable: " + err.Error())
	}
	return b
}
