package domain

import (
	"github.com/golang-jwt/jwt/v5"
)

type Role int

const (
    _ Role = iota
    USR
    ADM
)

type User struct {
    ID int `json:"user_id"`
    Username string `json:"username"`
    Password string `json:"password"`
    Role Role `json:"role"`
    TokenExpiresAt int64 `json:"eat"`
}

type TokenClaims struct {
    UID int
    Role Role
    jwt.RegisteredClaims
}

type TokenPair struct {
    RefreshToken string `json:"refresh_token"`
    AccessToken string `json:"access_token"`
}
