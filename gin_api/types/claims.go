package types

import (
    "github.com/golang-jwt/jwt"
)

type Claims struct {
    Id    Int32      `json:"id"`
    Name  string   `json:"name"`
    Roles []string `json:"roles"`
    jwt.StandardClaims
}