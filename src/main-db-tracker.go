package main

import (
	"database/sql"
	"fmt"
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

func Db_Tracker_Table_Create(db *sql.DB) (err error) {
	_, err = db.Exec(`
		PRAGMA foreign_keys = ON;

		CREATE TABLE IF NOT EXISTS tracker (
			tracker_id INTEGER NOT NULL UNIQUE,
			tracker_name TEXT NOT NULL UNIQUE,
			tracker_notes TEXT NOT NULL DEFAULT "",
			PRIMARY KEY (tracker_id)
		);
		
		CREATE TABLE IF NOT EXISTS field (
			field_id INTEGER NOT NULL UNIQUE,
			tracker_id INTEGER NOT NULL,
			field_type TEXT CHECK(field_type in ('number', 'option')) NOT NULL DEFAULT 'number',
			field_name TEXT NOT NULL,
			field_notes TEXT NOT NULL DEFAULT "",
			UNIQUE(tracker_id, field_name),
			PRIMARY KEY(field_id),
			FOREIGN KEY(tracker_id) REFERENCES tracker (tracker_id) ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS number (
			number_id INTEGER NOT NULL UNIQUE,
			field_id INTEGER NOT NULL,
			decimal_places INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY(number_id),
			FOREIGN KEY(field_id) REFERENCES field (field_id) ON DELETE CASCADE
		);
		
		CREATE TABLE IF NOT EXISTS option (
			option_id INTEGER NOT NULL UNIQUE,
			field_id INTEGER NOT NULL,
			option_value INTEGER NOT NULL DEFAULT 0,
			option_name TEXT NOT NULL DEFAULT "value",
			UNIQUE(field_id, option_name),
			PRIMARY KEY(option_id),
			FOREIGN KEY(field_id) REFERENCES field (field_id) ON DELETE CASCADE
		);
	`)
	return err
}

func Db_Tracker_Create(db *sql.DB, tracker_name string, tracker_notes string) (tracker_id int, err error) {
	sql_string := fmt.Sprintf(
		`INSERT INTO tracker (tracker_name, tracker_notes) VALUES ('%s', '%s');`,
		tracker_name, tracker_notes)

	fmt.Println("SQL:", sql_string)
	result, err := db.Exec(sql_string)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func Db_Tracker_Get(db *sql.DB, tracker_id int) (tracker Db_Tracker, err error) {
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
		return tracker, err
	}
	defer rows.Close()

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
			return tracker, err
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
		return tracker, err
	}

	return tracker, nil
}

func Db_Tracker_All_Get(db *sql.DB) (trackers []Db_Tracker, err error) {
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
		return trackers, err
	}
	defer rows.Close()

	var tracker_scan_last_id = 0
	var field_scan_last_id = 0

	for rows.Next() {
		var tracker_scan Db_Tracker
		var field_scan Db_Field
		var number_scan Db_Number
		var option_scan Db_Option
		err = rows.Scan(
			&tracker_scan.Id, &tracker_scan.Name, &tracker_scan.Notes,
			&field_scan.Id, &field_scan.Type, &field_scan.Name, &field_scan.Notes,
			&number_scan.Id, &number_scan.Decimal_Places,
			&option_scan.Id, &option_scan.Value, &option_scan.Name,
		)
		if err != nil {
			return trackers, err
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
		return trackers, err
	}

	return trackers, nil
}

func Db_Tracker_Id_Get(db *sql.DB, tracker_name string) (tracker_id int, err error) {
	sql_string := fmt.Sprintf(
		`SELECT tracker_id FROM tracker WHERE tracker_name = '%s';`,
		tracker_name)

	row := db.QueryRow(sql_string)

	err = row.Scan(&tracker_id)
	if err != nil {
		return 0, err
	}

	return tracker_id, nil
}

func Db_Tracker_Name_Update(db *sql.DB, tracker_id int, tracker_name string) (err error) {
	sql_string := fmt.Sprintf(
		`UPDATE tracker SET tracker_name = '%s' WHERE tracker_id = %d;`,
		tracker_name, tracker_id)

	fmt.Println("SQL:", sql_string)
	_, err = db.Exec(sql_string)

	return err
}

func Db_Tracker_Notes_Update(db *sql.DB, tracker_id int, tracker_notes string) (err error) {
	sql_string := fmt.Sprintf(
		`UPDATE tracker SET tracker_notes = '%s' WHERE tracker_id = %d;`,
		tracker_notes, tracker_id)

	fmt.Println("SQL:", sql_string)
	_, err = db.Exec(sql_string)

	return err
}

func Db_Tracker_Delete(db *sql.DB, tracker_id int) (err error) {
	sql_string := fmt.Sprintf(
		`DELETE FROM tracker WHERE tracker_id = "%d";`,
		tracker_id)

	fmt.Println("SQL:", sql_string)
	_, err = db.Exec(sql_string)
	if err != nil {
		return err
	}

	return nil
}

