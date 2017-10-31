package tools

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	mrand "math/rand"
	"time"
)

const (
	ConfirmationCodeRandomCharKeySpace = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-!$&:;=?@_~"
	DeafultPWHashCost                  = 8
	letterIdxBits                      = 6                    // 6 bits to represent 64 possibilities / indexes
	letterIdxMask                      = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

var ConfirmationCodeLengthRange = []int{24, 40}

func RandomIntInRange(minMax ...int) int {
	if len(minMax) != 2 {
		panic("A valid range must be an int slice of len==2")
	}
	mrand.Seed(time.Now().UTC().UnixNano())
	return mrand.Intn(minMax[1]-minMax[0]) + minMax[0]
}

func SecureRandomBytes(length int) []byte {
	var randomBytes = make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("Unable to generate random bytes")
	}
	return randomBytes
}

func SecureRandomString(availableCharBytes string, length int) string {
	// Compute bitMask
	availableCharLength := len(availableCharBytes)
	if availableCharLength == 0 || availableCharLength > 256 {
		panic("availableCharBytes length must be greater than 0 and less than or equal to 256")
	}
	var bitLength byte
	var bitMask byte
	for bits := availableCharLength - 1; bits != 0; {
		bits = bits >> 1
		bitLength++
	}
	bitMask = 1<<bitLength - 1

	// Compute bufferSize
	bufferSize := length + length/3

	// Create random string
	result := make([]byte, length)
	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			// Random byte buffer is empty, get a new one
			randomBytes = SecureRandomBytes(bufferSize)
		}
		// Mask bytes to get an index into the character slice
		if idx := int(randomBytes[j%length] & bitMask); idx < availableCharLength {
			result[i] = availableCharBytes[idx]
			i++
		}
	}

	return string(result)
}

func GenerateConfirmationCode() string {
	return SecureRandomString(
		ConfirmationCodeRandomCharKeySpace,
		RandomIntInRange(ConfirmationCodeLengthRange...),
	)
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
