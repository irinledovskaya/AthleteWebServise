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
	r.HandleFunc("/finding", surnameFinding).Methods("GET", "POST")
	r.HandleFunc("/add", newAthlete).Methods("GET", "POST")
	r.HandleFunc("/update", updateAthlete).Methods("GET", "POST")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./main/templates/")))
	log.Fatal(http.ListenAndServe(":8080", r))
}
