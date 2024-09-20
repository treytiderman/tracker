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
	defer db.Close()

	fmt.Printf("DATABASE SQLite: %s\n", db_path)

	Create_Db_Tables(db)

	Start_Web_Server(db)
}
