package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", hello).Methods("GET")
	r.HandleFunc("/athlete", getAllAthletes).Methods("GET")
	r.HandleFunc("/athlete/{id}", getAthleteByID).Methods("GET")
	r.HandleFunc("/add", newAthlete).Methods("GET", "POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}
