package main

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func main() {
	db_path := os.Getenv("DB_PATH")
	if db_path == "" {
		db_path = "../data/data.db"
	}

	var err error
	db, err = sql.Open("sqlite", db_path)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("DATABASE SQLite: %s\n", db_path)
	defer db.Close()

	err = Create_Tracker_Tables(db)
	if err != nil {
		log.Fatal(err)
	}

	err = Create_Entry_Tables(db)
	if err != nil {
		log.Fatal(err)
	}

	http_server_start()
}
