package jwt_helper

import (
	"errors"
	"strings"
	"time"

	jwts "github.com/dgrijalva/jwt-go"
)

func ClaimTokenData(token, publicKey string) (*IDTokenRequest, error) {
	key, err := jwts.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return nil, errors.New("the token is invalid")
	}

	parts := strings.Split(token, ".")
	if len(parts) < 3 {
		return nil, errors.New("the token is invalid")
	}
	err = jwts.SigningMethodRS256.Verify(strings.Join(parts[0:2], "."), parts[2], key)
	if err != nil {
		return nil, errors.New("the token is invalid")
	}

	claims := jwts.MapClaims{}
	jwts.ParseWithClaims(token, claims, func(token *jwts.Token) (interface{}, error) {
		return nil, nil
	})

	return toIDTokenRequest(claims), nil
}

func Verify(token, publicKey string) (*IDTokenRequest, error) {
	data, err := ClaimTokenData(token, publicKey)
	if err != nil {
		return nil, err
	}

	if float64(time.Now().Unix()) >= data.Exp {
		return nil, errors.New("the token has expired")
	}

	return data, nil
}

func toIDTokenRequest(claims jwts.MapClaims) *IDTokenRequest {
	return &IDTokenRequest{
		Nbf:        claims["nbf"].(float64),
		Iss:        claims["iss"].(string),
		Sub:        claims["sub"].(string),
		Iat:        claims["iat"].(float64),
		Exp:        claims["exp"].(float64),
		Aud:        claims["aud"].(string),
		Email:      claims["email"].(string),
		Username:   claims["username"].(string),
		Nonce:      claims["nonce"].(string),
		FamilyName: claims["family_name"].(string),
		GivenName:  claims["given_name"].(string),
		Token:      claims["token"].(string),
	}
}
