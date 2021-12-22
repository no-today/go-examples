package jwt

import (
	"github.com/golang-jwt/jwt"
	"reflect"
	"testing"
)

func TestTokenGenerateAndParse(t *testing.T) {
	claims := Claims{
		Principal:      "no-today",
		Authorities:    []string{"ADMIN", "USER"},
		StandardClaims: jwt.StandardClaims{},
	}

	secret := "AAA"
	token, err := generateToken(claims, secret, 10*60)
	assertNoError(t, err)

	token_, err := parseToken(token.Token, secret)
	assertNoError(t, err)

	c := token_.Claims.(*Claims)
	assertEquals(t, c.Principal, claims.Principal)
	assertDeepEquals(t, c.Authorities, claims.Authorities)
}

func assertDeepEquals(t *testing.T, got []string, want []string) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func assertEquals(t *testing.T, got string, want string) {
	if got != want {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("No exception expected, but appeared: %v", err)
	}
}
