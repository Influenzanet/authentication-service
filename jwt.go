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
	privateKeyPath   = "keys/token-signing.rsa"     // openssl genrsa -out token-signing.rsa 2048
	publicKeyPath    = "keys/token-signing.rsa.pub" // openssl rsa -in token-signing.rsa -pubout > token-signing.rsa.pub
	oldPublicKeyPath = "keys/old-token-signing.rsa.pub"
	retrySleepTime   = 100 * time.Millisecond
)

var (
	verifyKey           *rsa.PublicKey
	tokenValidityPeriod = time.Hour * 72   // in hours
	minTokenAge         = time.Minute * 30 // don't allow token renewal before that time
)

type userClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func loadSignKey(keyPath string, retry int) (*rsa.PrivateKey, error) {
	signBytes, err := ioutil.ReadFile(keyPath)

	// Retry when file reading failed:
	trials := 1
	for err != nil && retry > trials {
		log.Println(err.Error())
		signBytes, err = ioutil.ReadFile(keyPath)
		time.Sleep(retrySleepTime)
		trials++
	}
	if err != nil {
		return nil, err
	}

	// Parse key;
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, err
	}
	return signKey, nil
}

func loadVerifyKey(keyPath string, retry int) (*rsa.PublicKey, error) {
	verifyBytes, err := ioutil.ReadFile(keyPath)
	// Retry when file reading failed:
	trials := 1
	for err != nil && retry > trials {
		log.Println(err.Error())
		verifyBytes, err = ioutil.ReadFile(keyPath)
		time.Sleep(retrySleepTime)
		trials++
	}

	if err != nil {
		return nil, err
	}

	currentVerifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, err
	}
	return currentVerifyKey, nil
}

// GenerateNewToken create and signes a new token
func generateNewToken(userID string, userRole string) (string, error) {
	// Create the Claims
	claims := userClaims{
		userID,
		userRole,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenValidityPeriod).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	signKey, err := loadSignKey(privateKeyPath, 5)
	if err != nil {
		return "", err
	}

	// Also reload verification key directly
	verifyKey, err = loadVerifyKey(publicKeyPath, 5)
	if err != nil {
		return "", err
	}

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(signKey)
	return tokenString, err
}

func parseWithOldKey(tokenString string) (*jwt.Token, error) {
	oldVerifyKey, err := loadVerifyKey(oldPublicKeyPath, 5)
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return oldVerifyKey, nil
	})
	return token, err
}

// validateToken parses and validates the token string
func validateToken(tokenString string) (claims *userClaims, valid bool, usingOldKey bool, err error) {
	// If key not loaded yet, load it
	if verifyKey == nil {
		verifyKey, err = loadVerifyKey(publicKeyPath, 5)
		if err != nil {
			return
		}
	}

	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return verifyKey, nil
	})

	if err != nil {
		err2 := err
		token, err = parseWithOldKey(tokenString)
		if err != nil {
			verifyKey, err = loadVerifyKey(publicKeyPath, 5)
			err = err2
			return
		}
		usingOldKey = true

		claims, valid = token.Claims.(*userClaims)
		valid = valid && token.Valid
		return
	}

	claims, valid = token.Claims.(*userClaims)
	valid = valid && token.Valid
	return
}
