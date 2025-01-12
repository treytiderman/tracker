package main

import (
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

func handle_routes_page(mux *http.ServeMux) {
	mux.Handle("/login", mw_logger(http.HandlerFunc(page_login)))

	mux.Handle("/trackers", mw_logger(mw_auth(http.HandlerFunc(page_trackers))))
	mux.Handle("/tracker-create", mw_logger(mw_auth(http.HandlerFunc(page_tracker_create))))

	mux.Handle("/tracker-info", mw_logger(mw_auth(http.HandlerFunc(page_tracker_info))))
	mux.Handle("/tracker-log", mw_logger(mw_auth(http.HandlerFunc(page_tracker_log))))
	mux.Handle("/tracker-records", mw_logger(mw_auth(http.HandlerFunc(page_tracker_records))))
	mux.Handle("/tracker-history", mw_logger(mw_auth(http.HandlerFunc(page_tracker_history))))
	mux.Handle("POST /htmx/entry/history", mw_logger(mw_auth(http.HandlerFunc(htmx_entry_history))))
	mux.Handle("POST /htmx/entry/checkbox-toggle", mw_logger(mw_auth(http.HandlerFunc(htmx_entry_checkbox_toggle))))

	mux.Handle("/entry-view", mw_logger(mw_auth(http.HandlerFunc(page_entry_view))))
	mux.Handle("/entry-editor", mw_logger(mw_auth(http.HandlerFunc(page_entry_editor))))

	mux.Handle("/content", mw_logger(mw_auth(http.HandlerFunc(page_content))))
	mux.Handle("/settings", mw_logger(mw_auth(http.HandlerFunc(page_settings))))
	mux.Handle("/test", mw_logger(mw_auth(http.HandlerFunc(page_test))))
}

func replace_with_checkbox(arr []byte, i int, task_count int) []byte {
	end_pattern_length := 5
	if i > len(arr)-end_pattern_length {
		return arr
	}

	before_i := 4
	after_i := 4
	start_tag := string(arr[i-before_i : i+after_i])

	if start_tag == "<li>[x] " || start_tag == "<li>[ ] " {
		task_count++

		checked := " "
		if start_tag == "<li>[x] " {
			checked = "checked"
		}

		j := 0
		j_offset := 0
		a := arr[i+after_i:]
		for j = 0; j < len(a); j++ {
			next_tag := string(a[j : j+end_pattern_length])
			if next_tag == "<br /" {
				j_offset = 7
				break
			} else if next_tag == "</li>" {
				break
			}
		}

		before := arr[:i-before_i]
		after := a[j+j_offset:]
		label := a[:j]

		div := fmt.Sprintf(`<li style="list-style: none;">
		<div class="flex items-baseline gap-4">
		    <input type="checkbox" class="tt-input" onclick="event.preventDefault();" id="task_%d" name="task_%d" %s >
		    <label class="tt-label-inline" style="font-size: var(--font-size);" for="task_%d"> %s </label>
		</div>`, task_count, task_count, checked, task_count, label)

		arr = []byte(string(before) + div + string(after))
	}

	i++
	return replace_with_checkbox(arr, i, task_count)
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

			arr = replace_with_checkbox(arr, 4, 0)

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

	trackers, err := Get_Trackers(db)
	if err != nil {
		log.Fatal(err)
	}

	entries, err := Get_All_Entries(db)
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

	trackers, err := Get_Trackers(db)
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

	trackers, err := Get_Trackers(db)
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

	tracker, err := Get_Tracker(db, id)
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

	trackers, err := Get_Trackers(db)
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

	tracker, err := Get_Tracker(db, id)
	if err != nil {
		log.Fatal(err)
	}

	entries, err := Get_Entries(db, id)
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

	trackers, err := Get_Trackers(db)
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

	tracker, err := Get_Tracker(db, id)
	if err != nil {
		log.Fatal(err)
	}

	entries, err := Get_Entries(db, id)
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

	trackers, err := Get_Trackers(db)
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

	tracker, err := Get_Tracker(db, id)
	if err != nil {
		log.Fatal(err)
	}

	entries, err := Get_Entries(db, id)
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

func htmx_entry_history(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-tracker-create")

	err := r.ParseForm()
	if err != nil {
		return
	}

	tracker_id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return
	}
	r.Form.Del("id")

	search := r.Form.Get("search")
	r.Form.Del("search")

	tracker, err := Get_Tracker(db, tracker_id)
	if err != nil {
		log.Fatal(err)
	}

	entries, err := Get_Entries_Filter(db, tracker_id, search)
	if err != nil {
		log.Fatal(err)
		return
	}

	tmp.ExecuteTemplate(w, "tracker_history", struct {
		Tracker Db_Tracker
		Entries []Db_Entry
	}{
		Tracker: tracker,
		Entries: entries,
	})
}

func htmx_entry_checkbox_toggle(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func page_entry_view(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-entry-view")

	tracker_id, err := strconv.Atoi(r.URL.Query().Get("tracker_id"))
	if err != nil {
		tracker_id = 0
	}

	tracker, err := Get_Tracker(db, tracker_id)
	if err != nil {
		log.Fatal(err)
	}

	entry_id, err := strconv.Atoi(r.URL.Query().Get("entry_id"))
	if err != nil {
		entry_id = 0
	}

	entries, err := Get_Entries(db, tracker_id)
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
		Title   string
		Tracker Db_Tracker
		Entry   Db_Entry
	}{
		Title:   "Entry",
		Tracker: tracker,
		Entry:   entry,
	})
}

