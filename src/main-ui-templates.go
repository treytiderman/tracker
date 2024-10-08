package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"text/template"

	bf_chroma "github.com/Depado/bfchroma/v2"
	bf_html "github.com/alecthomas/chroma/v2/formatters/html"
	bf "github.com/russross/blackfriday/v2"
)

func Routes_pages(db *sql.DB) {
	page_Settings(db)
	page_Trackers(db)
	page_Tracker_Create(db)
	page_Tracker_Info(db)
	page_Tracker_Log(db)
	page_Tracker_Records(db)
	page_Tracker_History(db)
	page_Entry_View(db)
}

func page_Trackers(db *sql.DB) {
	funcMap := template.FuncMap{
		"increment": func(i int) int {
			return i + 1
		},
	}

	tmp, err := template.New("").Funcs(funcMap).ParseFiles("./templates/trackers.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/trackers", func(w http.ResponseWriter, r *http.Request) {

		// Get All Trackers
		trackers, err := Db_Tracker_All_Get(db)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("GET: /trackers")

		// Get Records
		entries, err := Db_Entry_All_Get(db)
		if err != nil {
			log.Fatal(err)
			w.Write([]byte(err.Error()))
			return
		}

		// Page Data
		data := struct {
			Trackers []Db_Tracker
			Entries  []Db_Entry
		}{
			Trackers: trackers,
			Entries:  entries,
		}

		tmp.ExecuteTemplate(w, "trackers.html", data)
	})
}

func page_Tracker_Info(db *sql.DB) {
	funcMap := template.FuncMap{
		"decimal_places_to_step_size": func(x int) float32 {
			return 1 / float32(math.Pow10(x))
		},
	}

	tmp, err := template.New("").Funcs(funcMap).ParseFiles("./templates/tracker-info.html")
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
			if len(trackers) > 0 {
				// Set id to first tracker's id if not set in the URL
				id = trackers[0].Id
			} else {
				id = 1
			}
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

		tmp.ExecuteTemplate(w, "tracker-info.html", data)
	})
}

func page_Tracker_Create(db *sql.DB) {
	tmp, err := template.New("").ParseFiles("./templates/tracker-create.html")
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

		tmp.ExecuteTemplate(w, "tracker-create.html", data)
	})
}

func page_Tracker_Log(db *sql.DB) {
	funcMap := template.FuncMap{
		"decimal_places_to_step_size": func(x int) float32 {
			return 1 / float32(math.Pow10(x))
		},
	}

	tmp, err := template.New("").Funcs(funcMap).ParseFiles("./templates/tracker-log.html")
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
			if len(trackers) > 0 {
				// Set id to first tracker's id if not set in the URL
				id = trackers[0].Id
			} else {
				id = 1
			}
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

		tmp.ExecuteTemplate(w, "tracker-log.html", data)
	})
}

func page_Tracker_Records(db *sql.DB) {
	tmp, err := template.New("").ParseFiles("./templates/tracker-records.html")
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
			if len(trackers) > 0 {
				// Set id to first tracker's id if not set in the URL
				id = trackers[0].Id
			} else {
				id = 1
			}
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

		tmp.ExecuteTemplate(w, "tracker-records.html", data)
	})
}

func page_Tracker_History(db *sql.DB) {
	funcMap := template.FuncMap{
		"decimal_places_to_step_size": func(x int) float32 {
			return 1 / float32(math.Pow10(x))
		},
		"render_markdown": func(md string) string {
			var b []byte
			for _, bb := range []byte(md) {
				// Parser doesn't like \r (byte: 13)
				if bb != 13 {
					b = append(b, bb)
				}
			}
			arr := bf.Run([]byte(b), bf.WithRenderer(bf_chroma.NewRenderer(
				bf_chroma.Style("vulcan"),
				bf_chroma.ChromaOptions(bf_html.WithLineNumbers(true), bf_html.WithClasses(true)),
			)), bf.WithExtensions(bf.HardLineBreak | bf.CommonExtensions))
			str := string(arr)
			return str
		},
	}

	tmp, err := template.New("").Funcs(funcMap).ParseFiles("./templates/tracker-history.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/tracker-history", func(w http.ResponseWriter, r *http.Request) {

		// Get All Trackers
		trackers, err := Db_Tracker_All_Get(db)
		if err != nil {
			log.Fatal(err)
		}

		// Get Id from URL
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			if len(trackers) > 0 {
				// Set id to first tracker's id if not set in the URL
				id = trackers[0].Id
			} else {
				id = 1
			}
		}
		fmt.Printf("GET: /tracker-history?id=%d\n", id)

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

		tmp.ExecuteTemplate(w, "tracker-history.html", data)
	})
}

func page_Entry_View(db *sql.DB) {
	funcMap := template.FuncMap{
		"decimal_places_to_step_size": func(x int) float32 {
			return 1 / float32(math.Pow10(x))
		},
		"render_markdown": func(md string) string {
			var b []byte
			for _, bb := range []byte(md) {
				// Parser doesn't like \r (byte: 13)
				if bb != 13 {
					b = append(b, bb)
				}
			}
			arr := bf.Run([]byte(b), bf.WithRenderer(bf_chroma.NewRenderer(
				bf_chroma.Style("vulcan"),
				bf_chroma.ChromaOptions(bf_html.WithLineNumbers(true), bf_html.WithClasses(true)),
			)), bf.WithExtensions(bf.HardLineBreak | bf.CommonExtensions))
			str := string(arr)
			return str
		},
	}

	tmp, err := template.New("").Funcs(funcMap).ParseFiles("./templates/entry-view.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/entry-view", func(w http.ResponseWriter, r *http.Request) {

		tracker_id, err := strconv.Atoi(r.URL.Query().Get("tracker_id"))
		if err != nil {
			tracker_id = 0
		}

		entry_id, err := strconv.Atoi(r.URL.Query().Get("entry_id"))
		if err != nil {
			entry_id = 0
		}

		fmt.Printf("GET: /entry-view?tracker_id=%dentry_id=%d\n", tracker_id, entry_id)

		// Get Records by Id
		entries, err := Db_Entry_Get(db, tracker_id)
		if err != nil {
			log.Fatal(err)
			return
		}

		var entry Db_Entry
		for _, e := range entries {
			if e.Id == entry_id {
				entry = e
			}
		}

		// Page Data
		data := struct {
			Entry  Db_Entry
		}{
			Entry:  entry,
		}

		tmp.ExecuteTemplate(w, "entry-view.html", data)
	})
}

func page_Settings(db *sql.DB) {
	tmp, err := template.New("").ParseFiles("./templates/settings.html")
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

		tmp.ExecuteTemplate(w, "settings.html", data)
	})
}
