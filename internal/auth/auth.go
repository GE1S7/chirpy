package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader, ok := headers["Authorization"]
	if !ok {
		return "", fmt.Errorf("No Authorization header")
	}

	for _, e := range authHeader {
		if strings.Contains(e, "Bearer ") {
			tokenString := strings.Split(e, " ")[1]

			return tokenString, nil
		}
	}

	return "", fmt.Errorf("http Header does not contain authentication token")

}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil

}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader, ok := headers["Authorization"]
	if !ok {
		return "", fmt.Errorf("No Authorization header")
	}

	for _, e := range authHeader {
		if strings.Contains(e, "ApiKey ") {
			tokenString := strings.Split(e, " ")[1]

			return tokenString, nil
		}
	}

	return "", fmt.Errorf("http Header does not contain API key")
}
