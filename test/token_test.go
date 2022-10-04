package test

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserId int
	jwt.StandardClaims
}

func TestToken(t *testing.T) {
	claims := &Claims{}

	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsImV4cCI6MTY2NDg1NjcxNCwiaWF0IjoxNjY0ODU2NzE0LCJpc3MiOiJmcmVlZGIuY29tIiwic3ViIjoidXNlciB0b2tlbiJ9.5hf4Vmsf6IMlsV4qc0WZT6e6wy06_O0RGe9HS_aAKss"

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte("abcdefg"), nil
	})

	if err != nil {
		t.Errorf("err: %v\n", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		t.Logf("claims: %v\n", claims)
	}

	t.Error("invalid token")
}
