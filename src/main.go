package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func main() {
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// slog.SetDefault(logger)

	db_path := os.Getenv("DB_PATH")
	if db_path == "" {
		db_path = "../data/data.db"
	}

	var err error
	db, err = sql.Open("sqlite", db_path)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("DATABASE opened: %s\n", db_path)
	defer db.Close()

	err = Create_Tracker_Tables(db)
	if err != nil {
		fmt.Print(err)
	}

	err = Create_Entry_Tables(db)
	if err != nil {
		fmt.Print(err)
	}

	trackers, err := Get_Trackers(db)
	if err != nil {
		fmt.Print(err)
	}

	// Create a default tracker if none exist
	if len(trackers) == 0 {
		tracker_id, err := Create_Tracker(db, "Notes", "Notes, Memos, Journal, etc.")
		if err != nil {
			fmt.Print(err)
		}

		_, err = Create_Entry(db, tracker_id,
`# Heading H1

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
			fmt.Print(err)
		}
	}

	http_server_start()
}
