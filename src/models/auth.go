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

type UserRolesResponse struct {
	Access map[string]string
	Stores map[string]string `bson:"store_map" json:"store_map"`
}

type CustomClaims struct {
	Perms            map[string]string
	Confirmed        bool
	StripeCustomerID string
	jwt.StandardClaims
}

func GenerateTokenClaims(access map[string]string, confirmed bool, stripeID string) CustomClaims {
	return CustomClaims{
		Perms:            access,
		Confirmed:        confirmed,
		StripeCustomerID: stripeID,
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

func (u *User) GenetateLoginToken(pw string) (string, error) {
	if !u.Confirmed {
		return "", errors.New(UNCONFIRMED_USER)
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
	if pw == "" || err != nil {
		return "", errors.New("Unauthorized access")
	}

	token, err := u.UpdateToken()
	return token, err
}

func GetJWTContent(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(JWT_SIGNATURE), nil
		},
	)
	return token, err
}

func (u *User) GenetateLoginTokenAndSetHeaders(pw string, w http.ResponseWriter) (string, error) {
	if !u.Confirmed {
		return "", errors.New(UNCONFIRMED_USER)
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pw))
	if pw == "" || err != nil {
		return "", errors.New("Unauthorized access")
	}

	token, err := u.UpdateTokenAndCookie(w)
	return token, err
}

func (u *User) UpdateToken() (string, error) {
	claims := GenerateTokenClaims(u.Roles.Access, u.Confirmed, u.StripeCustomerID)
	token, err := claims.CreateToken()

	u.Login.Token = token
	u.Login.UID = u.ID
	return token, err
}

func (ur *User) UpdateTokenAndCookie(w http.ResponseWriter) (string, error) {
	token, err := ur.UpdateToken()
	w.Header().Set(JWT_COOKIE_NAME, token)
	w.Header().Set(USERID_COOKIE_NAME, ur.ID.Hex())
	return token, err
}
