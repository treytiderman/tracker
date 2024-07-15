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

	Page_Trackers(db)
	Page_Tracker_Create(db)
	Page_Tracker_Records(db)
	Page_Tracker_Log(db)
	Page_Tracker_Chart(db)
	Page_Names()

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

func Page_Trackers(db *sql.DB) {
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

	http.HandleFunc("/htmx/tracker/delete", func(w http.ResponseWriter, r *http.Request) {
		id, err1 := strconv.Atoi(r.URL.Query().Get("id"))
		if err1 != nil {
			w.Write([]byte(err1.Error()))
			return
		}

		tracker, err2 := Tracker_By_Id(db, id)
		if err2 != nil {
			w.Write([]byte(err2.Error()))
			return
		}

		// err4 := Record_Table_Delete(db, tracker.Name)
		// if err4 != nil {
		// 	w.Write([]byte(err4.Error()))
		// 	return
		// }

		err3 := Tracker_Delete(db, tracker.Name)
		if err3 != nil {
			w.Write([]byte(err3.Error()))
			return
		}

		w.Write([]byte(fmt.Sprintf("<div class='text-red-300'>Deleted: '%s'</div>", tracker.Name)))
	})
}

func Page_Tracker_Create(db *sql.DB) {
	page, err1 := template.New("").ParseFS(Templates_Embed, "templates/tracker-create.html")
	if err1 != nil {
		log.Fatal(err1)
	}

	http.HandleFunc("/tracker-create", func(w http.ResponseWriter, r *http.Request) {
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
		// w.Write([]byte(fmt.Sprintf("create tracker '%s' with notes '%s'\n", tracker_name, tracker_notes)))

		_, err2 := Tracker_New(db, tracker_name)
		if err2 != nil {
			w.Write([]byte(err2.Error()))
			return
		}

		err3 := Tracker_Update_Notes(db, tracker_name, tracker_notes)
		if err3 != nil {
			w.Write([]byte(err3.Error()))
			return
		}

		for field_id := 0; field_id < 100; field_id++ {
			if r.Form.Has(fmt.Sprintf("field_%d_name", field_id)) {
				field_name := r.Form.Get(fmt.Sprintf("field_%d_name", field_id))
				field_type := r.Form.Get(fmt.Sprintf("field_%d_type", field_id))
				// w.Write([]byte(fmt.Sprintf("  create field [%d] '%s' of type '%s'\n", field_id, field_name, field_type)))

				if field_type == "number" {
					max_flag := false
					max_value := 1000
					min_flag := false
					min_value := 0
					decimal_places := 0

					if r.Form.Has(fmt.Sprintf("field_%d_max_flag", field_id)) {
						max_flag = true
						max_value, _ = strconv.Atoi(r.Form.Get(fmt.Sprintf("field_%d_max_value", field_id)))
						// w.Write([]byte(fmt.Sprintf("    max %d\n", max_value)))
					}
					if r.Form.Has(fmt.Sprintf("field_%d_min_flag", field_id)) {
						min_flag = true
						min_value, _ = strconv.Atoi(r.Form.Get(fmt.Sprintf("field_%d_min_value", field_id)))
						// w.Write([]byte(fmt.Sprintf("    min %d\n", min_value)))
					}
					if r.Form.Has(fmt.Sprintf("field_%d_decimal_places", field_id)) {
						decimal_places, _ = strconv.Atoi(r.Form.Get(fmt.Sprintf("field_%d_decimal_places", field_id)))
						// w.Write([]byte(fmt.Sprintf("    decimal_places %d\n", decimal_places)))
					}

					_, err4 := Tracker_Add_Number_Field(db, tracker_name, field_name, max_flag, max_value, min_flag, min_value, decimal_places)
					if err4 != nil {
						w.Write([]byte(err4.Error()))
						return
					}

				} else if field_type == "option" {
					option_values := []int{}
					option_names := []string{}

					for option_id := 0; option_id < 100; option_id++ {
						if r.Form.Has(fmt.Sprintf("field_%d_option_%d_value", field_id, option_id)) {
							option_name := r.Form.Get(fmt.Sprintf("field_%d_option_%d_name", field_id, option_id))
							option_value, _ := strconv.Atoi(r.Form.Get(fmt.Sprintf("field_%d_option_%d_value", field_id, option_id)))
							option_names = append(option_names, option_name)
							option_values = append(option_values, option_value)
							// w.Write([]byte(fmt.Sprintf("    option [%d] '%s'\n", option_value, option_name)))
						}
					}

					_, err5 := Tracker_Add_Option_Field(db, tracker_name, field_name, option_values, option_names)
					if err5 != nil {
						w.Write([]byte(err5.Error()))
						return
					}
				}
			}
		}

		err6 := Record_Table_Create(db, tracker_name)
		if err6 != nil {
			w.Write([]byte(err6.Error()))
			return
		}

		// w.Write([]byte("\nsuccess"))

		http.Redirect(w, r, "/trackers", http.StatusSeeOther)
	})
}

