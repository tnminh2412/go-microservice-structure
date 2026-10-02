package jwt_helper

import (
	"encoding/json"
	"go-pkg/log"
	"go-pkg/log/base"

	jwtgo "github.com/dgrijalva/jwt-go"
)

type IDTokenRequest struct {
	Nbf        float64
	Iss        string
	Sub        string
	Iat        float64
	Exp        float64
	Aud        string
	Email      string
	Username   string
	Nonce      string
	State      string
	FamilyName string
	GivenName  string
	Token      string
	Metadata   map[string]interface{}
}

func GenerateIDToken(req *IDTokenRequest, privateKey string) (string, error) {
	metadata, err := json.Marshal(req.Metadata)
	if err != nil {
		return "", err
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodRS256, jwtgo.MapClaims{
		"nbf":         req.Nbf,
		"iss":         req.Iss,
		"sub":         req.Sub,
		"iat":         req.Iat,
		"exp":         req.Exp,
		"aud":         req.Aud,
		"username":    req.Username,
		"email":       req.Email,
		"nonce":       req.Nonce,
		"state":       req.State,
		"family_name": req.FamilyName,
		"given_name":  req.GivenName,
		"token":       req.Token,
		"metadata":    string(metadata),
	})

	rsaKey, err := GetRSAKey(privateKey)
	if err != nil {
		log.Error(func() string {
			return "jwt_helper.GetRSAKey failed"
		}, base.Error(err))
		return "", err
	}

	tokenStr, err := token.SignedString(rsaKey)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
