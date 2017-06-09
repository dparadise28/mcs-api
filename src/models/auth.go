package models

import (
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
	"log"
)

// perm keys for read write record privelages
const (
	OrdersReadWritePerm = "o_rw"
	OrdersReadPerm      = "o_r"
	StoreReadWrite      = "s_rw"
	StoreRead           = "s_r"
	ProductReadWrite    = "p_rw"
	ProductRead         = "p_r"

	// owner = owner of store; admin = access to everything
	OWNER = "owner"
	ADMIN = "admin"
)

type UserRoles struct {
	ID     bson.ObjectId `bson:"_id,omitempty" json:"user_id"`
	Access map[bson.ObjectId]string
}

type CustomClaims struct {
	Roles *UserRoles
	jwt.StandardClaims
}

type LoginResponse struct {
	Token string
}

func (ur *UserRoles) GenetateLoginToken() (string, error) {
	// Sign and get the complete encoded token as a string
	claims := CustomClaims{
		ur,
		jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    "test",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("temp_signiture_key"))
	log.Println(err)
	return tokenString, err
}