func Page_Tracker_Records(db *sql.DB) {
	t, err1 := template.New("").ParseFS(Templates_Embed, "templates/tracker-records.html")
	if err1 != nil {
		log.Fatal(err1)
	}

	http.HandleFunc("/tracker/records", func(w http.ResponseWriter, r *http.Request) {
		id, err1 := strconv.Atoi(r.URL.Query().Get("id"))
		if err1 != nil {
			log.Fatal(err1)
			w.Write([]byte(err1.Error()))
			return
		}

		tracker, err2 := Tracker_By_Id(db, id)
		if err2 != nil {
			log.Fatal(err2)
			w.Write([]byte(err2.Error()))
			return
		}

		record_table, err3 := Record_Get_Deep(db, tracker.Name)
		if err3 != nil {
			log.Fatal(err3)
			w.Write([]byte(err3.Error()))
			return
		}

		trackers, err4 := Tracker_Get_All(db)
		if err4 != nil {
			log.Fatal(err4)
		}

		data := struct {
			Trackers []Tracker
			Tracker  Tracker
			Fields   []Field_Deep
			Records  []struct {
				Id        int
				Timestamp string
				Data      []string
				Notes     string
			}
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

		data.Trackers = trackers
		data.Tracker = tracker
		data.Fields = record_table.Fields

		for _, record := range record_table.Records {
			record_to_print := struct {
				Id        int
				Timestamp string
				Data      []string
				Notes     string
			}{
				Id:        int(record.Id),
				Timestamp: record.Timestamp,
				Data:      []string{},
				Notes:     record.Notes,
			}

			for i, data := range record.Data {
				field := record_table.Fields[i]
				if field.Type == "number" {
					data_moved := float32(data) / float32(math.Pow10(field.Type_Number.Decimal_Places))
					data_string := fmt.Sprintf("%.2f", data_moved)
					record_to_print.Data = append(record_to_print.Data, data_string)
				} else if field.Type == "option" {
					for j, val := range field.Type_Option.Option_Values {
						if val == int(data) {
							data_string := fmt.Sprintf("%s", field.Type_Option.Option_Names[j])
							record_to_print.Data = append(record_to_print.Data, data_string)
							break
						}
					}
				}
			}

			data.Records = append(data.Records, record_to_print)
		}

		t.ExecuteTemplate(w, "tracker-records.html", data)
	})
}

func Page_Tracker_Log(db *sql.DB) {
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
			log.Fatal(err1)
			w.Write([]byte(err1.Error()))
			return
		}

		tracker, err2 := Tracker_By_Id(db, id)
		if err2 != nil {
			log.Fatal(err2)
			w.Write([]byte(err2.Error()))
			return
		}

		fields, err3 := Tracker_Get_Fields_Deep(db, tracker.Name)
		if err3 != nil {
			log.Fatal(err3)
			w.Write([]byte(err3.Error()))
			return
		}

		trackers, err4 := Tracker_Get_All(db)
		if err4 != nil {
			log.Fatal(err4)
		}

		data := struct {
			Trackers []Tracker
			Tracker  Tracker
			Fields   []Field_Deep
		}{}

		data.Trackers = trackers
		data.Tracker = tracker
		data.Fields = fields

		t.ExecuteTemplate(w, "tracker-log.html", data)
	})

	http.HandleFunc("/htmx/tracker/record", func(w http.ResponseWriter, r *http.Request) {
		id, err1 := strconv.Atoi(r.URL.Query().Get("id"))
		if err1 != nil {
			w.Write([]byte(err1.Error()))
			return
		}

		tracker, err2 := Tracker_By_Id(db, id)
		if err2 != nil {
			w.Write([]byte(err2.Error()))
			return
		}

		err3 := r.ParseForm()
		if err3 != nil {
			w.Write([]byte(err3.Error()))
			return
		}

		w.Write([]byte(fmt.Sprintf("add record to tracker '%s'\n", tracker.Name)))

		field_names := []string{}
		field_values := []string{}

		for field_id := 0; field_id < 100; field_id++ {
			if r.Form.Has(fmt.Sprintf("field_%d", field_id)) {
				field, err4 := Field_By_Id(db, field_id)
				if err4 != nil {
					w.Write([]byte(err4.Error()))
					return
				}
				
				field_name := field.Name
				field_type := field.Type
				field_value_string := r.Form.Get(fmt.Sprintf("field_%d", field_id))
				
				
				field_number_options, err5 := Tracker_Get_Number(db, field_id)
				if err5 != nil {
					w.Write([]byte(err5.Error()))
					return
				}

				if field_type == "number" {
					field_value_float, _ := strconv.ParseFloat(field_value_string, 64)
					field_value_adjusted := float32(field_value_float) * float32(math.Pow10(field_number_options.Decimal_Places))
					field_value_string = fmt.Sprintf("%.f", field_value_adjusted)
				}

				w.Write([]byte(fmt.Sprintf("  record '%s' of type '%s' as '%s'\n", field_name, field_type, field_value_string)))

				field_names = append(field_names, field_name)
				field_values = append(field_values, field_value_string)
			}
		}

		_, err6 := Record_Add(db, tracker.Name, "", field_names, field_values)
		if err6 != nil {
			w.Write([]byte(err6.Error()))
			return
		}

		w.Write([]byte("\nsuccess"))
	})
}

