package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	req, err := http.NewRequest("GET", "https://api.example.com/data", nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer token123")

	_, err = GetBearerToken(req.Header)
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestNoBearerToken(t *testing.T) {
	req, err := http.NewRequest("GET", "https://api.example.com/data", nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "")

	_, err = GetBearerToken(req.Header)
	if err == nil {
		t.Errorf("header without token string accepted")
	}
}
