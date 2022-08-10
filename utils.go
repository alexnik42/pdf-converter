package main

import (
	"crypto/rand"
	"encoding/hex"
)

func generateUniqueToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "default"
	}
	return hex.EncodeToString(b)
}
