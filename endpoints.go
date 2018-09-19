package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type userCredentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type tokenMessage struct {
	Token string `json:"token"`
}

func sendJSONResponse(w http.ResponseWriter, payload interface{}) {
	enc := json.NewEncoder(w)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")

	if err := enc.Encode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}

func loginParticipantHandl(w http.ResponseWriter, req *http.Request) {
	// Only Post method is allowed:
	if req.Method != "POST" {
		err := errors.New("wrong method")
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	// Get user credentials fro request's body:
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	var creds userCredentials

	err = json.Unmarshal(body, &creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if creds == (userCredentials{}) {
		http.Error(w, errors.New("wrong json format").Error(), http.StatusBadRequest)
		return
	}

	// TODO: check credentials
	log.Println(creds)
	userID := uint(1)

	// generate token
	token, err := generateNewToken(userID, "participant")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send response
	tokenResp := tokenMessage{
		Token: token,
	}
	sendJSONResponse(w, tokenResp)
}

func loginResearcher() {

}

func loginAdmin() {

}

func signup() {

}

func validateTokenHandl(w http.ResponseWriter, req *http.Request) {
	tokenString := req.Header.Get("Authorization")
	if len(tokenString) == 0 {
		http.Error(w, errors.New("missing authorization header").Error(), http.StatusInternalServerError)
		return
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse and validate token
	parsedToken, ok, err := validateToken(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		http.Error(w, errors.New("token not valid").Error(), http.StatusInternalServerError)
		return
	}
	log.Println(parsedToken)
	sendJSONResponse(w, parsedToken)
}

func renewToken() {

}
