package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func handle_routes_api_htmx(mux *http.ServeMux) {
	mux.Handle("GET /htmx/token", http.HandlerFunc(auth_token_get))   // test cookie
	mux.Handle("POST /htmx/token", http.HandlerFunc(auth_token_post)) // set cookie if password correct

	mux.Handle("/htmx/test/form", mw_logger(http.HandlerFunc(htmx_test_form)))

	mux.Handle("POST /htmx/tracker/create", mw_logger(mw_read_only(mw_auth(http.HandlerFunc(htmx_tracker_create)))))
	mux.Handle("POST /htmx/tracker/name", mw_logger(mw_read_only(mw_auth(http.HandlerFunc(htmx_tracker_name)))))
	mux.Handle("POST /htmx/tracker/notes", mw_logger(mw_read_only(mw_auth(http.HandlerFunc(htmx_tracker_notes)))))
	mux.Handle("GET /htmx/tracker/delete", mw_logger(mw_read_only(mw_auth(http.HandlerFunc(htmx_tracker_delete)))))

	mux.Handle("POST /htmx/entry/create", mw_logger(mw_read_only(mw_auth(http.HandlerFunc(htmx_entry_create)))))
	mux.Handle("POST /htmx/entry/update", mw_logger(mw_read_only(mw_auth(http.HandlerFunc(htmx_entry_update)))))

	mux.Handle("POST /htmx/tracker/log", mw_logger(mw_read_only(mw_auth(http.HandlerFunc(htmx_log_create)))))
	mux.Handle("POST /htmx/tracker/log-update", mw_logger(mw_read_only(mw_auth(http.HandlerFunc(htmx_log_update)))))
	mux.Handle("GET /htmx/tracker/log-delete", mw_logger(mw_read_only(mw_auth(http.HandlerFunc(htmx_log_delete)))))

	mux.Handle("GET /content/{content_path}", mw_auth(http.HandlerFunc(content_get)))
	mux.Handle("POST /content-upload", mw_logger(mw_read_only(mw_auth(http.HandlerFunc(content_upload)))))
	mux.Handle("DELETE /content/{content_path}", mw_logger(mw_read_only(mw_auth(http.HandlerFunc(content_delete)))))
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
			Name:     "token",
			Value:    form_password,
			Path:     "/",
			MaxAge:   0,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &cookie)
		slog.Info("password correct, cookie set")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		slog.Warn("password not correct")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

// Test

func htmx_test_form(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

// Tracker

func htmx_tracker_create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}

	tracker_name := r.Form.Get("tracker_name")
	tracker_notes := r.Form.Get("tracker_notes")

	tracker_id, err := Create_Tracker(db, tracker_name, tracker_notes)
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

				_, err := Add_Number_Field(db, tracker_id, field_name, field_notes, decimal_places)
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

				_, err := Add_Option_Field_With_Options(db, tracker_id, field_name, field_notes, options)
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

	// Get Id from URL
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return
	}

	// Get Tracker Name from Form Data
	tracker_name := r.Form.Get("tracker_name")

	// Update Tracker Name
	err = Update_Tracker_Name(db, id, tracker_name)
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

	// Get Id from URL
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		return
	}

	// Get Tracker Notes from Form Data
	tracker_notes := r.Form.Get("tracker_notes")

	// Update Tracker Notes
	err = Update_Tracker_Notes(db, id, tracker_notes)
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
	err = Delete_Tracker(db, id)
	if err != nil {
		return
	}

	// Reload without Id
	http.Redirect(w, r, "/tracker-info", http.StatusSeeOther)
}

// Entry

func htmx_entry_create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}

	// Get Id from URL
	id, err := strconv.Atoi(r.URL.Query().Get("tracker_id"))
	if err != nil {
		return
	}
	r.Form.Del("tracker_id")

	// Get then Delete non field stuff from form
	entry_notes := r.Form.Get("entry_notes")
	r.Form.Del("entry_notes")

	entry_date := r.Form.Get("entry_date")
	r.Form.Del("entry_date")

	entry_time := r.Form.Get("entry_time")
	r.Form.Del("entry_time")
	if len([]rune(entry_time)) == 5 {
		entry_time = entry_time + ":00"
	}

	entry_timezone := r.Form.Get("entry_timezone")
	r.Form.Del("entry_timezone")

	time_string := fmt.Sprintf("%s %s %s", entry_date, entry_time, entry_timezone)
	dt, err := time.Parse("2006-01-02 15:04:05 -0700", time_string)
	if err != nil {
		return
	}

	timestamp := dt.UTC().Format("2006-01-02 15:04:05")

	// Get Tracker by Id
	tracker, err := Get_Tracker(db, id)
	if err != nil {
		log.Fatal(err)
	}

	entry_id, err := Create_Entry(db, tracker.Id, entry_notes)
	if err != nil {
		return
	}

	err = Update_Entry_Timestamp(db, entry_id, timestamp)
	if err != nil {
		return
	}

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

		Add_Log_To_Entry(db, entry_id, field_id, value)
	}

	// Reload page
	url := fmt.Sprintf("/tracker-history?id=%d", tracker.Id)
	w.Header().Add("HX-Redirect", url)
	w.Write([]byte("ok"))
}

