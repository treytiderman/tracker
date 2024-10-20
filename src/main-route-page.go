package main

import (
	"log"
	"math"
	"net/http"
	"strconv"
	"text/template"

	bf_chroma "github.com/Depado/bfchroma/v2"
	bf_html "github.com/alecthomas/chroma/v2/formatters/html"
	bf "github.com/russross/blackfriday/v2"
)

func handle_routes_page(mux *http.ServeMux) {
	mux.Handle("/login", mw_logger(http.HandlerFunc(page_login)))

	mux.Handle("/trackers", mw_logger(mw_auth(http.HandlerFunc(page_trackers))))
	mux.Handle("/tracker-create", mw_logger(mw_auth(http.HandlerFunc(page_tracker_create))))

	mux.Handle("/tracker-info", mw_logger(mw_auth(http.HandlerFunc(page_tracker_info))))
	mux.Handle("/tracker-log", mw_logger(mw_auth(http.HandlerFunc(page_tracker_log))))
	mux.Handle("/tracker-records", mw_logger(mw_auth(http.HandlerFunc(page_tracker_records))))
	mux.Handle("/tracker-history", mw_logger(mw_auth(http.HandlerFunc(page_tracker_history))))

	mux.Handle("/entry-view", mw_logger(mw_auth(http.HandlerFunc(page_entry_view))))

	mux.Handle("/settings", mw_logger(mw_auth(http.HandlerFunc(page_settings))))
	mux.Handle("/test", mw_logger(mw_auth(http.HandlerFunc(page_test))))
}

func parse_templates(page string) *template.Template {
	funcMap := template.FuncMap{
		"increment": func(i int) int {
			return i + 1
		},
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
				bf_chroma.ChromaOptions(
					bf_html.WithLineNumbers(false),
					bf_html.WithClasses(true),
				),
			)), bf.WithExtensions(bf.HardLineBreak|bf.CommonExtensions))
			str := string(arr)
			return str
		},
	}

	tmp, err := template.New("").Funcs(funcMap).ParseGlob("./components/*")
	if err != nil {
		log.Fatal(err)
	}

	tmp, err = tmp.ParseFiles("./pages/" + page + ".html")
	if err != nil {
		log.Fatal(err)
	}

	return tmp
}

// Pages

func page_login(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-login")
	tmp.ExecuteTemplate(w, "app_page_only", struct {
		Title string
	}{
		Title: "Login",
	})
}

func page_trackers(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-trackers")

	trackers, err := Db_Tracker_All_Get(db)
	if err != nil {
		log.Fatal(err)
	}

	entries, err := Db_Entry_All_Get(db)
	if err != nil {
		log.Fatal(err)
		return
	}

	tmp.ExecuteTemplate(w, "app", struct {
		Title    string
		Trackers []Db_Tracker
		Tracker  Db_Tracker
		Entries  []Db_Entry
	}{
		Title:    "Trackers",
		Trackers: trackers,
		Tracker:  Db_Tracker{Id: 1, Name: ""},
		Entries:  entries,
	})
}

func page_tracker_create(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-tracker-create")

	trackers, err := Db_Tracker_All_Get(db)
	if err != nil {
		log.Fatal(err)
	}

	tmp.ExecuteTemplate(w, "app", struct {
		Title    string
		Trackers []Db_Tracker
		Tracker  Db_Tracker
	}{
		Title:    "New Tracker",
		Trackers: trackers,
		Tracker:  Db_Tracker{Id: 1, Name: ""},
	})
}

func page_tracker_info(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-tracker-info")

	trackers, err := Db_Tracker_All_Get(db)
	if err != nil {
		log.Fatal(err)
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		if len(trackers) > 0 {
			id = trackers[0].Id
		} else {
			id = 1
		}
	}

	tracker, err := Db_Tracker_Get(db, id)
	if err != nil {
		log.Fatal(err)
	}

	tmp.ExecuteTemplate(w, "app", struct {
		Title    string
		Tracker  Db_Tracker
		Trackers []Db_Tracker
	}{
		Title:    "Info / " + tracker.Name,
		Trackers: trackers,
		Tracker:  tracker,
	})
}

