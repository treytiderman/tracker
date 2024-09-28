package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func Start_Web_Server(db *sql.DB) {

	// Setup Public Routes
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	// Base URL Redirect
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/tracker-info", http.StatusSeeOther)
		}
	})

	// All Other Routes
	Routes_pages(db)
	Routes_htmx(db)

	// Test Route
	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, time.Now().Format(time.UnixDate))
	})

	// Start Web Server
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8000"
	}
	port = fmt.Sprintf(":%s", port)
	fmt.Printf("HTTP SERVER: http://%s%s\n", "localhost", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
