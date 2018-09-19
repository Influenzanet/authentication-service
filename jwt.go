package main

import (
	"crypto/rsa"
	"io/ioutil"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// location of the files used for signing and verification
const (
	privateKeyPath      = "token.rsa"     // openssl genrsa -out token-signing.rsa 2048
	publicKeyPath       = "token.rsa.pub" // openssl rsa -in token-signing.rsa -pubout > token-signing.rsa.pub
	tokenValidityPeriod = 72              // in hours
	retrySleepTime      = 100 * time.Millisecond
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

type userClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func loadSignKey() error {
	signBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return err
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return err
	}
	return nil
}

// GenerateNewToken create and signes a new token
func generateNewToken(userID uint, userRole string) (string, error) {
	// Create the Claims
	claims := userClaims{
		userID,
		userRole,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * tokenValidityPeriod).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	if err := loadSignKey(); err != nil {
		counter := 0
		for counter < 5 {
			if err := loadSignKey(); err == nil {
				continue
			}
		}
		log.Fatal(err)
	}

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(signKey)
	return tokenString, err
}
