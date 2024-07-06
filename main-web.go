package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"
)

// https://stackoverflow.com/questions/77124124/how-to-parse-embed-fs-templates-with-the-template-parsefs-function

//go:embed public
var Public_Embed embed.FS

//go:embed templates
var Templates_Embed embed.FS

func Get_Embed_Templates() *template.Template {
	t, err := template.New("").
		ParseFS(Templates_Embed,
			"templates/*.html",
			"templates/*/*.html",
		)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func Start_Web_Server(db *sql.DB) {
	Setup_Public_Routes()
	Setup_Test_Routes()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/trackers", http.StatusSeeOther)
	})

	Page_Tracker(db)
	Page_Names()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	port = fmt.Sprintf(":%s", port)

	hostname, _ := os.Hostname()
	log.Println("Web Server: started")
	log.Printf("- http://%s%s\n", "localhost", port)
	log.Printf("- http://%s%s\n", hostname, port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func Setup_Public_Routes() {
	var public_fs = fs.FS(Public_Embed)
	public_files, err := fs.Sub(public_fs, "public")
	if err != nil {
		log.Fatal(err)
	}
	fs := http.FileServer(http.FS(public_files))
	// http.Handle("/", fs)
	http.Handle("/public/", http.StripPrefix("/public/", fs))
}

func Setup_Test_Routes() {
	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, time.Now().Format(time.UnixDate))
	})
}

func Page_Tracker(db *sql.DB) {
	t, err1 := template.New("").ParseFS(Templates_Embed, "templates/trackers.html")
	if err1 != nil {
		log.Fatal(err1)
	}

	http.HandleFunc("/trackers", func(w http.ResponseWriter, r *http.Request) {
		data := []struct {
			Tracker Tracker
			Fields  []Field
		}{
			// {
			// 	Tracker: Tracker{1, "Walk Dog 🐶", "Each time I take the dog out"},
			// 	Fields:  []Field{},
			// },
			// {
			// 	Tracker: Tracker{2, "Brush Teeth 🦷", "Each time I brush my teeth"},
			// 	Fields:  []Field{},
			// },
			// {
			// 	Tracker: Tracker{3, "My Weight ⚖️", "Weight in pounds"},
			// 	Fields: []Field{
			// 		{Id: 1, Type: "number", Name: "Weight", Notes: ""},
			// 	},
			// },
			// {
			// 	Tracker: Tracker{4, "Lift 🏋️", "My Lifting Habits"},
			// 	Fields: []Field{
			// 		{Id: 2, Type: "option", Name: "Exersise", Notes: ""},
			// 		{Id: 3, Type: "number", Name: "Weight", Notes: ""},
			// 		{Id: 4, Type: "number", Name: "Reps", Notes: ""},
			// 	},
			// },
		}

		trackers, err2 := Tracker_Get_All(db)
		if err2 != nil {
			log.Fatal(err2)
		}

		for _, tracker := range trackers {
			fields, err3 := Tracker_Get_Fields(db, tracker.Name)
			if err3 != nil {
				log.Fatal(err2)
			}

			data = append(data, struct {
				Tracker Tracker
				Fields  []Field
			}{
				Tracker: tracker,
				Fields:  fields,
			})
		}

		t.ExecuteTemplate(w, "trackers.html", data)
	})
}

func Page_Names() {
	t, err := template.New("").ParseFS(Templates_Embed, "templates/names.html")
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Title string
		Names []string
	}{
		Title: "Names",
		Names: []string{"bob", "jim", "arlo"},
	}

	http.HandleFunc("/names", func(w http.ResponseWriter, r *http.Request) {
		t.ExecuteTemplate(w, "names.html", data)
	})

	http.HandleFunc("/names/add", func(w http.ResponseWriter, r *http.Request) {
		name := r.PostFormValue("name-add")
		data.Names = append(data.Names, name)
		t.ExecuteTemplate(w, "block-name", name)
	})
}
