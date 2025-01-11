package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func handle_routes_ui(mux *http.ServeMux) {
	mux.Handle("/notes", mw_logger(mw_auth(http.HandlerFunc(notes_home_page))))
	mux.Handle("/notes/search", mw_logger(mw_auth(http.HandlerFunc(notes_search_results))))

	mux.Handle("/notes/entry", mw_logger(mw_auth(http.HandlerFunc(notes_entry))))

	mux.Handle("/notes/hello", mw_logger(mw_auth(http.HandlerFunc(notes_hello))))
}

func notes_hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func notes_home_page(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-notes")

	entry_id, err := strconv.Atoi(r.URL.Query().Get("entry"))
	if err != nil {
		entry_id = 0
	}
	log.Println("entry_id", entry_id)

	entries, err := Get_Entries(db, 1)
	if err != nil {
		log.Fatal(err)
	}

	var entry Db_Entry
	for _, ent := range entries {
		if ent.Id == entry_id {
			entry = ent
		}
	}

	tmp.ExecuteTemplate(w, "app_page_only", struct {
		Title string
		Entry Db_Entry
	}{
		Title: "Log",
		Entry: entry,
	})
}

func notes_search_results(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-notes")

	trackers, err := Get_Trackers(db)
	if err != nil {
		log.Fatal(err)
	}

	err = r.ParseForm()
	if err != nil {
		return
	}

	search := r.Form.Get("search")
	log.Println("SEARCH:", search)

	entries, err := Get_Entries_Filter(db, 1, search)
	if err != nil {
		log.Fatal(err)
	}

	tmp.ExecuteTemplate(w, "notes_history", struct {
		Search   string
		Tracker  Db_Tracker
		Trackers []Db_Tracker
		Entries  []Db_Entry
	}{
		Search:   search,
		Tracker:  Db_Tracker{Id: 1, Name: "Log"},
		Trackers: trackers,
		Entries:  entries,
	})
}

func notes_entry(w http.ResponseWriter, r *http.Request) {
	entry_id, err := strconv.Atoi(r.URL.Query().Get("entry"))
	if err != nil {
		log.Fatal(err)
	}

	err = r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	entry_note := r.Form.Get("notes")

	if entry_id == 0 {
		entry_id, err = Create_Entry(db, 1, entry_note)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err = Update_Entry_Notes(db, entry_id, entry_note)
		if err != nil {
			log.Fatal(err)
		}
	}

	url := fmt.Sprintf("/notes?entry=%d", entry_id)
	log.Println("url", url)
	w.Header().Add("HX-Redirect", url)
	w.Write([]byte("ok"))
}
