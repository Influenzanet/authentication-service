package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Hello World")
	r := mux.NewRouter()

	r.HandleFunc("/login", loginParticipantHandl)
	r.HandleFunc("/validate", validateTokenHandl)
	log.Fatal(http.ListenAndServe(":3100", r))
}
