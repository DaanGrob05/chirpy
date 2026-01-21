package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	secret := "my-super-secret-key"
	userID := uuid.New()

	t.Run("Create and Validate Valid Token", func(t *testing.T) {
		duration := time.Hour
		token, err := MakeJWT(userID, secret, duration)
		if err != nil {
			t.Fatalf("failed to make JWT: %v", err)
		}

		validatedID, err := ValidateJWT(token, secret)
		if err != nil {
			t.Fatalf("failed to validate JWT: %v", err)
		}

		if validatedID != userID {
			t.Errorf("expected userID %v, got %v", userID, validatedID)
		}
	})

	t.Run("Reject Expired Token", func(t *testing.T) {
		duration := -time.Hour
		token, err := MakeJWT(userID, secret, duration)
		if err != nil {
			t.Fatalf("failed to make JWT: %v", err)
		}

		_, err = ValidateJWT(token, secret)
		if err == nil {
			t.Error("expected error for expired token, but got none")
		}
	})

	t.Run("Reject Wrong Secret", func(t *testing.T) {
		duration := time.Hour
		token, err := MakeJWT(userID, secret, duration)
		if err != nil {
			t.Fatalf("failed to make JWT: %v", err)
		}

		wrongSecret := "definitely-the-wrong-key"
		_, err = ValidateJWT(token, wrongSecret)
		if err == nil {
			t.Error("expected error for wrong secret, but got none")
		}
	})
}
