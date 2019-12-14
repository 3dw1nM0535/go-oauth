package handlers

import (
	"crypto/rand"
	"encoding/base64"
)

// RandToken generates a random @l length token
func RandToken(l int) string {
	b := make([]byte, l)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
