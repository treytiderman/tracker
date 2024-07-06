package main

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("SQLITE file ./data.db opened")
	defer db.Close()

	Tables_Create(db)

	Start_Web_Server(db)
}
