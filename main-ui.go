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
	// http.Handle("/", fs)
	http.Handle("/public/", http.StripPrefix("/public/", fs))

	// Base URL Redirect
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/trackers", http.StatusSeeOther)
	})

	// Test Route
	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, time.Now().Format(time.UnixDate))
	})

	// All Other Routes
	routes(db)

	// Start Web Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	port = fmt.Sprintf(":%s", port)
	hostname, _ := os.Hostname()

	fmt.Println("Web Server: started")
	fmt.Printf("- http://%s%s\n", "localhost", port)
	fmt.Printf("- http://%s%s\n", hostname, port)
	fmt.Println()

	log.Fatal(http.ListenAndServe(port, nil))
}

func routes(db *sql.DB) {
	page_Trackers(db)
	page_Tracker_Create(db)
	page_Tracker_Records(db)
	page_Tracker_Log(db)
	page_Tracker_Chart(db)
}

func page_Trackers(db *sql.DB) {
	t, err1 := template.New("").ParseFS(Templates_Embed, "templates/trackers.html")
	if err1 != nil {
		log.Fatal(err1)
	}

	http.HandleFunc("/trackers", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET: /trackers")
		data := []Db_Tracker{
			{
				1, "Walk Dog 🐶", "Each time I take the dog out", nil,
			},
			{
				2, "Brush Teeth 🦷", "Each time I brush my teeth", nil,
			},
			{
				3, "Money 💰", "Money used by different cards", []Db_Field{
					{1, "number", "Transactions", "Money Spent", Db_Number{1, 2}, nil},
					{2, "option", "Card", "Which Credit or Debit Card was used?", Db_Number{0, 0}, []Db_Option{
						{1, -1, "Discover"},
						{2, 0, "Visa"},
						{3, 2, "American Express"},
					}},
				},
			},
			{
				4, "Lift 🏋️", "My Lifting Habits", []Db_Field{
					{2, "option", "Exercise", "", Db_Number{0, 0}, []Db_Option{
						{1, 1, "Bench Press"},
						{2, 3, "Squat"},
						{3, 5, "Deadlift"},
					}},
					{3, "number", "Weight", "", Db_Number{2, 0}, nil},
					{4, "number", "Reps", "", Db_Number{3, 0}, nil},
				},
			},
		}

		trackers, err2 := Get_Trackers(db)
		if err2 != nil {
			log.Fatal(err2)
		}
		data = trackers

		t.ExecuteTemplate(w, "trackers.html", data)
	})

	http.HandleFunc("/htmx/tracker/delete", func(w http.ResponseWriter, r *http.Request) {
		id, err1 := strconv.Atoi(r.URL.Query().Get("id"))
		if err1 != nil {
			w.Write([]byte(err1.Error()))
			return
		}

		err2 := Delete_Tracker_By_Id(db, id)
		if err2 != nil {
			w.Write([]byte(err2.Error()))
			return
		}

		w.Write([]byte("<div class='text-red-300'>Deleted</div>"))
	})
}

func page_Tracker_Create(db *sql.DB) {
	page, err1 := template.New("").ParseFS(Templates_Embed, "templates/tracker-create.html")
	if err1 != nil {
		log.Fatal(err1)
	}

	http.HandleFunc("/tracker-create", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET: /tracker-create")
		page.ExecuteTemplate(w, "tracker-create.html", "")
	})

	http.HandleFunc("/htmx/tracker/create", func(w http.ResponseWriter, r *http.Request) {
		err1 := r.ParseForm()
		if err1 != nil {
			w.Write([]byte(err1.Error()))
			return
		}

		tracker_name := r.Form.Get("tracker_name")
		tracker_notes := r.Form.Get("tracker_notes")

		_, err2 := Create_Tracker(db, tracker_name, tracker_notes)
		if err2 != nil {
			w.Write([]byte(err2.Error()))
			return
		}

		for field_id := 0; field_id < 100; field_id++ {
			if r.Form.Has(fmt.Sprintf("field_%d_name", field_id)) {
				field_name := r.Form.Get(fmt.Sprintf("field_%d_name", field_id))
				field_type := r.Form.Get(fmt.Sprintf("field_%d_type", field_id))
				// field_notes := r.Form.Get(fmt.Sprintf("field_%d_notes", field_id))
				field_notes := ""

				if field_type == "number" {
					decimal_places := 0

					if r.Form.Has(fmt.Sprintf("field_%d_decimal_places", field_id)) {
						decimal_places, _ = strconv.Atoi(r.Form.Get(fmt.Sprintf("field_%d_decimal_places", field_id)))
					}

					_, err4 := Add_Number_Field(db, tracker_name, field_name, field_notes, decimal_places)
					if err4 != nil {
						w.Write([]byte(err4.Error()))
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

					_, err5 := Add_Option_Field(db, tracker_name, field_name, field_notes, options)
					if err5 != nil {
						w.Write([]byte(err5.Error()))
						return
					}
				}
			}
		}

		http.Redirect(w, r, "/trackers", http.StatusSeeOther)
	})
}

