package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "3572033"
	DB_NAME     = "Athletes"
)

func dbconn() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)
	//defer db.Close()
	if err != nil {
		log.Fatal("open database: ", err)
	}
	return db
}
