package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("main/templates/mainpage.html")
	if err != nil {
		fmt.Println("parsing main template: ", err)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		fmt.Println("executing main template: ", err)
		return
	}
}

func getAllAthletes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	db := dbconn()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM athlete")
	defer rows.Close()
	if err != nil {
		fmt.Println("get query: ", err)
	}

	var b []Athlete
	for rows.Next() {
		a := Athlete{}
		err := rows.Scan(&a.Id, &a.Birth, &a.Country, &a.Name, &a.Surname, &a.Weight)
		if err != nil {
			fmt.Println("retrieving athlete: ", err)
		}
		b = append(b, a)
	}

	json.NewEncoder(w).Encode(b)
}

func getAthleteByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println("param is not int: ", err)
	}

	db := dbconn()
	defer db.Close()

	row := db.QueryRow("SELECT * FROM athlete WHERE id = $1", id)
	a := Athlete{}
	err = row.Scan(&a.Id, &a.Birth, &a.Country, &a.Name, &a.Surname, &a.Weight)
	if err != nil {
		fmt.Println("retrieving athlete: ", err)
	}
	json.NewEncoder(w).Encode(a)
}

func newAthlete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("main/templates/add.html")
		if err != nil {
			fmt.Println("parsing add template: ", err)
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			fmt.Println("executing add template: ", err)
			return
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("parsing add form: ", err)
			return
		}
		a := Athlete{}
		a.Name, a.Surname, a.Country = strings.Join(r.Form["name"], ""),
			strings.Join(r.Form["surname"], ""), strings.Join(r.Form["country"], "")
		a.Birth, err = time.Parse("2006-01-02", strings.Join(r.Form["birth"], ""))
		if err != nil {
			fmt.Println("birth ", err)
		}
		a.Weight, err = strconv.ParseFloat(strings.Join(r.Form["weight"], ""), 64)
		if err != nil {
			fmt.Println("weight ", err)
		}

		db := dbconn()
		defer db.Close()
	}
}
