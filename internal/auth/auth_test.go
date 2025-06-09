package auth

import "testing"

func TestHashPassword(t *testing.T) {
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

// func TestCheckPasswordHash(t *testing.T) {}
