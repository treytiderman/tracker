package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func handle_routes_ui(mux *http.ServeMux) {
	mux.Handle("/hx", mw_logger(mw_auth(http.HandlerFunc(hx_home_page))))
	mux.Handle("/hx/search", mw_logger(mw_auth(http.HandlerFunc(hx_search_results))))

	mux.Handle("/hx/entry", mw_logger(mw_auth(http.HandlerFunc(hx_entry))))

	mux.Handle("/hx/hello", mw_logger(mw_auth(http.HandlerFunc(hx_hello))))
}

func hx_hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func hx_home_page(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-hx")

	entry_id, err := strconv.Atoi(r.URL.Query().Get("entry"))
	if err != nil {
		entry_id = 1
	}
	log.Println("entry_id", entry_id)

	entries, err := Db_Entry_All_Get(db)
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

func hx_search_results(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-hx")

	trackers, err := Db_Tracker_All_Get(db)
	if err != nil {
		log.Fatal(err)
	}

	err = r.ParseForm()
	if err != nil {
		return
	}

	search := r.Form.Get("search")
	log.Println("SEARCH:", search)

	entries, err := Db_Entry_Filter_Notes_Get(db, 1, search)
	if err != nil {
		log.Fatal(err)
	}

	tmp.ExecuteTemplate(w, "hx_search_results", struct {
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

func hx_entry(w http.ResponseWriter, r *http.Request) {
	entry_id, err := strconv.Atoi(r.URL.Query().Get("entry"))
	if err != nil {
		log.Fatal(err)
	}

	err = r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	entry_note := r.Form.Get("notes")

	logs := make([]struct{Field_Id int; Value int}, 0)

	if entry_id == 0 {
		entry_id, err = Db_Entry_Create(db, 1, entry_note, logs)
		if err != nil {
			log.Fatal(err)
		}
	} else {		
		err = Db_Entry_Notes_Update(db, entry_id, entry_note)
		if err != nil {
			log.Fatal(err)
		}
	}

	// url := fmt.Sprintf("/hx")
	url := fmt.Sprintf("/hx?entry=%d", entry_id)
	log.Println("url", url)
	w.Header().Add("HX-Redirect", url)
	w.Write([]byte("ok"))
}
