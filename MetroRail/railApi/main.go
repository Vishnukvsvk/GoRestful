package main

import (
	"database/sql"
	"log"

	"github.com/Vishnukvsvk/Metrorail/dbutils"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	//Connect to database
	db, err := sql.Open("sqlite3", "./railapi.db")
	if err != nil {
		log.Println("driver creation failed")
	}

	//Create tables
	dbutils.Initialize(db)
}
