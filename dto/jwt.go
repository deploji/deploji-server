package dto

type JWT struct {
	Token string
}

type JWTClaims struct {
	UserID uint
	Sub    string
}
