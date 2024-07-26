package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database file (SQLite) './data.db' opened")
	fmt.Println()
	defer db.Close()

	Create_Tables(db)

	Start_Web_Server(db)
}