func page_entry_editor(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-entry-editor")

	trackers, err := Get_Trackers(db)
	if err != nil {
		log.Fatal(err)
	}

	tracker_id, err := strconv.Atoi(r.URL.Query().Get("tracker_id"))
	if err != nil {
		tracker_id = 0
	}

	tracker, err := Get_Tracker(db, tracker_id)
	if err != nil {
		log.Fatal(err)
	}

	entry_id, err := strconv.Atoi(r.URL.Query().Get("entry_id"))
	if err != nil {
		entry_id = 0
	}

	entries, err := Get_Entries(db, tracker_id)
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

	tmp.ExecuteTemplate(w, "app", struct {
		Title    string
		Trackers []Db_Tracker
		Tracker  Db_Tracker
		Entry    Db_Entry
	}{
		Title:    "Entry",
		Trackers: trackers,
		Tracker:  tracker,
		Entry:    entry,
	})
}

func page_settings(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-settings")

	trackers, err := Get_Trackers(db)
	if err != nil {
		log.Fatal(err)
	}

	tmp.ExecuteTemplate(w, "app", struct {
		Title    string
		Trackers []Db_Tracker
		Tracker  Db_Tracker
	}{
		Title:    "Trackers",
		Tracker:  Db_Tracker{Id: 1, Name: ""},
		Trackers: trackers,
	})
}

func page_content(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-content")

	trackers, err := Get_Trackers(db)
	if err != nil {
		log.Fatal(err)
	}

	content_list, err := Get_Content_List()
	if err != nil {
		log.Fatal(err)
	}

	tmp.ExecuteTemplate(w, "app", struct {
		Title        string
		Tracker      Db_Tracker
		Trackers     []Db_Tracker
		Content_List []Content
	}{
		Title:        "Log",
		Tracker:      Db_Tracker{Id: 1, Name: ""},
		Trackers:     trackers,
		Content_List: content_list,
	})
}

func page_test(w http.ResponseWriter, r *http.Request) {
	tmp := parse_templates("page-test")

	trackers, err := Get_Trackers(db)
	if err != nil {
		log.Fatal(err)
	}

	tracker_id, err := strconv.Atoi(r.URL.Query().Get("tracker_id"))
	if err != nil {
		tracker_id = 1
	}

	tracker, err := Get_Tracker(db, tracker_id)
	if err != nil {
		log.Fatal(err)
	}

	entries, err := Get_Entries(db, tracker_id)
	if err != nil {
		log.Fatal(err)
	}

	// entry_id, err := strconv.Atoi(r.URL.Query().Get("entry_id"))
	// if err != nil {
	// 	entry_id = 1
	// }

	entry := entries[0]

	tmp.ExecuteTemplate(w, "app", struct {
		Title    string
		Tracker  Db_Tracker
		Trackers []Db_Tracker
		Entries  []Db_Entry
		Entry    Db_Entry
	}{
		Title:    "Test",
		Trackers: trackers,
		Tracker:  tracker,
		Entries:  entries,
		Entry:    entry,
	})
}
