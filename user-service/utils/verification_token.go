package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func GenerateVerificationToken() (string, time.Time) {
	bytes := make([]byte, 16)
	_, _ = rand.Read(bytes)
	expiration := time.Now().Add(30 * time.Minute)
	return hex.EncodeToString(bytes), expiration

}