func Db_Tracker_Field_Number_Create(db *sql.DB, tracker_id int, field_name string, field_notes string, decimal_places int) (field_id int, err error) {
	sql_string_field := fmt.Sprintf(
		`INSERT INTO field (tracker_id, field_type, field_name, field_notes) VALUES (%d,"number",'%s','%s');`,
		tracker_id, field_name, field_notes)

	fmt.Println("SQL:", sql_string_field)
	result, err := db.Exec(sql_string_field)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	field_id = int(id)

	sql_string_number := fmt.Sprintf(
		`INSERT INTO number (field_id, decimal_places) VALUES (%d,%d);`,
		field_id, decimal_places)

	fmt.Println("SQL:", sql_string_number)
	_, err = db.Exec(sql_string_number)
	if err != nil {
		return 0, err
	}

	return field_id, nil
}

func Db_Tracker_Field_Number_Update(db *sql.DB, field_id int, decimal_places int) (err error) {
	sql_string := fmt.Sprintf(
		`UPDATE log
		SET log_value = log_value * power(10, %d - (SELECT decimal_places FROM number WHERE field_id = %d))
		WHERE field_id = %d;

		UPDATE number
		SET decimal_places = %d
		WHERE field_id = %d;`,
		decimal_places, field_id, field_id, decimal_places, field_id)

	fmt.Println("SQL:", sql_string)
	_, err = db.Exec(sql_string)

	return err
}

func Db_Tracker_Field_Option_Create(db *sql.DB, tracker_id int, field_name string, field_notes string, options []struct { Value int; Name string }) (field_id int, err error) {
	sql_string_field := fmt.Sprintf(
		`INSERT INTO field (tracker_id, field_type, field_name, field_notes) VALUES (%d,"option",'%s','%s');`,
		tracker_id, field_name, field_notes)

	fmt.Println("SQL:", sql_string_field)
	result, err := db.Exec(sql_string_field)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	field_id = int(id)

	for _, option := range options {
		sql_string_option := fmt.Sprintf(
			`INSERT INTO option (field_id, option_value, option_name) VALUES (%d,%d,'%s');`,
			field_id, option.Value, option.Name)

		fmt.Println("SQL:", sql_string_option)
		_, err = db.Exec(sql_string_option)
		if err != nil {
			return 0, err
		}
	}

	return field_id, nil
}

func Db_Tracker_Field_Option_Add(db *sql.DB, field_id int, option_value int, option_name string) (option_id int, err error) {
	sql_string := fmt.Sprintf(
		`INSERT INTO option (field_id, option_value, option_name) VALUES (%d,%d,'%s');`,
		field_id, option_value, option_name)

	fmt.Println("SQL:", sql_string)
	result, err := db.Exec(sql_string)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	option_id = int(id)

	return option_id, nil
}

func Db_Tracker_Field_Option_Value_Update(db *sql.DB, option_id int, option_value int) (err error) {
	sql_string := fmt.Sprintf(
		`UPDATE log
		SET log_value = %d
		FROM (SELECT field_id, option_value FROM option WHERE option_id = %d) AS o
		WHERE log.field_id = o.field_id AND log.log_value = o.option_value;

		UPDATE option SET option_value = %d WHERE option_id = %d;`,
		option_value, option_id, option_value, option_id)

	fmt.Println("SQL:", sql_string)
	_, err = db.Exec(sql_string)

	return err
}

func Db_Tracker_Field_Option_Name_Update(db *sql.DB, option_id int, option_name string) (err error) {
	sql_string := fmt.Sprintf(
		`UPDATE option SET option_name = '%s' WHERE option_id = %d;`,
		option_name, option_id)

	fmt.Println("SQL:", sql_string)
	_, err = db.Exec(sql_string)

	return err
}

func Db_Tracker_Field_Option_Remove(db *sql.DB, option_id int) (err error) {
	sql_string := fmt.Sprintf(
		`DELETE FROM entry WHERE entry_id = (
			SELECT entry_id FROM log
			WHERE log.field_id = (SELECT field_id FROM option WHERE option_id = %d)
			AND log.log_value = (SELECT option_value FROM option WHERE option_id = %d)
		);
		DELETE FROM option WHERE option_id = %d;`,
		option_id, option_id, option_id)

	fmt.Println("SQL:", sql_string)
	_, err = db.Exec(sql_string)

	return err
}

func Db_Tracker_Field_Name_Update(db *sql.DB, field_id int, field_name string) (err error) {
	sql_string := fmt.Sprintf(
		`UPDATE field SET field_name = '%s' WHERE field_id = %d;`,
		field_name, field_id)

	fmt.Println("SQL:", sql_string)
	_, err = db.Exec(sql_string)

	return err
}

func Db_Tracker_Field_Notes_Update(db *sql.DB, field_id int, field_notes string) (err error) {
	sql_string := fmt.Sprintf(
		`UPDATE field SET field_notes = '%s' WHERE field_id = %d;`,
		field_notes, field_id)

	fmt.Println("SQL:", sql_string)
	_, err = db.Exec(sql_string)

	return err
}

func Db_Tracker_Field_Delete(db *sql.DB, field_id int) (err error) {
	sql_string := fmt.Sprintf(
		`DELETE FROM field WHERE field_id = "%d";
		DELETE FROM log WHERE field_id = "%d";`,
		field_id, field_id)

	fmt.Println("SQL:", sql_string)
	_, err = db.Exec(sql_string)
	if err != nil {
		return err
	}

	return nil
}
