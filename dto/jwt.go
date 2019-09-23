package dto

import "github.com/sotomskir/mastermind-server/models"

type JWT struct {
	Token string
}

type JWTClaims struct {
	UserID uint
	Sub    string
	Type   models.UserType
}
