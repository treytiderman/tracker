package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db_path := "./data/data.db"
	db, err := sql.Open("sqlite", db_path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("DATABASE SQLite: %s\n", db_path)
	defer db.Close()

	Create_Tables(db)

	Start_Web_Server(db)
}
