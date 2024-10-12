package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func Routes_upload() {
	http.HandleFunc("/ui/upload", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		if method == "GET" {
			fmt.Printf("GET: %s\n", r.URL)
		}
		if method == "DELETE" {
			fmt.Printf("DELETE: %s\n", r.URL)
		}

		READ_ONLY := os.Getenv("READ_ONLY")
		if READ_ONLY == "true" {
			w.Write([]byte("Nope, this is READ_ONLY mode"))
			return
		}

		fmt.Printf("POST: %s\n", r.URL)

		timestamp := time.Now().Format(time.DateTime)
		timestamp = strings.ReplaceAll(timestamp, " ", "_")
		path := "../public/upload/" + timestamp + ".png"
		img, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		defer img.Close()
		fmt.Printf("FILE SAVED: %s\n", path)

		_, err = io.Copy(img, r.Body)
		if err != nil {
			log.Fatal(err)
		}

		path = strings.ReplaceAll(path, "../", "/")
		w.Write([]byte(path))
	})
}
