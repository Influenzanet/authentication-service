package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	b64 "encoding/base64"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	secretKey    []byte
	secretKeyEnc string
)

type userClaims struct {
	ID         string            `json:"id,omitempty"`
	InstanceID string            `json:"instance_id,omitempty"`
	Payload    map[string]string `json:"payload,omitempty"`
	jwt.StandardClaims
}

func checkTokenAgeMaturity(issuedAt int64) bool {
	return time.Now().Unix() < time.Unix(issuedAt, 0).Add(conf.JWT.TokenMinimumAgeMin).Unix()
}

func getSecretKey() (newSecretKey []byte, err error) {
	newSecretKeyEnc := os.Getenv("JWT_TOKEN_KEY")
	if secretKeyEnc == newSecretKeyEnc {
		return newSecretKey, nil
	}
	secretKeyEnc = newSecretKeyEnc
	newSecretKey, err = b64.StdEncoding.DecodeString(newSecretKeyEnc)
	if err != nil {
		return newSecretKey, err
	}
	if len(newSecretKey) < 32 {
		return newSecretKey, errors.New("couldn't find proper secret key")
	}
	secretKey = newSecretKey
	return
}

// GenerateNewToken create and signes a new token
func generateNewToken(userID string, userRoles []string, instanceID string) (string, error) {
	payload := map[string]string{}

	if len(userRoles) > 0 {
		payload["roles"] = strings.Join(userRoles, ",")
	}

	// Create the Claims
	claims := userClaims{
		userID,
		instanceID,
		payload,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(conf.JWT.TokenExpiryInterval).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey, err := getSecretKey()
	if err != nil {
		return "", err
	}

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secretKey)
	return tokenString, err
}

// validateToken parses and validates the token string
func validateToken(tokenString string) (claims *userClaims, valid bool, err error) {
	secretKey, err := getSecretKey()
	if err != nil {
		return nil, false, err
	}

	token, err := jwt.ParseWithClaims(tokenString, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if token == nil {
		return
	}
	claims, valid = token.Claims.(*userClaims)
	valid = valid && token.Valid
	return
}
