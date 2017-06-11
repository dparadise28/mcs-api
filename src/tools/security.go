package tools

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	//"log"
)

var (
	DeafultPWHashCost      = 14
	ConfirmationCodeLength = 32
)

func GenerateConfirmationCode() string {
	n := ConfirmationCodeLength
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), DeafultPWHashCost)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
