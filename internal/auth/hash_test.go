package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashing(t *testing.T) {
	passwords := []string{"AsDfAsD", "%9Je136", "pdfjnme$4", "@kli59"}

	for _, e := range passwords {
		hashed, err := HashPassword(e)
		if err != nil {
			t.Errorf("Hashing error: %v", err)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(e)); err != nil {
			t.Errorf("Hashing error: %v", err)
		}
	}

}

func TestInvalidPairs(t *testing.T) {
	wrongPassHash := map[string]string{
		"test": "jj9",
		"dsae": "8989898989dsf",
		"888":  "aaAAasa8888888",
	}

	for key, value := range wrongPassHash {
		err := CheckPasswordHash(value, key)
		if err == nil {
			t.Errorf("False positive hash check. password: %v, invalid hash: %v", key, value)
		}
	}
}

func TestHashingFuncs(t *testing.T) {
	var passHash = make(map[string]string)
	passwords := []string{"test", "asdfasd", "%6je134", "aPldfjnme$3", "!?@klQiQ56.-"}

	for _, e := range passwords {
		hashed, err := HashPassword(e)
		passHash[e] = hashed
		if err != nil {
			t.Errorf("Hashing error: %v", err)
		}
	}

	for key, value := range passHash {
		err := CheckPasswordHash(value, key)
		if err != nil {
			t.Errorf("Checking error: %v", err)
		}

	}
}
