package model

import "github.com/golang-jwt/jwt"

type User struct {
	Uuid  string `json:"uuid"`
	Token string `json:"token"`
}

type Response struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type JWTClaim struct {
	Uuid           string `json:"uuid"`
	StandardClaims jwt.StandardClaims
}

func (J *JWTClaim) Valid() error {
	//TODO implement me
	panic("implement me")
}
