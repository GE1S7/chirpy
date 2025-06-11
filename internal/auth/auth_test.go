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

func TestTokenSame(t *testing.T) {
	expected := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiI0NTIyYjc2Mi1kNzAxLTQ1M2MtYTYyMy1jMzE3NDUxMWJhM2QiLCJleHAiOjE3NDk2NTEzNzUsImlhdCI6MTc0OTY0Nzc3NX0.waapqvDhwZiV58xr0M0eeKFhTgFsZawNtpD9v7zhTJk"
	req, err := http.NewRequest("GET", "https://api.example.com/data", nil)
	if err != nil {
		t.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiI0NTIyYjc2Mi1kNzAxLTQ1M2MtYTYyMy1jMzE3NDUxMWJhM2QiLCJleHAiOjE3NDk2NTEzNzUsImlhdCI6MTc0OTY0Nzc3NX0.waapqvDhwZiV58xr0M0eeKFhTgFsZawNtpD9v7zhTJk")

	actual, err := GetBearerToken(req.Header)
	if err != nil {
		t.Errorf("%v", err)
	}

	if expected != actual {
		t.Errorf("expected token (%v) does not match retrieved (%v)", expected, actual)
	}
}