func htmx_entry_update(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}

	// Get Id from URL
	entry_id, err := strconv.Atoi(r.URL.Query().Get("entry_id"))
	if err != nil {
		return
	}
	r.Form.Del("entry_id")

	// Get Id from URL
	tracker_id, err := strconv.Atoi(r.URL.Query().Get("tracker_id"))
	if err != nil {
		return
	}
	r.Form.Del("tracker_id")

	tracker, err := Get_Tracker(db, tracker_id)
	if err != nil {
		log.Fatal(err)
	}

	// Get then Delete non field stuff from form
	entry_notes := r.Form.Get("entry_notes")
	r.Form.Del("entry_notes")

	entry_date := r.Form.Get("entry_date")
	r.Form.Del("entry_date")

	entry_time := r.Form.Get("entry_time")
	r.Form.Del("entry_time")
	if len([]rune(entry_time)) == 5 {
		entry_time = entry_time + ":00"
	}

	entry_timezone := r.Form.Get("entry_timezone")
	r.Form.Del("entry_timezone")

	time_string := fmt.Sprintf("%s %s %s", entry_date, entry_time, entry_timezone)
	dt, err := time.Parse("2006-01-02 15:04:05 -0700", time_string)
	if err != nil {
		return
	}

	timestamp := dt.UTC().Format("2006-01-02 15:04:05")

	for k, v := range r.Form {

		ids := strings.Split(k, "__")

		log_id, err := strconv.Atoi(strings.ReplaceAll(ids[0], "log_", ""))
		if err != nil {
			return
		}

		field_id, err := strconv.Atoi(strings.ReplaceAll(ids[1], "field_", ""))
		if err != nil {
			return
		}

		var field Db_Field
		for _, f := range tracker.Fields {
			if f.Id == field_id {
				field = f
			}
		}

		log_value := 0
		if field.Type == "number" {
			field_value_float, _ := strconv.ParseFloat(v[0], 64)
			field_value_adjusted := float64(field_value_float) * float64(math.Pow10(field.Number.Decimal_Places))
			log_value = int(math.Floor(field_value_adjusted))
			Update_Log(db, log_id, log_value)
		} else if field.Type == "option" {
			log_value, err = strconv.Atoi(v[0])
			if err != nil {
				return
			}
			Update_Log(db, log_id, log_value)
		}
	}

	Update_Entry_Notes(db, entry_id, entry_notes)
	Update_Entry_Timestamp(db, entry_id, timestamp)

	// Reload page
	url := fmt.Sprintf("/tracker-history?id=%d", tracker.Id)
	w.Header().Add("HX-Redirect", url)
	w.Write([]byte("ok"))
}

func htmx_log_create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}

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
	tracker, err := Get_Tracker(db, id)
	if err != nil {
		log.Fatal(err)
	}

	entry_id, err := Create_Entry(db, tracker.Id, entry_notes)
	if err != nil {
		return
	}

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

		Add_Log_To_Entry(db, entry_id, field_id, value)
	}

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

	entry_note := r.Form.Get("entry_note")

	err = Update_Entry_Notes(db, entry_id, entry_note)
	if err != nil {
		log.Fatalln(err)
	}

	r.Form.Del("tracker_id")
	r.Form.Del("entry_id")
	r.Form.Del("entry_note")

	// Get Tracker by Id
	entries, err := Get_Entries(db, tracker_id)
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

		err = Update_Log(db, log_id, log_value)
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
	err = Delete_Entry(db, entry_id)
	if err != nil {
		log.Fatalln(err)
	}

	// Reload page
	url := fmt.Sprintf("/tracker-history?id=%d", tracker_id)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// File

func content_get(w http.ResponseWriter, r *http.Request) {
	content_path := "../content/" + r.PathValue("content_path")

	file, err := os.ReadFile(content_path)
	if err != nil {
		w.Write([]byte("file not found\n"))
		w.Write([]byte(content_path))
	}

	w.Write(file)
}

func content_delete(w http.ResponseWriter, r *http.Request) {
	content_path := "../content/" + r.PathValue("content_path")

	err := os.Remove(content_path)
	if err != nil {
		w.Write([]byte("file not found"))
	}

	// remove content path from db

	slog.Info("FILE DELETED:", "path", content_path)
	w.Write([]byte("ok"))
}

func content_upload(w http.ResponseWriter, r *http.Request) {
	body_bytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	r.Body.Close()                                     // must close
	r.Body = io.NopCloser(bytes.NewBuffer(body_bytes)) // recreate the reader

	mime_type := http.DetectContentType(body_bytes)

	ext := ".png"
	if mime_type == "image/png" {
		ext = ".png"
	} else if mime_type == "image/jpeg" {
		ext = ".jpg"
	} else if mime_type == "audio/mpeg" || mime_type == "application/octet-stream" {
		ext = ".mp3"
	} else if mime_type == "video/mp4" {
		ext = ".mp4"
	} else if mime_type == "text/plain; charset=utf-8" {
		ext = ".txt" // also ".svg" or ".csv"
	} else if mime_type == "application/zip" {
		ext = ".zip"
	} else if mime_type == "application/pdf" {
		ext = ".pdf"
	}

	timestamp := time.Now().Format(time.DateTime)
	timestamp = strings.ReplaceAll(timestamp, " ", "_")
	timestamp = strings.ReplaceAll(timestamp, ":", "-")
	path := "../content/" + timestamp + ext

	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("FILE Saved:", "path", path)
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		log.Fatal(err)
	}

	// add content path to db

	path = strings.ReplaceAll(path, "../data/", "/")
	w.Write([]byte(path))
}
