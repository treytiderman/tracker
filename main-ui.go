package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"
)

//go:embed public
var Public_Embed embed.FS

//go:embed templates
var Templates_Embed embed.FS

func Start_Web_Server(db *sql.DB) {

	// Setup Public Routes
	var public_fs = fs.FS(Public_Embed)
	public_files, err := fs.Sub(public_fs, "public")
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.FS(public_files))
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	// Base URL Redirect
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/tracker-info", http.StatusSeeOther)
		}
	})

	// Test Route
	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, time.Now().Format(time.UnixDate))
	})

	// All Other Routes
	routes(db)

	// Start Web Server
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8000"
	}
	port = fmt.Sprintf(":%s", port)
	fmt.Printf("HTTP SERVER: http://%s%s\n", "localhost", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func routes(db *sql.DB) {
	htmx_Routes(db)
	page_Settings(db)
	page_Tracker_Info(db)
	page_Tracker_Create(db)
	page_Tracker_Log(db)
	page_Tracker_Records(db)
}

func htmx_Routes(db *sql.DB) {
	http.HandleFunc("/htmx/tracker/create", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("POST: %s\n", r.URL)

		err := r.ParseForm()
		if err != nil {
			return
		}
		fmt.Printf("FORM: %s\n", r.Form.Encode())

		tracker_name := r.Form.Get("tracker_name")
		tracker_notes := r.Form.Get("tracker_notes")

		tracker_id, err := Db_Tracker_Create(db, tracker_name, tracker_notes)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		for field_id := 0; field_id < 100; field_id++ {
			if r.Form.Has(fmt.Sprintf("field_%d_name", field_id)) {
				field_name := r.Form.Get(fmt.Sprintf("field_%d_name", field_id))
				field_type := r.Form.Get(fmt.Sprintf("field_%d_type", field_id))
				field_notes := ""

				if field_type == "number" {
					decimal_places := 0

					if r.Form.Has(fmt.Sprintf("field_%d_decimal_places", field_id)) {
						decimal_places, _ = strconv.Atoi(r.Form.Get(fmt.Sprintf("field_%d_decimal_places", field_id)))
					}

					_, err := Db_Tracker_Field_Number_Create(db, tracker_id, field_name, field_notes, decimal_places)
					if err != nil {
						return
					}
				} else if field_type == "option" {
					var options []struct {
						Value int
						Name  string
					}

					for option_id := 0; option_id < 100; option_id++ {
						if r.Form.Has(fmt.Sprintf("field_%d_option_%d_value", field_id, option_id)) {
							option_name := r.Form.Get(fmt.Sprintf("field_%d_option_%d_name", field_id, option_id))
							option_value, _ := strconv.Atoi(r.Form.Get(fmt.Sprintf("field_%d_option_%d_value", field_id, option_id)))
							options = append(options, struct {
								Value int
								Name  string
							}{
								option_value,
								option_name,
							})
						}
					}

					_, err := Db_Tracker_Field_Option_Create(db, tracker_id, field_name, field_notes, options)
					if err != nil {
						return
					}
				}
			}
		}

		url := fmt.Sprintf("/tracker-info?id=%d", tracker_id)
		http.Redirect(w, r, url, http.StatusSeeOther)
	})

	http.HandleFunc("/htmx/tracker/name", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("POST: %s\n", r.URL)

		// Parse Form Data
		err := r.ParseForm()
		if err != nil {
			return
		}
		fmt.Printf("FORM: %s\n", r.Form.Encode())

		// Get Id from URL
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			return
		}

		// Get Tracker Name from Form Data
		tracker_name := r.Form.Get("tracker_name")

		// Update Tracker Name
		err = Db_Tracker_Name_Update(db, id, tracker_name)
		if err != nil {
			return
		}

		// Reload page
		url := fmt.Sprintf("/tracker-info?id=%d", id)
		http.Redirect(w, r, url, http.StatusSeeOther)
	})

	http.HandleFunc("/htmx/tracker/notes", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("POST: %s\n", r.URL)

		// Parse Form Data
		err := r.ParseForm()
		if err != nil {
			return
		}
		fmt.Printf("FORM: %s\n", r.Form.Encode())

		// Get Id from URL
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			return
		}

		// Get Tracker Notes from Form Data
		tracker_notes := r.Form.Get("tracker_notes")
		fmt.Printf("POST: /htmx/tracker/notes?id=%d&tracker_notes=%s\n", id, tracker_notes)

		// Update Tracker Notes
		err = Db_Tracker_Notes_Update(db, id, tracker_notes)
		if err != nil {
			return
		}

		// Reload page
		url := fmt.Sprintf("/tracker-info?id=%d", id)
		http.Redirect(w, r, url, http.StatusSeeOther)
	})

	http.HandleFunc("/htmx/tracker/log", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("POST: %s\n", r.URL)

		// Parse Form Data
		err := r.ParseForm()
		if err != nil {
			return
		}
		fmt.Printf("FORM: %s\n", r.Form.Encode())

		// Get Id from URL
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			return
		}

		// Get Notes and delete extra fields from Form
		entry_notes := r.Form.Get("entry_notes")
		r.Form.Del("entry_notes")
		r.Form.Del("entry_date")
		r.Form.Del("entry_time")
		r.Form.Del("id")

		// Get Tracker by Id
		tracker, err := Db_Tracker_Get(db, id)
		if err != nil {
			log.Fatal(err)
		}

		// Logs Memory
		var logs = make([]struct {
			Field_Id int
			Value    int
		}, 0)

		for k, v := range r.Form {

			field_id, err := strconv.Atoi(strings.ReplaceAll(k, "field_", ""))
			if err != nil {
				return
			}

			var field Db_Field
			for _, f := range tracker.Fields {
				if f.Id == field_id {
					field = f
				}
			}

			value := 0
			if field.Type == "number" {
				field_value_float, _ := strconv.ParseFloat(v[0], 64)
				field_value_adjusted := float64(field_value_float) * float64(math.Pow10(field.Number.Decimal_Places))
				value = int(math.Floor(field_value_adjusted))
			} else if field.Type == "option" {
				value, err = strconv.Atoi(v[0])
				if err != nil {
					return
				}
			}

			logs = append(logs, struct {
				Field_Id int
				Value    int
			}{
				field_id,
				value,
			})
		}

		Db_Entry_Create(db, tracker.Id, entry_notes, logs)

		// Reload page
		url := fmt.Sprintf("/tracker-log?id=%d", id)
		http.Redirect(w, r, url, http.StatusSeeOther)
	})

	http.HandleFunc("/htmx/tracker/delete", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("POST: %s\n", r.URL)

		// Get Id from URL
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			return
		}

		// Delete Tracker
		err = Db_Tracker_Delete(db, id)
		if err != nil {
			return
		}

		// Reload without Id
		http.Redirect(w, r, "/tracker-info", http.StatusSeeOther)
	})
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
