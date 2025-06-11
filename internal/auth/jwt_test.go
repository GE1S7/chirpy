package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestInvalidSecret(t *testing.T) {
	userIDIn := uuid.New()
	tokenSecret := "examplesecret"
	expiresIn, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf("%v", err)
	}

	jwt, err := MakeJWT(userIDIn, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("%v", err)
	}

	invalidSecret := "falsehood"

	_, err = ValidateJWT(jwt, invalidSecret)
	if err == nil {
		t.Errorf("Invalid secret accepted")
	}
}

func TestExpiredSecret(t *testing.T) {
	userIDIn := uuid.New()
	tokenSecret := "examplesecret"
	expiresIn, err := time.ParseDuration("1s")
	if err != nil {
		t.Errorf("%v", err)
	}

	jwt, err := MakeJWT(userIDIn, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("%v", err)
	}

	wait, err := time.ParseDuration("2s")
	if err != nil {
		t.Errorf("%v", err)
	}

	time.Sleep(wait)

	_, err = ValidateJWT(jwt, tokenSecret)
	if err == nil {
		t.Errorf("Expired JWT validated")
	}

}

func TestJWTFuncs(t *testing.T) {
	userIDIn := uuid.New()
	tokenSecret := "examplesecret"
	expiresIn, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf("%v", err)
	}

	jwt, err := MakeJWT(userIDIn, tokenSecret, expiresIn)
	if err != nil {
		t.Errorf("%v", err)
	}
	fmt.Println(jwt)

	userIDOut, err := ValidateJWT(jwt, tokenSecret)
	if err != nil {
		t.Errorf("%v", err)
	}

	if userIDIn != userIDOut {
		t.Errorf("UserID mismatch: %v (in) != %v (out)", userIDIn, userIDOut)
	}

}
