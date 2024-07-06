package main

import (
	"database/sql"
	// "fmt"
	"log"
	_ "modernc.org/sqlite"
)

// tables

func Tracker_Test(db *sql.DB) {
	Tracker_New(db, "test-1")
	Tracker_Update_Notes(db, "test-1", "this is just a test")
	Tracker_Add_Number_Field(db, "test-1", "num", false, 500, true, 100, 2)
	Tracker_Add_Option_Field(db, "test-1", "opt", []int{0, 1}, []string{"ok", "error"})

	Record_Table_Create(db, "test-1")
	// Record_Add(db, "test-1", "", []string{"num", "opt"}, []int{42, 0})
	// Record_Add(db, "test-1", "", []string{"num", "opt"}, []int{0, 1})
	// Record_Add(db, "test-1", "tehe", []string{"num", "opt"}, []int{69, 0})
	// Record_Add(db, "test-1", "", []string{"num", "opt"}, []int{3, 1})

	// log.Println("Tracker_Get")
	// resultA, errA := Tracker_Get_All(db)
	// if errA != nil {
	// 	log.Fatal(errA)
	// }
	// log.Println(resultA)

	// log.Println("Tracker_Get_Fields", "test-1")
	// resultB, errB := Tracker_Get_Fields(db, "test-1")
	// if errB != nil {
	// 	log.Fatal(errB)
	// }
	// log.Println(resultB)

	log.Println("Tracker_Get_Fields_Deep", "test-1")
	resultC, errC := Tracker_Get_Fields_Deep(db, "test-1")
	if errC != nil {
		log.Fatal(errC)
	}
	log.Println(resultC)

	// log.Println("Record_Get_Deep", "test-1")
	// resultD, errD := Record_Get_Deep(db, "test-1")
	// if errD != nil {
	// 	log.Fatal(errD)
	// }
	// // log.Println(resultD)
	// log.Printf("tracker '%s' (%d)\n", resultD.tracker.name, resultD.tracker.id)

	// var columes_string string
	// for _, field := range resultD.fields {
	// 	columes_string += fmt.Sprintf("\t%s", field.field_name)
	// }
	// log.Printf("id\ttimestamp\t%s\tnotes\n", columes_string)

	// for _, record := range resultD.records {
	// 	var data_string string
	// 	for i, d := range record.data {
	// 		field := resultD.fields[i]
	// 		if field.field_type == "number" {
	// 			data_string += fmt.Sprintf("\t%d", d)
	// 		} else if field.field_type == "option" {
	// 			data_string += fmt.Sprintf("\t%s", field.type_option.option_names[d])
	// 		}
	// 	}
	// 	log.Printf("%d\t[%s]%s\t'%s'\n", record.id, record.timestamp, data_string, record.notes)
	// }
}
