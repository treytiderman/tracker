package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

// tables

func Tables_Create(db *sql.DB) {
	log.Println("create tables 'tracker', 'field', 'number', and 'option' if they do not exist")

	_, err := db.Exec(`
		PRAGMA foreign_keys = ON;
		CREATE TABLE IF NOT EXISTS tracker (
			tracker_id INTEGER NOT NULL UNIQUE,
			tracker_name TEXT NOT NULL UNIQUE,
			tracker_notes TEXT,
			PRIMARY KEY (tracker_id)
		);
		CREATE TABLE IF NOT EXISTS field (
			field_id INTEGER NOT NULL UNIQUE,
			tracker_id INTEGER NOT NULL,
			field_type TEXT CHECK(field_type in ('number', 'option')) NOT NULL DEFAULT 'number',
			field_name TEXT NOT NULL,
			field_notes TEXT,
			PRIMARY KEY(field_id),
			FOREIGN KEY(tracker_id) REFERENCES tracker (tracker_id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS number (
			number_id INTEGER NOT NULL UNIQUE,
			field_id INTEGER NOT NULL,
			max_flag INTEGER NOT NULL DEFAULT false,
			max_value INTEGER NOT NULL DEFAULT 1000,
			min_flag INTEGER NOT NULL DEFAULT false,
			min_value INTEGER NOT NULL DEFAULT 1,
			decimal_places INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY(number_id)
			FOREIGN KEY(field_id) REFERENCES field (field_id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS option (
			option_id INTEGER NOT NULL UNIQUE,
			field_id INTEGER NOT NULL,
			option_value INTEGER NOT NULL DEFAULT 0,
			option_name TEXT NOT NULL DEFAULT "value",
			PRIMARY KEY(option_id)
			FOREIGN KEY(field_id) REFERENCES field (field_id) ON DELETE CASCADE
		);
	`)

	if err != nil {
		log.Fatal(err)
	}
}

func Tables_Drop(db *sql.DB) {
	log.Println("drop tables 'tracker', 'field', 'number', and 'option' if they exist")

	_, err := db.Exec(`
		DROP TABLE IF EXISTS option;
		DROP TABLE IF EXISTS number;
		DROP TABLE IF EXISTS field;
		DROP TABLE IF EXISTS tracker;
	`)

	if err != nil {
		log.Fatal(err)
	}
}

// trackers

type Tracker struct {
	id    int
	name  string
	notes string
}

func Tracker_Get_All(db *sql.DB) ([]Tracker, error) {
	rows, err1 := db.Query(`SELECT * FROM tracker;`)
	if err1 != nil {
		return nil, err1
	}
	defer rows.Close()

	var trackers []Tracker
	for rows.Next() {
		var tracker Tracker
		err2 := rows.Scan(&tracker.id, &tracker.name, &tracker.notes)
		if err2 != nil {
			return nil, err2
		}
		trackers = append(trackers, tracker)
	}

	return trackers, nil
}

func Tracker_By_Name(db *sql.DB, tracker_name string) (Tracker, error) {
	row := db.QueryRow(`SELECT * FROM tracker WHERE tracker_name = ?;`, tracker_name)

	var tracker Tracker
	err := row.Scan(&tracker.id, &tracker.name, &tracker.notes)
	if err != nil {
		return tracker, err
	}

	return tracker, nil
}

