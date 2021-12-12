package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	id     int
	name   string
	author string
}

func main() {
	db, err := sql.Open("sqlite3", "./books.db")
	if err != nil {
		log.Println(err)
	}

	//Create table
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS books(id INTEGER PRIMARY KEY, isbn INTEGER, author VARCHAR(64),name VARCHAR(64) NULL)")
	if err != nil {
		log.Println("Error in creating table")
	} else {
		log.Println("Successfully created table")
	}

	statement.Exec()
	dbOperations(db)
}

func dbOperations(db *sql.DB) {
	//CRUD

	//Create
	statement, _ := db.Prepare("INSERT INTO books(isbn,name,author) VALUES(?,?,?)")
	statement.Exec("124566", "A Tale of Two cities", "Charles Dickens")
	log.Println("Inserted book into database")

	//READ
	rows, _ := db.Query("SELECT id, name, author FROM books")
	var tempbook Book
	for rows.Next() {
		rows.Scan(&tempbook.id, &tempbook.name, &tempbook.author)
		log.Printf("ID:%d, Book:%s, Author:%s\n", tempbook.id, tempbook.name, tempbook.author)
	}

	//Update
	statement, _ = db.Prepare("UPDATE books set name=? where id=?")
	statement.Exec("A tale of One city", 1)
	log.Println("Update name of book")

	//Delete
	statement, _ = db.Prepare("delete from books where id=?")
	statement.Exec(1)
	log.Println("Successfully deleted the book in database!")
}