func page_tracker_log(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-tracker-log")

	trackers, err := Db_Tracker_All_Get(db)
	if err != nil {
		log.Fatal(err)
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		if len(trackers) > 0 {
			id = trackers[0].Id
		} else {
			id = 1
		}
	}

	tracker, err := Db_Tracker_Get(db, id)
	if err != nil {
		log.Fatal(err)
	}

	entries, err := Db_Entry_Get(db, id)
	if err != nil {
		log.Fatal(err)
		return
	}

	tmp.ExecuteTemplate(w, "app", struct {
		Title    string
		Trackers []Db_Tracker
		Tracker  Db_Tracker
		Entries  []Db_Entry
	}{
		Title:    "Log / " + tracker.Name,
		Trackers: trackers,
		Tracker:  tracker,
		Entries:  entries,
	})
}

func page_tracker_records(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-tracker-records")

	trackers, err := Db_Tracker_All_Get(db)
	if err != nil {
		log.Fatal(err)
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		if len(trackers) > 0 {
			id = trackers[0].Id
		} else {
			id = 1
		}
	}

	tracker, err := Db_Tracker_Get(db, id)
	if err != nil {
		log.Fatal(err)
	}

	entries, err := Db_Entry_Get(db, id)
	if err != nil {
		log.Fatal(err)
		return
	}

	tmp.ExecuteTemplate(w, "app", struct {
		Title    string
		Trackers []Db_Tracker
		Tracker  Db_Tracker
		Entries  []Db_Entry
	}{
		Title:    "Records / " + tracker.Name,
		Trackers: trackers,
		Tracker:  tracker,
		Entries:  entries,
	})
}

func page_tracker_history(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-tracker-history")

	trackers, err := Db_Tracker_All_Get(db)
	if err != nil {
		log.Fatal(err)
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		if len(trackers) > 0 {
			id = trackers[0].Id
		} else {
			id = 1
		}
	}

	tracker, err := Db_Tracker_Get(db, id)
	if err != nil {
		log.Fatal(err)
	}

	entries, err := Db_Entry_Get(db, id)
	if err != nil {
		log.Fatal(err)
		return
	}

	tmp.ExecuteTemplate(w, "app", struct {
		Title    string
		Trackers []Db_Tracker
		Tracker  Db_Tracker
		Entries  []Db_Entry
	}{
		Title:    "History / " + tracker.Name,
		Trackers: trackers,
		Tracker:  tracker,
		Entries:  entries,
	})
}

func page_entry_view(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-entry-view")

	tracker_id, err := strconv.Atoi(r.URL.Query().Get("tracker_id"))
	if err != nil {
		tracker_id = 0
	}

	entry_id, err := strconv.Atoi(r.URL.Query().Get("entry_id"))
	if err != nil {
		entry_id = 0
	}

	log.Printf("GET: /entry-view?tracker_id=%dentry_id=%d\n", tracker_id, entry_id)

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

	tmp.ExecuteTemplate(w, "app_page_only", struct {
		Title string
		Entry Db_Entry
	}{
		Title: "Entry",
		Entry: entry,
	})
}

func page_settings(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-settings")

	trackers, err := Db_Tracker_All_Get(db)
	if err != nil {
		log.Fatal(err)
	}

	tmp.ExecuteTemplate(w, "app", struct {
		Title    string
		Tracker  Db_Tracker
		Trackers []Db_Tracker
	}{
		Title:    "Trackers",
		Trackers: trackers,
		Tracker:  Db_Tracker{Id: 1, Name: ""},
	})
}

func page_test(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-test")
	tmp.ExecuteTemplate(w, "app", struct {
		Title string
	}{
		Title: "Test",
	})
}
