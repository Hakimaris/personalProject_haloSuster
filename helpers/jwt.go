package helpers

import (
	"github.com/golang-jwt/jwt/v5"
)

type DB_USERNAMEClaims struct {
	jwt.RegisteredClaims
	Id  string `json:"id"`
	Nip string `json:"nip"`
}

