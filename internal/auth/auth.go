package auth

import (
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