func Page_Tracker_Chart(db *sql.DB) {
	t, err1 := template.New("").ParseFS(Templates_Embed, "templates/tracker-chart.html")
	if err1 != nil {
		log.Fatal(err1)
	}

	http.HandleFunc("/tracker/chart", func(w http.ResponseWriter, r *http.Request) {
		id, err1 := strconv.Atoi(r.URL.Query().Get("id"))
		if err1 != nil {
			log.Fatal(err1)
			w.Write([]byte(err1.Error()))
			return
		}

		tracker, err2 := Tracker_By_Id(db, id)
		if err2 != nil {
			log.Fatal(err2)
			w.Write([]byte(err2.Error()))
			return
		}

		record_table, err3 := Record_Get_Deep(db, tracker.Name)
		if err3 != nil {
			log.Fatal(err3)
			w.Write([]byte(err3.Error()))
			return
		}

		trackers, err4 := Tracker_Get_All(db)
		if err4 != nil {
			log.Fatal(err4)
		}

		data := struct {
			Trackers []Tracker
			Tracker  Tracker
			Fields   []Field_Deep
			Records  []struct {
				Id        int
				Timestamp string
				Data      []string
				Notes     string
			}
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

		data.Trackers = trackers
		data.Tracker = tracker
		data.Fields = record_table.Fields

		for _, record := range record_table.Records {
			record_to_print := struct {
				Id        int
				Timestamp string
				Data      []string
				Notes     string
			}{
				Id:        int(record.Id),
				Timestamp: record.Timestamp,
				Data:      []string{},
				Notes:     record.Notes,
			}

			for i, data := range record.Data {
				field := record_table.Fields[i]
				if field.Type == "number" {
					data_moved := float32(data) / float32(math.Pow10(field.Type_Number.Decimal_Places))
					data_string := fmt.Sprintf("%.2f", data_moved)
					record_to_print.Data = append(record_to_print.Data, data_string)
				} else if field.Type == "option" {
					for j, val := range field.Type_Option.Option_Values {
						if val == int(data) {
							data_string := fmt.Sprintf("%s", field.Type_Option.Option_Names[j])
							record_to_print.Data = append(record_to_print.Data, data_string)
							break
						}
					}
				}
			}

			data.Records = append(data.Records, record_to_print)
		}

		t.ExecuteTemplate(w, "tracker-chart.html", data)
	})
}
