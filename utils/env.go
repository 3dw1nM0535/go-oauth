package utils

import (
	"log"
	"os"
)

// MustGet : returns Environment variable otherwise error
func MustGet(envVar string) string {
	v := os.Getenv(envVar)
	if v == "" {
		log.Panicf("Env variable missing: " + v)
	}
	return v
}
