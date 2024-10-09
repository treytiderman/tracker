package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	db_path := os.Getenv("DB_PATH")
	if db_path == "" {
		db_path = "../data/data.db"
	}

	fmt.Printf("DATABASE SQLite: %s\n", db_path)
	db, err := sql.Open("sqlite", db_path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = Db_Tracker_Table_Create(db)
	if err != nil {
		log.Fatal(err)
	}

	err = Db_Entry_Table_Create(db)
	if err != nil {
		log.Fatal(err)
	}

	Start_Web_Server(db)
}
