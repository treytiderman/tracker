package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func handle_routes_api_htmx(mux *http.ServeMux) {
	mux.Handle("GET /htmx/token", mw_logger(http.HandlerFunc(auth_token_get)))   // test cookie
	mux.Handle("POST /htmx/token", mw_logger(http.HandlerFunc(auth_token_post))) // set cookie if password correct

	mux.Handle("POST /htmx/tracker/create", mw_logger(mw_auth(http.HandlerFunc(htmx_tracker_create))))
	mux.Handle("POST /htmx/tracker/name", mw_logger(mw_auth(http.HandlerFunc(htmx_tracker_name))))
	mux.Handle("POST /htmx/tracker/notes", mw_logger(mw_auth(http.HandlerFunc(htmx_tracker_notes))))
	mux.Handle("POST /htmx/tracker/delete", mw_logger(mw_auth(http.HandlerFunc(htmx_tracker_delete))))

	mux.Handle("POST /htmx/tracker/log", mw_logger(mw_auth(http.HandlerFunc(htmx_log_create))))
	mux.Handle("POST /htmx/tracker/log-update", mw_logger(mw_auth(http.HandlerFunc(htmx_log_update))))
	mux.Handle("POST /htmx/tracker/log-delete", mw_logger(mw_auth(http.HandlerFunc(htmx_log_delete))))

	mux.Handle("POST /ui/upload", mw_logger(mw_auth(http.HandlerFunc(htmx_file_upload))))
}

// Auth

func auth_token_valid(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return false
	}

	password := os.Getenv("PASSWORD")
	if password == "" {
		password = "password"
	}

	if cookie.Value == password {
		return true
	} else {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return false
	}
}

func auth_token_get(w http.ResponseWriter, r *http.Request) {
	token_is_valid := auth_token_valid(w, r)
	if !token_is_valid {
		return
	}

	w.Write([]byte("cookie valid"))
}

func auth_token_post(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}

	form_password := r.Form.Get("password")

	password := os.Getenv("PASSWORD")
	if password == "" {
		password = "password"
	}

	if form_password == password {
		cookie := http.Cookie{
			Name:   "token",
			Value:  form_password,
			Path:   "/",
			MaxAge: 0,
			// HttpOnly: true, // can't access via javascript
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &cookie)
		log.Printf("FORM: password correct\n")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		log.Printf("FORM: password not correct\n")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// Tracker

func htmx_tracker_create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}
	log.Printf("FORM: %s\n", r.Form.Encode())

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
}

func htmx_tracker_name(w http.ResponseWriter, r *http.Request) {
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
}

func htmx_tracker_notes(w http.ResponseWriter, r *http.Request) {
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
}

func htmx_tracker_delete(w http.ResponseWriter, r *http.Request) {
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
}

// Log

func htmx_log_create(w http.ResponseWriter, r *http.Request) {
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
}

func htmx_log_update(w http.ResponseWriter, r *http.Request) {
	tracker_id, err := strconv.Atoi(r.URL.Query().Get("tracker_id"))
	if err != nil {
		log.Fatalln(err)
	}

	entry_id, err := strconv.Atoi(r.URL.Query().Get("entry_id"))
	if err != nil {
		log.Fatalln(err)
	}

	err = r.ParseForm()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("FORM: %s\n", r.Form.Encode())

	entry_note := r.Form.Get("entry_note")

	err = Db_Entry_Notes_Update(db, entry_id, entry_note)
	if err != nil {
		log.Fatalln(err)
	}

	r.Form.Del("tracker_id")
	r.Form.Del("entry_id")
	r.Form.Del("entry_note")

	// Get Tracker by Id
	entries, err := Db_Entry_Get(db, tracker_id)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range r.Form {

		log_id, err := strconv.Atoi(strings.ReplaceAll(k, "log_", ""))
		if err != nil {
			log.Fatalln(err)
		}

		var ll Db_Log
		for _, entry := range entries {
			for _, l := range entry.Logs {
				if l.Id == log_id {
					ll = l
				}
			}
		}

		log_value := 0
		if ll.Field_Type == "number" {
			log_value_float, _ := strconv.ParseFloat(v[0], 64)
			log_value_adjusted := float64(log_value_float) * float64(math.Pow10(ll.Decimal_Places))
			log_value = int(math.Floor(log_value_adjusted))
		} else if ll.Field_Type == "option" {
			log_value, err = strconv.Atoi(v[0])
			if err != nil {
				return
			}
		}

		err = Db_Entry_Log_Update(db, log_id, log_value)
		if err != nil {
			log.Fatalln(err)
		}

	}

	// Reload page
	url := fmt.Sprintf("/tracker-history?id=%d", tracker_id)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func htmx_log_delete(w http.ResponseWriter, r *http.Request) {
	tracker_id, err := strconv.Atoi(r.URL.Query().Get("tracker_id"))
	if err != nil {
		log.Fatalln(err)
	}

	entry_id, err := strconv.Atoi(r.URL.Query().Get("entry_id"))
	if err != nil {
		log.Fatalln(err)
	}

	// Delete Entry
	err = Db_Entry_Delete(db, entry_id)
	if err != nil {
		log.Fatalln(err)
	}

	// Reload page
	url := fmt.Sprintf("/tracker-history?id=%d", tracker_id)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// File

func htmx_file_upload(w http.ResponseWriter, r *http.Request) {
	timestamp := time.Now().Format(time.DateTime)
	timestamp = strings.ReplaceAll(timestamp, " ", "_")
	path := "../public/upload/" + timestamp + ".png"

	img, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("FILE SAVED: %s\n", path)
	defer img.Close()

	_, err = io.Copy(img, r.Body)
	if err != nil {
		log.Fatal(err)
	}

	path = strings.ReplaceAll(path, "../", "/")
	w.Write([]byte(path))
}
