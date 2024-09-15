package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Db_Tracker struct {
	Id     int
	Name   string
	Notes  string
	Fields []Db_Field
}

type Db_Field struct {
	Id      int
	Type    string
	Name    string
	Notes   string
	Number  Db_Number
	Options []Db_Option
}

type Db_Number struct {
	Id             int
	Decimal_Places int
}

type Db_Option struct {
	Id    int
	Value int
	Name  string
}

type Db_Entry struct {
	Id        int
	Timestamp string
	Notes     string
	Logs      []Db_Log
}

type Db_Log struct {
	Id             int
	Value          int
	Field_Id       int
	Field_Type     string
	Field_Name     string
	Decimal_Places int
	Option_Value   int
	Option_Name    string
	Present        string
}

// SQL Tables

func Create_Tables(db *sql.DB) error {
	// fmt.Println("Create SQL tables if they do not exist")

	_, err := db.Exec(`
		-- foreign_keys constraints are not on by default
		PRAGMA foreign_keys = ON;
		
		CREATE TABLE IF NOT EXISTS tracker (
			tracker_id INTEGER NOT NULL UNIQUE,
		
			-- name to identify this tracker
			tracker_name TEXT NOT NULL UNIQUE,
		
			-- markdown formatted notes
			tracker_notes TEXT NOT NULL DEFAULT "",
		
			PRIMARY KEY (tracker_id)
		);
		
		CREATE TABLE IF NOT EXISTS field (
			field_id INTEGER NOT NULL UNIQUE,
		
			-- parent tracker
			tracker_id INTEGER NOT NULL,
		
			-- use "number" to track a signed whole number
			-- examples: weight, height...
			-- use "option" to a list of options
			-- examples: exercise, read status
			field_type TEXT CHECK(field_type in ('number', 'option')) NOT NULL DEFAULT 'number',
		
			-- name to identify this field
			field_name TEXT NOT NULL,
		
			-- markdown formatted notes
			field_notes TEXT NOT NULL DEFAULT "",
		
			-- a tracker can not have duplicate field_name's
			-- but multiple trackers can have the same field_name
			UNIQUE(tracker_id, field_name),
		
			PRIMARY KEY(field_id),
			FOREIGN KEY(tracker_id) REFERENCES tracker (tracker_id) ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS number (
			number_id INTEGER NOT NULL UNIQUE,
		
			-- parent field
			field_id INTEGER NOT NULL,
		
			-- 0 = round to integer
			-- 2 = round to 2 decimal places. example money
			-- -3 = round to thousands
			decimal_places INTEGER NOT NULL DEFAULT 0,
		
			PRIMARY KEY(number_id),
			FOREIGN KEY(field_id) REFERENCES field (field_id) ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS option (
			option_id INTEGER NOT NULL UNIQUE,
		
			-- parent field
			field_id INTEGER NOT NULL,
		
			-- key value pair
			option_value INTEGER NOT NULL DEFAULT 0,
			option_name TEXT NOT NULL DEFAULT "value",
		
			-- an option can not have duplicate option_name's
			UNIQUE(field_id, option_name),
		
			PRIMARY KEY(option_id),
			FOREIGN KEY(field_id) REFERENCES field (field_id) ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS entry (
			entry_id INTEGER NOT NULL UNIQUE,
		
			-- parent tracker
			tracker_id INTEGER NOT NULL,
		
			-- when this was inserted
			timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		
			-- markdown formatted notes
			entry_notes TEXT NOT NULL DEFAULT "",
		
			PRIMARY KEY(entry_id),
			FOREIGN KEY(tracker_id) REFERENCES tracker (tracker_id) ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS log (
			log_id INTEGER NOT NULL UNIQUE,
		
			-- parent field
			entry_id INTEGER NOT NULL,
		
			-- field info
			field_id INTEGER NOT NULL,
		
			-- value
			log_value INTEGER NOT NULL,
		
			PRIMARY KEY(log_id),
			FOREIGN KEY(entry_id) REFERENCES entry (entry_id) ON DELETE CASCADE,
			FOREIGN KEY(field_id) REFERENCES field (field_id) ON DELETE CASCADE
		);
	`)

	return err
}

