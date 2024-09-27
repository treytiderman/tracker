package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"text/template"
)

func Routes_pages(db *sql.DB) {
	page_Settings(db)
	page_Tracker_Create(db)
	page_Tracker_Info(db)
	page_Tracker_Log(db)
	page_Tracker_Records(db)
}

func page_Tracker_Info(db *sql.DB) {
	funcMap := template.FuncMap{
		"decimal_places_to_step_size": func(x int) float32 {
			return 1 / float32(math.Pow10(x))
		},
	}

	t, err := template.New("").Funcs(funcMap).ParseFS(Templates_Embed, "templates/tracker-info.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/tracker-info", func(w http.ResponseWriter, r *http.Request) {

		// Get All Trackers
		trackers, err := Db_Tracker_All_Get(db)
		if err != nil {
			log.Fatal(err)
		}

		// Get Id from URL
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			// Set id to first tracker's id if not set in the URL
			id = trackers[0].Id
		}
		fmt.Printf("GET: /tracker-info?id=%d\n", id)

		// Get Tracker by Id
		tracker, err := Db_Tracker_Get(db, id)
		if err != nil {
			log.Fatal(err)
		}

		// Page Data
		data := struct {
			Trackers []Db_Tracker
			Tracker  Db_Tracker
		}{
			Tracker:  tracker,
			Trackers: trackers,
		}

		t.ExecuteTemplate(w, "tracker-info.html", data)
	})
}

func page_Tracker_Create(db *sql.DB) {
	t, err := template.New("").ParseFS(Templates_Embed, "templates/tracker-create.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/tracker-create", func(w http.ResponseWriter, r *http.Request) {

		// Get All Trackers
		trackers, err := Db_Tracker_All_Get(db)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("GET: /tracker-create")

		// Page Data
		data := struct {
			Trackers []Db_Tracker
		}{
			Trackers: trackers,
		}

		t.ExecuteTemplate(w, "tracker-create.html", data)
	})
}

func page_Tracker_Log(db *sql.DB) {
	funcMap := template.FuncMap{
		"decimal_places_to_step_size": func(x int) float32 {
			return 1 / float32(math.Pow10(x))
		},
	}

	t, err := template.New("").Funcs(funcMap).ParseFS(Templates_Embed, "templates/tracker-log.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/tracker-log", func(w http.ResponseWriter, r *http.Request) {

		// Get All Trackers
		trackers, err := Db_Tracker_All_Get(db)
		if err != nil {
			log.Fatal(err)
		}

		// Get Id from URL
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			// Set id to first tracker's id if not set in the URL
			id = trackers[0].Id
		}
		fmt.Printf("GET: /tracker-log?id=%d\n", id)

		// Get Tracker by Id
		tracker, err := Db_Tracker_Get(db, id)
		if err != nil {
			log.Fatal(err)
		}

		// Page Data
		data := struct {
			Trackers []Db_Tracker
			Tracker  Db_Tracker
		}{
			Tracker:  tracker,
			Trackers: trackers,
		}

		t.ExecuteTemplate(w, "tracker-log.html", data)
	})
}

func page_Tracker_Records(db *sql.DB) {
	t, err := template.New("").ParseFS(Templates_Embed, "templates/tracker-records.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/tracker-records", func(w http.ResponseWriter, r *http.Request) {

		// Get All Trackers
		trackers, err := Db_Tracker_All_Get(db)
		if err != nil {
			log.Fatal(err)
		}

		// Get Id from URL
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			// Set id to first tracker's id if not set in the URL
			id = trackers[0].Id
		}
		fmt.Printf("GET: /tracker-records?id=%d\n", id)

		// Get Tracker by Id
		tracker, err := Db_Tracker_Get(db, id)
		if err != nil {
			log.Fatal(err)
		}

		// Get Records by Id
		entries, err := Db_Entry_Get(db, id)
		if err != nil {
			log.Fatal(err)
			w.Write([]byte(err.Error()))
			return
		}

		// Page Data
		data := struct {
			Tracker  Db_Tracker
			Entries  []Db_Entry
			Trackers []Db_Tracker
		}{
			Trackers: trackers,
			Tracker:  tracker,
			Entries:  entries,
		}

		t.ExecuteTemplate(w, "tracker-records.html", data)
	})
}

func page_Settings(db *sql.DB) {
	t, err := template.New("").ParseFS(Templates_Embed, "templates/settings.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {

		// Get All Trackers
		trackers, err := Db_Tracker_All_Get(db)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("GET: /settings")

		// Page Data
		data := struct {
			Trackers []Db_Tracker
		}{
			Trackers: trackers,
		}

		t.ExecuteTemplate(w, "settings.html", data)
	})
}
