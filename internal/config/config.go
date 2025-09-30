package config

import (
	"os"
)

var JWTSecret = getJWTSecret()

func getJWTSecret() []byte {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return []byte(secret)
	}
	return []byte("my-very-secret-key")
}
