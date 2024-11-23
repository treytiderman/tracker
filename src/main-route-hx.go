package main

import (
	"log"
	"net/http"
)

func handle_routes_ui(mux *http.ServeMux) {
	mux.Handle("/hx", mw_logger(mw_auth(http.HandlerFunc(hx_home_page))))
	mux.Handle("/hx/search", mw_logger(mw_auth(http.HandlerFunc(hx_search_results))))

	mux.Handle("/hx/new", mw_logger(mw_auth(http.HandlerFunc(hx_home_page))))
	mux.Handle("/hx/edit", mw_logger(mw_auth(http.HandlerFunc(hx_home_page))))

	mux.Handle("/hx/hello", mw_logger(mw_auth(http.HandlerFunc(hx_hello))))
}

func hx_hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello"))
}

func hx_home_page(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-hx")

	tmp.ExecuteTemplate(w, "app_page_only", struct {
		Title string
	}{
		Title: "Log",
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

	entries, err := Db_Entry_All_Filter_Notes_Get(db, search)
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
