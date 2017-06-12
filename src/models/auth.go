package models

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// perm keys for read write record privelages
const (
	ACCESSROLE_OrdersReadWritePerm  = "o_rw"
	ACCESSROLE_OrdersReadPerm       = "o_r"
	ACCESSROLE_StoreReadWritePerm   = "s_rw"
	ACCESSROLE_StoreReadPerm        = "s_r"
	ACCESSROLE_ProductReadWritePerm = "p_rw"
	ACCESSROLE_ProductReadPerm      = "p_r"

	// owner = owner of store; admin = access to everything
	ACCESSROLE_STOREOWNER     = "owner"
	ACCESSROLE_ADMIN          = "admin"
	ACCESSROLE_CONFIRMED_USER = "confirmed"

	UNCONFIRMED_USER = "Unconfirmed User"

	JWT_SIGNATURE = "temp_signiture_key"
	JWT_ISSUER    = "MCS-API"

	JWT_COOKIE_NAME    = "authtoken"
	USERID_COOKIE_NAME = "userID"
	USERID_HEADER_NAME = USERID_COOKIE_NAME
)

func JWT_TTL() int64 {
	return time.Now().Add(time.Minute * 60 * 24 * 15).Unix()
}

func COOKIE_TTL() time.Time {
	return time.Now().Add(3 * 24 * time.Hour)
}

type UserRoles struct {
	Access map[string]string
}

type CustomClaims struct {
	Perms     map[string]string
	Confirmed bool
	jwt.StandardClaims
}

func GenerateTokenClaims(access map[string]string, confirmed bool) CustomClaims {
	return CustomClaims{
		Perms:     access,
		Confirmed: confirmed,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: JWT_TTL(),
			Issuer:    JWT_ISSUER,
		},
	}
}

func (c *CustomClaims) CreateToken() (string, error) {
	// Sign and get the complete encoded token as a string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	tokenString, tokenErr := token.SignedString([]byte(JWT_SIGNATURE))
	if tokenErr != nil {
		return tokenString, tokenErr
	}
	return tokenString, tokenErr
}

func (ur *User) GenetateLoginToken(pw string) (string, error) {
	if !ur.Confirmed {
		return "", errors.New(UNCONFIRMED_USER)
	}
	err := bcrypt.CompareHashAndPassword([]byte(ur.Password), []byte(pw))
	if pw == "" || err != nil {
		return "", errors.New("Unauthorized access")
	}

	// Sign and get the complete encoded token as a string
	claims := GenerateTokenClaims(ur.Roles.Access, ur.Confirmed)
	token, err := claims.CreateToken()
	return token, err
}

func (ur *User) UpdateTokenAndCookie(w http.ResponseWriter) (string, error) {
	// Sign and get the complete encoded token as a string
	claims := GenerateTokenClaims(ur.Roles.Access, ur.Confirmed)
	token, err := claims.CreateToken()
	authCookie := http.Cookie{
		Name:    JWT_COOKIE_NAME,
		Value:   token,
		Expires: COOKIE_TTL(),
	}
	http.SetCookie(w, &authCookie)
	return token, err
}

//func (ur *User) GenetateLoginToken(pw string) (string, error) {