func Reset_Tables(db *sql.DB) error {
	// fmt.Println("Drop SQL tables if they exist")

	_, err := db.Exec(`
		DROP TABLE IF EXISTS log;
		DROP TABLE IF EXISTS entry;
		DROP TABLE IF EXISTS option;
		DROP TABLE IF EXISTS number;
		DROP TABLE IF EXISTS field;
		DROP TABLE IF EXISTS tracker;
	`)

	if err != nil {
		return err
	}

	err = Create_Tables(db)

	return err
}

// Get Data

func Get_Tracker_Id_By_Name(db *sql.DB, tracker_name string) (int, error) {
	sql_string := fmt.Sprintf(
		`SELECT tracker_id FROM tracker WHERE tracker_name = "%s";`,
		tracker_name)

	// fmt.Println("SQL:", sql_string)
	row := db.QueryRow(sql_string)

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func Get_Tracker_By_Id(db *sql.DB, tracker_id int) (Db_Tracker, error) {
	sql_string := fmt.Sprintf(
		`SELECT
			tracker_id, tracker_name, tracker_notes,
		
			IFNULL(field_id, 0) field_id,
			IFNULL(field_type, "") field_type,
			IFNULL(field_name, "") field_name,
			IFNULL(field_notes, "") field_notes,
		
			IFNULL(number_id, 0) number_id,
			IFNULL(decimal_places, 0) decimal_places,
		
			IFNULL(option_id, 0) option_id,
			IFNULL(option_value, 0) option_value,
			IFNULL(option_name, "") option_name
		FROM (
			SELECT * FROM tracker
			LEFT JOIN field USING (tracker_id)
			WHERE tracker_id = %d
		) AS tf
		LEFT JOIN number AS n USING (field_id)
		LEFT JOIN option AS o USING (field_id)
		ORDER BY tf.field_id, n.number_id, o.option_id;`,
		tracker_id)

	rows, err := db.Query(sql_string)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var tracker Db_Tracker
	var field_scan_last_id = 0

	for rows.Next() {
		var tracker_scan Db_Tracker
		var field_scan Db_Field
		var number_scan Db_Number
		var option_scan Db_Option
		err := rows.Scan(
			&tracker_scan.Id, &tracker_scan.Name, &tracker_scan.Notes,
			&field_scan.Id, &field_scan.Type, &field_scan.Name, &field_scan.Notes,
			&number_scan.Id, &number_scan.Decimal_Places,
			&option_scan.Id, &option_scan.Value, &option_scan.Name,
		)
		if err != nil {
			log.Fatal(err)
		}

		if tracker.Id == 0 {
			tracker = tracker_scan
		}

		if field_scan.Type == "number" {
			field_scan.Number = number_scan
			tracker.Fields = append(tracker.Fields, field_scan)
		} else if field_scan.Type == "option" {

			// Field already added
			if field_scan_last_id == field_scan.Id {
				tracker.Fields[len(tracker.Fields)-1].Options =
					append(tracker.Fields[len(tracker.Fields)-1].Options, option_scan)
			} else {
				field_scan.Options = append(field_scan.Options, option_scan)
				tracker.Fields = append(tracker.Fields, field_scan)
			}
		}

		field_scan_last_id = field_scan.Id
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return tracker, nil
}

func Get_Tracker_By_Name(db *sql.DB, tracker_name string) (Db_Tracker, error) {
	tracker_id, err1 := Get_Tracker_Id_By_Name(db, tracker_name)
	if err1 != nil {
		return Db_Tracker{}, err1
	}

	tracker, err2 := Get_Tracker_By_Id(db, tracker_id)
	if err2 != nil {
		return Db_Tracker{}, err2
	}

	return tracker, nil
}

func Get_Trackers(db *sql.DB) ([]Db_Tracker, error) {
	sql_string :=
		`SELECT
			tracker_id, tracker_name, tracker_notes,
		
			IFNULL(field_id, 0) field_id,
			IFNULL(field_type, "") field_type,
			IFNULL(field_name, "") field_name,
			IFNULL(field_notes, "") field_notes,
		
			IFNULL(number_id, 0) number_id,
			IFNULL(decimal_places, 0) decimal_places,
		
			IFNULL(option_id, 0) option_id,
			IFNULL(option_value, 0) option_value,
			IFNULL(option_name, "") option_name
		FROM (
			SELECT * FROM tracker
			LEFT JOIN field USING (tracker_id)
		) AS tf
		LEFT JOIN number AS n USING (field_id)
		LEFT JOIN option AS o USING (field_id)
		ORDER BY tf.tracker_id, tf.field_id, n.number_id, o.option_id;`

	rows, err := db.Query(sql_string)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var trackers []Db_Tracker
	var tracker_scan_last_id = 0
	var field_scan_last_id = 0

	for rows.Next() {
		var tracker_scan Db_Tracker
		var field_scan Db_Field
		var number_scan Db_Number
		var option_scan Db_Option
		err := rows.Scan(
			&tracker_scan.Id, &tracker_scan.Name, &tracker_scan.Notes,
			&field_scan.Id, &field_scan.Type, &field_scan.Name, &field_scan.Notes,
			&number_scan.Id, &number_scan.Decimal_Places,
			&option_scan.Id, &option_scan.Value, &option_scan.Name,
		)
		if err != nil {
			log.Fatal(err)
		}

		// New
		if tracker_scan_last_id != tracker_scan.Id {
			trackers = append(trackers, tracker_scan)
		}
		if field_scan.Id > 0 {
			if field_scan_last_id != field_scan.Id {
				trackers[len(trackers)-1].Fields = append(trackers[len(trackers)-1].Fields, field_scan)
				trackers[len(trackers)-1].Fields[len(trackers[len(trackers)-1].Fields)-1].Number = number_scan
				trackers[len(trackers)-1].Fields[len(trackers[len(trackers)-1].Fields)-1].Options =
					append(trackers[len(trackers)-1].Fields[len(trackers[len(trackers)-1].Fields)-1].Options, option_scan)
			} else {
				trackers[len(trackers)-1].Fields[len(trackers[len(trackers)-1].Fields)-1].Options =
					append(trackers[len(trackers)-1].Fields[len(trackers[len(trackers)-1].Fields)-1].Options, option_scan)
			}
		}

		tracker_scan_last_id = tracker_scan.Id
		field_scan_last_id = field_scan.Id
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return trackers, nil
}

func Get_Entries_By_Tracker_Id(db *sql.DB, tracker_id int) ([]Db_Entry, error) {
	sql_string := fmt.Sprintf(
		`SELECT
			entry.entry_id,
			entry.timestamp,
			entry.entry_notes,

			IFNULL(log.log_id, 0) AS log_id,
			IFNULL(log.log_value, 0) AS log_value,

			IFNULL(field.field_id, 0) AS field_id,
			IFNULL(field.field_type, "") AS field_type,
			IFNULL(field.field_name, "") AS field_name,

			IFNULL(number.decimal_places, 0) AS decimal_places,

			IFNULL(option.option_value, 0) AS option_value,
			IFNULL(option.option_name, "") AS option_name,

			IFNULL((CASE WHEN field.field_type == "number" THEN
				printf("%%.2f", log.log_value / power(10, number.decimal_places))
			ELSE
				option.option_name
			END), "") AS present
		FROM entry
		LEFT JOIN log USING (entry_id)
		LEFT JOIN field USING (field_id)
		LEFT JOIN number USING (field_id)
		LEFT JOIN option ON log.field_id = option.field_id AND log.log_value = option.option_value
		WHERE entry.tracker_id = %d
		ORDER BY entry.entry_id, field.field_id;`,
		tracker_id)

	rows, err := db.Query(sql_string)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var entries []Db_Entry
	var log_scan_last_id = 0
	var entry_scan_last_id = 0

	for rows.Next() {
		var entry_scan Db_Entry
		var log_scan Db_Log
		err := rows.Scan(
			&entry_scan.Id, &entry_scan.Timestamp, &entry_scan.Notes,
			&log_scan.Id, &log_scan.Value, &log_scan.Field_Id, &log_scan.Field_Type, &log_scan.Field_Name,
			&log_scan.Decimal_Places, &log_scan.Option_Value, &log_scan.Option_Name, &log_scan.Present,
		)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Println("entry_scan", entry_scan)
		// fmt.Println("log_scan", log_scan)

		if entry_scan_last_id != entry_scan.Id {
			entries = append(entries, entry_scan)
		}

		// Why am I getting duplicate log_scan ids?
		// This check for log_scan_last_id should not be needed
		if log_scan_last_id != log_scan.Id {

			if log_scan.Id > 0 {
				entries[len(entries)-1].Logs = append(entries[len(entries)-1].Logs, log_scan)
			}
		}

		entry_scan_last_id = entry_scan.Id
		log_scan_last_id = log_scan.Id
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return entries, nil
}

// Insert Data

func Create_Tracker(db *sql.DB, tracker_name string, tracker_notes string) (int, error) {
	sql_string := fmt.Sprintf(
		`INSERT INTO tracker (tracker_name, tracker_notes) VALUES ("%s", "%s");`,
		tracker_name, tracker_notes)

	fmt.Println("SQL:", sql_string)
	result, err1 := db.Exec(sql_string)
	if err1 != nil {
		return 0, err1
	}

	id, err2 := result.LastInsertId()
	if err2 != nil {
		return 0, err2
	}

	return int(id), nil
}

func Add_Number_Field(db *sql.DB, tracker_name string, field_name string, field_notes string, decimal_places int) (int, error) {
	tracker_id, err1 := Get_Tracker_Id_By_Name(db, tracker_name)
	if err1 != nil {
		return 0, err1
	}

	// Insert Field Row
	sql_string := fmt.Sprintf(
		`INSERT INTO field (tracker_id, field_type, field_name, field_notes) VALUES (%d,"number","%s","%s");`,
		tracker_id, field_name, field_notes)

	fmt.Println("SQL:", sql_string)
	result, err2 := db.Exec(sql_string)
	if err2 != nil {
		return 0, err2
	}

	field_id, err3 := result.LastInsertId()
	if err3 != nil {
		return 0, err3
	}

	// Insert Number Row
	sql_string2 := fmt.Sprintf(
		`INSERT INTO number (field_id, decimal_places) VALUES (%d,%d);`,
		field_id, decimal_places)

	fmt.Println("SQL:", sql_string2)
	_, err4 := db.Exec(sql_string2)
	if err4 != nil {
		return 0, err4
	}

	return int(field_id), nil
}

func Add_Option_Field(db *sql.DB, tracker_name string, field_name string, field_notes string, options []struct {
	Value int
	Name  string
}) (int, error) {
	tracker_id, err1 := Get_Tracker_Id_By_Name(db, tracker_name)
	if err1 != nil {
		return 0, err1
	}

	// Insert Field Row
	sql_string := fmt.Sprintf(
		`INSERT INTO field (tracker_id, field_type, field_name, field_notes) VALUES (%d,"option","%s","%s");`,
		tracker_id, field_name, field_notes)

	fmt.Println("SQL:", sql_string)
	result, err2 := db.Exec(sql_string)
	if err2 != nil {
		return 0, err2
	}

	field_id, err3 := result.LastInsertId()
	if err3 != nil {
		return 0, err3
	}

	// Insert Option Rows
	for _, option := range options {
		sql_string2 := fmt.Sprintf(
			`INSERT INTO option (field_id, option_value, option_name) VALUES (%d,%d,"%s");`,
			field_id, option.Value, option.Name)

		fmt.Println("SQL:", sql_string2)
		_, err4 := db.Exec(sql_string2)
		if err4 != nil {
			return 0, err4
		}
	}

	return int(field_id), nil
}

func Add_Entry(db *sql.DB, tracker_name string, entry_notes string, logs []struct {
	Field_Id int
	Value    int
}) (int, error) {
	tracker_id, err1 := Get_Tracker_Id_By_Name(db, tracker_name)
	if err1 != nil {
		return 0, err1
	}

	// Insert Entry Row
	sql_string := fmt.Sprintf(
		`INSERT INTO entry (tracker_id, entry_notes) VALUES (%d,"%s");`,
		tracker_id, entry_notes)

	fmt.Println("SQL:", sql_string)
	result, err2 := db.Exec(sql_string)
	if err2 != nil {
		return 0, err2
	}

	entry_id, err3 := result.LastInsertId()
	if err3 != nil {
		return 0, err3
	}

	// Insert Option Rows
	for _, log := range logs {
		sql_string2 := fmt.Sprintf(
			`INSERT INTO log (entry_id, field_id, log_value) VALUES (%d,%d,%d);`,
			entry_id, log.Field_Id, log.Value)

		fmt.Println("SQL:", sql_string2)
		_, err4 := db.Exec(sql_string2)
		if err4 != nil {
			return 0, err4
		}
	}

	return int(entry_id), nil
}

// Update Data

func Update_Tracker_Name_By_Id(db *sql.DB, tracker_id int, new_tracker_name string) error {
	sql_string := fmt.Sprintf(
		`UPDATE tracker SET tracker_name = "%s" WHERE tracker_id = %d;`,
		new_tracker_name, tracker_id)

	fmt.Println("SQL:", sql_string)
	_, err := db.Exec(sql_string)
	if err != nil {
		return err
	}

	return nil
}

func Update_Tracker_Notes_By_Name(db *sql.DB, tracker_name string, tracker_notes string) error {
	sql_string := fmt.Sprintf(
		`UPDATE tracker SET tracker_notes = "%s" WHERE tracker_name = "%s";`,
		tracker_notes, tracker_name)

	fmt.Println("SQL:", sql_string)
	_, err := db.Exec(sql_string)
	if err != nil {
		return err
	}

	return nil
}

func Update_Tracker_Notes_By_Id(db *sql.DB, tracker_id int, tracker_notes string) error {
	sql_string := fmt.Sprintf(
		`UPDATE tracker SET tracker_notes = "%s" WHERE tracker_id = %d;`,
		tracker_notes, tracker_id)

	fmt.Println("SQL:", sql_string)
	_, err := db.Exec(sql_string)
	if err != nil {
		return err
	}

	return nil
}

func Update_Number_Field() {
	// TODO
}

func Update_Option_Field() {
	// TODO
}

func Update_Entry(timestamp string) {
	// TODO
}

func Update_Log() {
	// TODO
}

// Delete Data

func Delete_Tracker_By_Id(db *sql.DB, tracker_id int) error {
	sql_string := fmt.Sprintf(
		`DELETE FROM tracker WHERE tracker_id = "%d";`,
		tracker_id)

	fmt.Println("SQL:", sql_string)
	_, err := db.Exec(sql_string)
	if err != nil {
		return err
	}

	return nil
}

func Delete_Tracker_By_Name(db *sql.DB, tracker_name string) error {
	sql_string := fmt.Sprintf(
		`DELETE FROM tracker WHERE tracker_name = "%s";`,
		tracker_name)

	fmt.Println("SQL:", sql_string)
	_, err := db.Exec(sql_string)
	if err != nil {
		return err
	}

	return nil
}

func Delete_Field_By_Id() {
	// TODO
}

func Delete_Entry_By_Id() {
	// TODO
}
