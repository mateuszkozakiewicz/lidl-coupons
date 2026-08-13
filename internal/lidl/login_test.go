package lidl

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func makeToken(t *testing.T, claims jwt.MapClaims) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("test-secret"))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}
	return signed
}

func TestCheckTokenValid(t *testing.T) {
	t.Run("empty token is invalid", func(t *testing.T) {
		if checkTokenValid("") {
			t.Error("expected empty token to be invalid")
		}
	})

	t.Run("malformed token is invalid", func(t *testing.T) {
		if checkTokenValid("not-a-jwt") {
			t.Error("expected malformed token to be invalid")
		}
	})

	t.Run("token without exp claim is invalid", func(t *testing.T) {
		token := makeToken(t, jwt.MapClaims{"sub": "user"})
		if checkTokenValid(token) {
			t.Error("expected token without exp claim to be invalid")
		}
	})

	t.Run("expired token is invalid", func(t *testing.T) {
		token := makeToken(t, jwt.MapClaims{
			"exp": time.Now().Add(-time.Hour).Unix(),
		})
		if checkTokenValid(token) {
			t.Error("expected expired token to be invalid")
		}
	})

	t.Run("unexpired token is valid", func(t *testing.T) {
		token := makeToken(t, jwt.MapClaims{
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		if !checkTokenValid(token) {
			t.Error("expected unexpired token to be valid")
		}
	})
}
