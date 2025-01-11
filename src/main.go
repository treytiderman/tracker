package main

import (
	"database/sql"
	"log"
	"log/slog"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func main() {

	// Set up structured logging | slog.NewJSONHandler or slog.NewTextHandler
	log_level_env := os.Getenv("LOG_LEVEL")
	log_level := slog.LevelDebug
	if log_level_env == "INFO" {
		log_level = slog.LevelInfo
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: log_level,
	})))

	// Set up database
	db_path := os.Getenv("DB_PATH")
	if db_path == "" {
		db_path = "../data/data.db"
	}

	var err error
	db, err = sql.Open("sqlite", db_path)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("database opened", "file", db_path)
	defer db.Close()

	// Initialize database table
	err = Create_Tracker_Tables(db)
	if err != nil {
		log.Fatal(err)
	}

	err = Create_Entry_Tables(db)
	if err != nil {
		log.Fatal(err)
	}

	trackers, err := Get_Trackers(db)
	if err != nil {
		log.Fatal(err)
	}

	// Create a default tracker if none exist
	if len(trackers) == 0 {
		tracker_id, err := Create_Tracker(db, "Notes", "Notes, Memos, Journal, etc.")
		if err != nil {
			log.Fatal(err)
		}

		_, err = Create_Entry(db, tracker_id, `# Heading H1

## Sub-Heading H2

Notes are writen with Markdown syntax
Text can be **bold** or *italic*
[Markdown Cheatsheet](https://www.markdownguide.org/cheat-sheet/)

- Unordered List Item 1
- Unordered List Item 2

1. Ordered List Item 1
2. Ordered List Item 2

- [x] Task List Item 1
- [ ] Task List Item 2
`)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = http_server_start()
	log.Fatal(err)
}