func page_Tracker_Records(db *sql.DB) {
	t, err := template.New("").ParseFS(Templates_Embed, "templates/tracker-records.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/tracker/records", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			log.Fatal(err)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Printf("GET: /tracker/records?id=%d\n", id)

		data := struct {
			Tracker  Db_Tracker
			Entries  []Db_Entry
			Trackers []Db_Tracker
		}{
			// Tracker: Tracker{
			// 	Id:    1,
			// 	Name:  "Commissions",
			// 	Notes: "Badger badger badger",
			// },
			// Fields: []Field_Deep{
			// 	{
			// 		Id: 1, Type: "number", Name: "Cost", Notes: "",
			// 		Type_Number: Field_Number{ Decimal_Places: 2 },
			// 		Type_Option: Field_Option{},
			// 	},
			// 	{
			// 		Id: 1, Type: "option", Name: "Status", Notes: "",
			// 		Type_Number: Field_Number{},
			// 		Type_Option: Field_Option{
			// 			Option_Values: []int{-1, 0, 1},
			// 			Option_Names:  []string{"Canceled", "In Progress", "Complete"},
			// 		},
			// 	},
			// },
			// Records: []struct {
			// 	Id        int
			// 	Timestamp string
			// 	Data      []string
			// 	Notes     string
			// }{
			// 	{
			// 		Id:        1,
			// 		Timestamp: "2024-05-07T18:56:44Z",
			// 		Data:      []string{"1000.00", "Complete"},
			// 		Notes:     "notes...",
			// 	},
			// 	{
			// 		Id:        2,
			// 		Timestamp: "2024-06-24T19:05:12Z",
			// 		Data:      []string{"250.50", "In Progress"},
			// 		Notes:     "notes 2...",
			// 	},
			// },
		}

		tracker, err := Get_Tracker_By_Id(db, id)
		if err != nil {
			log.Fatal(err)
			w.Write([]byte(err.Error()))
			return
		}
		data.Tracker = tracker

		entries, err := Get_Entries_By_Tracker_Id(db, id)
		if err != nil {
			log.Fatal(err)
			w.Write([]byte(err.Error()))
			return
		}
		data.Entries = entries

		trackers, err := Get_Trackers(db)
		if err != nil {
			log.Fatal(err)
		}
		data.Trackers = trackers

		// j, _ := json.Marshal(data.Entries)
		// fmt.Println("- DATA:", string(j))
		t.ExecuteTemplate(w, "tracker-records.html", data)
	})
}

func page_Tracker_Log(db *sql.DB) {
	funcMap := template.FuncMap{
		"decimal_places_to_step_size": func(x int) float32 {
			return 1 / float32(math.Pow10(x))
		},
	}

	t, err1 := template.New("").Funcs(funcMap).ParseFS(Templates_Embed, "templates/tracker-log.html")
	if err1 != nil {
		log.Fatal(err1)
	}

	http.HandleFunc("/tracker/log", func(w http.ResponseWriter, r *http.Request) {
		id, err1 := strconv.Atoi(r.URL.Query().Get("id"))
		if err1 != nil {
			w.Write([]byte(err1.Error()))
			return
		}
		fmt.Printf("GET: /tracker/log?id=%d\n", id)

		data := struct {
			Trackers []Db_Tracker
			Tracker  Db_Tracker
		}{}

		tracker, err2 := Get_Tracker_By_Id(db, id)
		if err2 != nil {
			log.Fatal(err2)
			w.Write([]byte(err2.Error()))
			return
		}
		data.Tracker = tracker

		trackers, err4 := Get_Trackers(db)
		if err4 != nil {
			log.Fatal(err4)
		}
		data.Trackers = trackers

		t.ExecuteTemplate(w, "tracker-log.html", data)
	})

	http.HandleFunc("/htmx/tracker/record", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("POST: %s\n", r.URL)

		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		tracker, err := Get_Tracker_By_Id(db, id)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		// fmt.Println("- tracker", tracker)

		err = r.ParseForm()
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Printf("FORM: %s\n", r.Form.Encode())

		var logs = make([]struct {
			Field_Id int
			Value    int
		}, 0)
		entry_notes := r.Form.Get("entry_notes")
		r.Form.Del("entry_notes")
		r.Form.Del("id")
		// fmt.Println("- entry_notes", entry_notes)

		for k, v := range r.Form {
			// fmt.Println("- raw", k, v[0])

			field_id, err := strconv.Atoi(strings.ReplaceAll(k, "field_", ""))
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			// fmt.Println("- field_id", field_id)

			var field Db_Field
			for _, f := range tracker.Fields {
				if f.Id == field_id {
					field = f
				}
			}
			// fmt.Println("- field", field)

			value := 0
			if field.Type == "number" {
				field_value_float, _ := strconv.ParseFloat(v[0], 64)
				field_value_adjusted := float64(field_value_float) * float64(math.Pow10(field.Number.Decimal_Places))
				value = int(math.Floor(field_value_adjusted))
			} else if field.Type == "option" {
				value, err = strconv.Atoi(v[0])
				if err != nil {
					w.Write([]byte(err.Error()))
					return
				}
			}
			// fmt.Println("- value", value)

			logs = append(logs, struct {
				Field_Id int
				Value    int
			}{
				field_id,
				value,
			})
		}

		// fmt.Printf("> Add_Entry %s %+v %s\n", tracker.Name, logs, entry_notes)
		Add_Entry(db, tracker.Name, entry_notes, logs)

		w.Write([]byte("ok"))
	})
}

func page_Tracker_Chart(db *sql.DB) {
	t, err := template.New("").ParseFS(Templates_Embed, "templates/tracker-chart.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/tracker/chart", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Printf("GET: /tracker/chart?id=%d\n", id)

		data := struct {
			Trackers []Db_Tracker
			Tracker  Db_Tracker
			Entries  []Db_Entry
		}{}

		tracker, err := Get_Tracker_By_Id(db, id)
		if err != nil {
			log.Fatal(err)
			w.Write([]byte(err.Error()))
			return
		}
		data.Tracker = tracker

		trackers, err := Get_Trackers(db)
		if err != nil {
			log.Fatal(err)
		}
		data.Trackers = trackers

		entries, err := Get_Entries_By_Tracker_Id(db, id)
		if err != nil {
			log.Fatal(err)
			w.Write([]byte(err.Error()))
			return
		}
		data.Entries = entries

		t.ExecuteTemplate(w, "tracker-chart.html", data)
	})
}
