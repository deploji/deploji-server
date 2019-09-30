package dto

import "github.com/deploji/deploji-server/models"

type JWT struct {
	Token string
}

type JWTClaims struct {
	UserID uint
	Sub    string
	Type   models.UserType
}