func Tracker_Get_Id(db *sql.DB, tracker_name string) (int, error) {
	row := db.QueryRow(`SELECT tracker_id FROM tracker WHERE tracker_name = ?;`, tracker_name)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func Tracker_New(db *sql.DB, tracker_name string) (int64, error) {
	log.Printf("create tracker '%s'", tracker_name)

	// sql call
	result, err1 := db.Exec(`INSERT INTO tracker (tracker_name) VALUES (?);`, tracker_name)
	if err1 != nil {
		return 0, err1
	}

	// get id of inserted row
	id, err2 := result.LastInsertId()
	if err2 != nil {
		return 0, err2
	}

	return id, nil
}

func Tracker_Delete(db *sql.DB, tracker_name string) error {
	log.Printf("delete tracker '%s'", tracker_name)

	_, err := db.Exec(`DELETE FROM tracker WHERE tracker_name = ?;`, tracker_name)
	if err != nil {
		return err
	}

	return nil
}

func Tracker_Update_Notes(db *sql.DB, tracker_name string, notes string) error {
	log.Printf("update tracker '%s' notes to: %s", tracker_name, notes)

	_, err := db.Exec(`UPDATE tracker SET tracker_notes = ? WHERE tracker_name = ?`, notes, tracker_name)
	if err != nil {
		return err
	}

	return nil
}

// field

type Field struct {
	field_id    int
	field_type  string
	field_name  string
	field_notes string
}

type Field_Deep struct {
	field_id    int
	field_type  string
	field_name  string
	field_notes string
	type_number Field_Number
	type_option Field_Option
}

type Field_Number struct {
	max_flag       bool
	max_value      int
	min_flag       bool
	min_value      int
	decimal_places int
}

type Field_Option struct {
	option_values []int
	option_names  []string
}

func Tracker_Get_Fields(db *sql.DB, tracker_name string) ([]Field, error) {

	// get tracker_id from tracker_name
	tracker_id, err1 := Tracker_Get_Id(db, tracker_name)
	if err1 != nil {
		return nil, err1
	}

	// get all fields of tracker_id
	rows, err1 := db.Query(`SELECT field_id, field_type, field_name, field_notes FROM field WHERE tracker_id = ?;`, tracker_id)
	if err1 != nil {
		return nil, err1
	}
	defer rows.Close()

	var fields []Field
	for rows.Next() {
		var field Field
		err2 := rows.Scan(&field.field_id, &field.field_type, &field.field_name, &field.field_notes)
		if err2 != nil {
			return nil, err2
		}
		fields = append(fields, field)
	}

	return fields, nil
}

func Tracker_Get_Number(db *sql.DB, field_id int) (Field_Number, error) {
	row := db.QueryRow(`SELECT max_flag, max_value, min_flag, min_value, decimal_places 
		FROM number WHERE field_id = ?;`, field_id)

	var number Field_Number
	err := row.Scan(&number.max_flag, &number.max_value, &number.min_flag, &number.min_value, &number.decimal_places)
	if err != nil {
		return number, err
	}

	return number, nil
}

func Tracker_Get_Option(db *sql.DB, field_id int) (Field_Option, error) {
	var option Field_Option

	rows, err1 := db.Query(`SELECT option_value, option_name
		FROM option WHERE field_id = ?;`, field_id)

	if err1 != nil {
		return option, err1
	}
	defer rows.Close()

	for rows.Next() {
		var option_value int
		var option_name string
		err2 := rows.Scan(&option_value, &option_name)
		if err2 != nil {
			return option, err2
		}
		option.option_values = append(option.option_values, option_value)
		option.option_names = append(option.option_names, option_name)
	}

	return option, nil
}

func Tracker_Get_Fields_Deep(db *sql.DB, tracker_name string) ([]Field_Deep, error) {
	var fields_deep []Field_Deep

	fields, err1 := Tracker_Get_Fields(db, tracker_name)
	if err1 != nil {
		return fields_deep, err1
	}

	for _, field := range fields {
		type_number := Field_Number{}
		type_option := Field_Option{}

		// get fields
		if field.field_type == "number" {
			field_number, err2 := Tracker_Get_Number(db, field.field_id)
			if err2 != nil {
				return fields_deep, err1
			}
			type_number = field_number
		} else if field.field_type == "option" {
			field_option, err2 := Tracker_Get_Option(db, field.field_id)
			if err2 != nil {
				return fields_deep, err1
			}
			type_option = field_option
		}

		// build up return object
		fields_deep = append(fields_deep, Field_Deep{
			field_id:    field.field_id,
			field_type:  field.field_type,
			field_name:  field.field_name,
			field_notes: field.field_notes,
			type_number: type_number,
			type_option: type_option,
		})
	}

	return fields_deep, nil
}

func Tracker_Add_Number_Field(db *sql.DB, tracker_name string, field_name string, max_flag bool, max_value int, min_flag bool, min_value int, decimal_places int) (int64, error) {
	log.Printf("add to tracker '%s' field type 'number' named '%s'", tracker_name, field_name)

	// get tracker_id from tracker_name
	tracker_id, err1 := Tracker_Get_Id(db, tracker_name)
	if err1 != nil {
		return 0, err1
	}

	// sql call - field
	result, err2 := db.Exec(`
		INSERT INTO field (tracker_id, field_type, field_name) VALUES (?,"number",?);`,
		tracker_id, field_name)
	if err2 != nil {
		return 0, err2
	}

	// get id of inserted row
	field_id, err3 := result.LastInsertId()
	if err3 != nil {
		return 0, err3
	}

	// sql call - number
	_, err4 := db.Exec(`
		INSERT INTO number (field_id, max_flag, max_value, min_flag, min_value, decimal_places) VALUES (?,?,?,?,?,?);`,
		field_id, max_flag, max_value, min_flag, min_value, decimal_places)
	if err4 != nil {
		return 0, err4
	}

	return field_id, nil
}

func Tracker_Add_Option_Field(db *sql.DB, tracker_name string, field_name string, option_values []int, option_names []string) (int64, error) {
	log.Printf("add to tracker '%s' field type 'option' named '%s'", tracker_name, field_name)

	// get tracker_id from tracker_name
	tracker_id, err1 := Tracker_Get_Id(db, tracker_name)
	if err1 != nil {
		return 0, err1
	}

	// sql call - field
	result, err2 := db.Exec(`
		INSERT INTO field (tracker_id, field_type, field_name) VALUES (?,"option",?);`,
		tracker_id, field_name)
	if err2 != nil {
		return 0, err2
	}

	// get id of inserted row
	field_id, err3 := result.LastInsertId()
	if err3 != nil {
		return 0, err3
	}

	// loop though options
	for i, option_value := range option_values {
		log.Printf("field_id '%d' option_value '%d' option_name '%s'", field_id, option_value, option_names[i])

		// sql call - number
		_, err4 := db.Exec(`INSERT INTO option (field_id, option_value, option_name)
			VALUES (?,?,?);`,
			field_id, option_value, option_names[i])
		if err4 != nil {
			return 0, err4
		}

	}

	return field_id, nil
}

// record

type Record_Table struct {
	tracker Tracker
	records []Record
	fields []Field_Deep
}

type Record struct {
	id        int64
	timestamp string
	notes     string

	data   []int64
}

func Record_Get_Deep(db *sql.DB, tracker_name string) (Record_Table, error) {
	var record_table Record_Table

	tracker, err1 := Tracker_By_Name(db, tracker_name)
	if err1 != nil {
		return record_table, err1
	}
	record_table.tracker = tracker

	fields, err2 := Tracker_Get_Fields_Deep(db, tracker_name)
	if err2 != nil {
		return record_table, err2
	}
	record_table.fields = fields

	var records []Record

	records_query := fmt.Sprintf("SELECT * FROM tracker_%d;", tracker.id)
	rows, err3 := db.Query(records_query)
	if err3 != nil {
		return record_table, err3
	}
	defer rows.Close()

	for rows.Next() {
		cols, err4 := rows.Columns()
		if err4 != nil {
			return record_table, err4
		}

		var record Record

		// make an object based on the number columns
		columns := make([]string, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		rows.Scan(columnPointers...)

		// map scan data to record
		for i, col_name := range cols {
			if col_name == "id" {
				record.id, _ = strconv.ParseInt(columns[i], 10, 0)
			} else if col_name == "tracker_id" {
			} else if col_name == "timestamp" {
				record.timestamp = columns[i]
			} else if col_name == "notes" {
				record.notes = columns[i]
			} else {
				col_int, _ := strconv.ParseInt(columns[i], 10, 0)
				record.data = append(record.data, col_int)
			}
		}
		
		records = append(records, record)
	}
	
	record_table.records = records

	return record_table, nil
}

func Record_Table_Create(db *sql.DB, tracker_name string) error {
	tracker_id, err1 := Tracker_Get_Id(db, tracker_name)
	if err1 != nil {
		return err1
	}

	log.Printf("create table 'tracker_%d' for tracker %s", tracker_id, tracker_name)

	fields, err2 := Tracker_Get_Fields_Deep(db, "test-1")
	if err2 != nil {
		log.Fatal(err2)
	}

	custom_fields_string := "-- custom fields from the field table"
	for _, field := range fields {
		custom_field_string := ""

		if field.field_type == "number" {
			custom_field_string = fmt.Sprintf("%s INT NOT NULL DEFAULT 0,", field.field_name)
		} else if field.field_type == "option" {
			custom_field_string = fmt.Sprintf("%s INT NOT NULL DEFAULT 0,", field.field_name)
		}

		custom_fields_string = strings.Join([]string{custom_fields_string, custom_field_string}, "\n\t")
	}

	create_table_string := fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS tracker_%d (
	id INTEGER NOT NULL UNIQUE,
	tracker_id INTEGER NOT NULL DEFAULT %d,
	timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	notes TEXT NOT NULL DEFAULT "",

	%s

	PRIMARY KEY(id),
	FOREIGN KEY(tracker_id) REFERENCES tracker (tracker_id)
);`, tracker_id, tracker_id, custom_fields_string)

	// sql call
	_, err3 := db.Exec(create_table_string)

	return err3
}

// func Record_Table_Migrate(db *sql.DB) {
// 	// unlock the tracker from adding changing fields
// 	// update fields
// 	// migrate data to new schema
// 	// re-lock the tracker
// }

func Record_Table_Delete(db *sql.DB, tracker_name string) error {
	return nil
}

func Record_Add(db *sql.DB, tracker_name string, notes string, data_names []string, data_values []int) (int64, error) {
	log.Printf("record in '%s'", tracker_name)

	tracker_id, err1 := Tracker_Get_Id(db, tracker_name)
	if err1 != nil {
		return 0, err1
	}

	field_names_string := strings.Join(data_names, ", ")

	var field_values []string
	for _, data_value := range data_values {
		field_values = append(field_values, strconv.Itoa(data_value))
	}

	field_values_string := strings.Join(field_values, ", ")

	insert_string := fmt.Sprintf(
		`INSERT INTO tracker_%d (notes, %s) VALUES ("%s", %s);`,
		tracker_id, field_names_string, notes, field_values_string)

	// sql call
	result, err2 := db.Exec(insert_string)
	if err2 != nil {
		return 0, err2
	}

	// get id of inserted row
	record_id, err3 := result.LastInsertId()
	if err3 != nil {
		return 0, err3
	}

	return record_id, nil
}
