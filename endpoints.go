package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type userCredentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type tokenResponse struct {
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

func login(w http.ResponseWriter, req *http.Request) {
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

	// TODO: generate token
	token := "test-token"

	/*
		err = errors.New("Test error")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	*/

	// Send response
	tokenResp := tokenResponse{
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

func checkToken() {

}

func renewToken() {

}
