package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func hello(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("main/templates/mainpage.html"))
	err := t.Execute(w, nil)
	if err != nil {
		fmt.Println("executing main template: ", err)
		return
	}
}

func table(w http.ResponseWriter, r *http.Request) {
	db := dbconn()
	defer db.Close()

	row := db.QueryRow("SELECT * FROM athlete WHERE surname = $1", "Petrov")
	a := Athlete{}
	err := row.Scan(&a.Id, &a.Birth, &a.SportClub, &a.Name, &a.Surname, &a.Weight)

	row = db.QueryRow("SELECT * FROM athlete WHERE surname = $1", "Juraev")
	b := Athlete{}
	err = row.Scan(&b.Id, &b.Birth, &b.SportClub, &b.Name, &b.Surname, &b.Weight)

	type AthleteTable struct {
		Table []Athlete
	}
	tab := AthleteTable{}
	tab.Table = append(tab.Table, a, b)

	fmt.Println(a.toString())
	t := template.Must(template.ParseFiles("main/templates/athletetable.html"))
	err = t.Execute(w, tab)
	if err != nil {
		fmt.Println("executing athletetable template: ", err)
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
		err := rows.Scan(&a.Id, &a.Birth, &a.SportClub, &a.Name, &a.Surname, &a.Weight)
		if err != nil {
			fmt.Println("retrieving athlete: ", err)
		}
		b = append(b, a)
	}

	json.NewEncoder(w).Encode(b)
}

func surnameFinding(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("main/templates/getsurname.html")
		if err != nil {
			fmt.Println("parsing getsurname template: ", err)
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			fmt.Println("executing getsurname template: ", err)
			return
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			fmt.Println("parsing getsurname form: ", err)
			return
		}
		surname := strings.Join(r.Form["surname"], "")
		if surname == "" {
			t, err := template.ParseFiles("main/templates/getfail.html")
			if err != nil {
				fmt.Println("parsing getfail template: ", err)
				return
			}
			err = t.Execute(w, nil)
			if err != nil {
				fmt.Println("executing getfail template: ", err)
				return
			}
			return
		}

		db := dbconn()
		defer db.Close()

		row := db.QueryRow("SELECT * FROM athlete WHERE surname = $1", surname)
		a := Athlete{}
		err = row.Scan(&a.Id, &a.Birth, &a.SportClub, &a.Name, &a.Surname, &a.Weight)
		if err != nil {
			fmt.Println("retrieving athlete: ", err)
			t, err := template.ParseFiles("main/templates/getfail.html")
			if err != nil {
				fmt.Println("parsing getfail template: ", err)
				return
			}
			err = t.Execute(w, nil)
			if err != nil {
				fmt.Println("executing getfail template: ", err)
				return
			}
			return
		}
		json.NewEncoder(w).Encode(a)
	}
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
		a.Name, a.Surname, a.SportClub = strings.Join(r.Form["name"], ""),
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

		row := db.QueryRow("SELECT max(id) FROM athlete")
		var maxid int
		err = row.Scan(&maxid)
		if err != nil {
			fmt.Println("retrieving max id : ", err)
		}
		a.Id = maxid + 1
		stmt := "INSERT INTO athlete (id, birth, country, name, surname, weight) VALUES($1, $2, $3, $4, $5, $6)"
		_, err = db.Query(stmt, a.Id, a.Birth, a.SportClub, a.Name, a.Surname, a.Weight)
		if err != nil {
			fmt.Println(err)
			return
		}

		t, err := template.ParseFiles("main/templates/addsuccess.html")
		if err != nil {
			fmt.Println("parsing addsuccess template: ", err)
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			fmt.Println("executing addsuccess template: ", err)
			return
		}
	}
}
