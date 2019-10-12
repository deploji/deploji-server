package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/deploji/deploji-server/dto"
	"github.com/deploji/deploji-server/models"
	"github.com/deploji/deploji-server/settings"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/ldap.v3"
	"log"
	"net/http"
	"strings"
	"time"
)

var ldapPort = uint16(389)
var ldapTLSPort = uint16(636)

func AuthenticateDatabase(user *models.User, password string) (bool, error) {
	if !CheckPasswordHash(password, string(user.Password)) {
		log.Printf("Pasword not match for user: %s", user.Username)
		return false, fmt.Errorf("bad password")
	}
	return true, nil
}

func AuthenticateLDAP(user *models.User, password string) (bool, error) {
	ldapServer := models.GetSettingValue("LDAP", "host", "localhost")
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapServer, ldapPort))
	if err != nil {
		log.Printf("Dial error: %s", err)
		return false, err
	}
	log.Printf("connection successfull: %v", l)
	domain := models.GetSettingValue("LDAP", "domain", "localhost")
	br, err := l.SimpleBind(ldap.NewSimpleBindRequest(fmt.Sprintf("%s\\%s", domain, user.Username), password, []ldap.Control{}))
	if err != nil {
		log.Printf("Bind error: %s", err)
		return false, err
	}
	log.Printf("bind successfull: %v", br)
	return true, nil
}

func HashPassword(password models.Password) (models.Password, error) {
	if password == "" {
		return "", nil
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return models.Password(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(user *models.User) (*dto.JWT, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Username,
		"uid": user.ID,
		"utp": user.Type,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(settings.Auth.TTL).Unix(),
		"nbf": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(settings.Auth.JWTSecret))
	if err != nil {
		return nil, err
	}
	return &dto.JWT{Token: tokenString}, nil
}

func RefreshToken(r *http.Request) (*dto.JWT, error) {
	oldToken, err := ParseToken(TokenGetter(r))
	if err != nil && err.Error() != "Token is expired" {
		log.Printf("Invalid token: %s", err)
		return nil, fmt.Errorf("invalid token")
	}

	if err := VerifyClaims(oldToken, true, true, false); err != nil {
		log.Printf("VerifyClaims error: %s", err)
		return nil, err
	}

	user := models.GetUserByUsername(oldToken.Claims.(jwt.MapClaims)["sub"].(string))
	if user == nil || user.IsActive == false {
		return nil, errors.New("user not found or inactive")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Username,
		"uid": user.ID,
		"iat": oldToken.Claims.(jwt.MapClaims)["iat"],
		"utp": user.Type,
		"exp": time.Now().Add(settings.Auth.TTL).Unix(),
		"nbf": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(settings.Auth.JWTSecret))
	if err != nil {
		return nil, err
	}
	return &dto.JWT{Token: tokenString}, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method: %v", token.Header["alg"])
			return false, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(settings.Auth.JWTSecret), nil
	})
}

func VerifyClaims(token *jwt.Token, verifyIAT bool, verifyNBF bool, verifyEXP bool) error {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("jwt error")
	}
	if verifyIAT && !claims.VerifyIssuedAt(time.Now().Unix(), true) {
		return fmt.Errorf("token used before issued")
	}
	if verifyNBF && !claims.VerifyNotBefore(time.Now().Unix(), true) {
		return fmt.Errorf("token from future")
	}
	if verifyEXP && !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return fmt.Errorf("token expired")
	}
	if !VerifyRefreshTTL(claims["iat"], time.Now().Add(-settings.Auth.RefreshTTL).Unix(), true) {
		return fmt.Errorf("refresh ttl expired")
	}
	return nil
}

func VerifyRefreshTTL(claim interface{}, cmp int64, req bool) bool {
	switch claim := claim.(type) {
	case float64:
		return verifyExp(int64(claim), cmp, req)
	case json.Number:
		v, _ := claim.Int64()
		return verifyExp(v, cmp, req)
	}
	return req == false
}

func verifyExp(exp int64, now int64, required bool) bool {
	if exp == 0 {
		return !required
	}
	return now <= exp
}

func TokenGetter(r *http.Request) string {
	return strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
}

func GetJWTClaims(r *http.Request) dto.JWTClaims {
	token, _ := ParseToken(TokenGetter(r))
	claims := token.Claims.(jwt.MapClaims)
	var jwtClaims dto.JWTClaims
	jwtClaims.Sub = claims["sub"].(string)
	jwtClaims.UserID = uint(claims["uid"].(float64))
	jwtClaims.Type = models.UserType(claims["utp"].(string))
	return jwtClaims
}
