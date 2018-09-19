package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// location of the files used for signing and verification
const (
	privateKeyPath      = "keys/token-signing.rsa"     // openssl genrsa -out token-signing.rsa 2048
	publicKeyPath       = "keys/token-signing.rsa.pub" // openssl rsa -in token-signing.rsa -pubout > token-signing.rsa.pub
	tokenValidityPeriod = 72                           // in hours
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

func loadVerifyKey() error {
	verifyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return err
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
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
			log.Println("retry loading token signing key")
			if err := loadSignKey(); err == nil {
				continue
			}
			counter++
			time.Sleep(retrySleepTime)
		}
		return "", err
	}

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(signKey)
	return tokenString, err
}

// validateToken parses and validates the token string
func validateToken(tokenString string) (claims *userClaims, valid bool, err error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple locals for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	if err := loadVerifyKey(); err != nil {
		counter := 0
		for counter < 5 {
			log.Println("retry loading token verify key")
			if err := loadVerifyKey(); err == nil {
				continue
			}
			counter++
			time.Sleep(retrySleepTime)
		}
		return nil, false, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return verifyKey, nil
	})
	if err != nil {
		return
	}

	claims, valid = token.Claims.(*userClaims)
	valid = valid && token.Valid
	return
}
