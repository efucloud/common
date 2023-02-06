package eauth

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt/v4"
)

func ParserToken(tokenStr string, verifyKeys []*rsa.PublicKey) (*AccountClaims, error) {
	for _, verifyKey := range verifyKeys {
		token, err := jwt.ParseWithClaims(tokenStr, &AccountClaims{}, func(token *jwt.Token) (i interface{}, e error) {
			return verifyKey, nil
		})
		claims, ok := token.Claims.(*AccountClaims)
		if ok && token.Valid {
			return claims, err
		}
	}
	return nil, errors.New("token invalid")
}
