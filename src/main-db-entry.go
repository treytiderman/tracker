package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

type Db_Entry struct {
	Id         int
	Tracker_Id int
	Timestamp  string
	Notes      string
	Logs       []Db_Log
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

func Db_Entry_Table_Create(db *sql.DB) (err error) {
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS entry (
			entry_id INTEGER NOT NULL UNIQUE,
			tracker_id INTEGER NOT NULL,
			timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			entry_notes TEXT NOT NULL DEFAULT "",
			PRIMARY KEY(entry_id),
			FOREIGN KEY(tracker_id) REFERENCES tracker (tracker_id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS log (
			log_id INTEGER NOT NULL UNIQUE,
			entry_id INTEGER NOT NULL,
			field_id INTEGER NOT NULL,
			log_value INTEGER NOT NULL,
			PRIMARY KEY(log_id),
			FOREIGN KEY(entry_id) REFERENCES entry (entry_id) ON DELETE CASCADE,
			FOREIGN KEY(field_id) REFERENCES field (field_id) ON DELETE CASCADE
		);
	`)
	return err
}

func Db_Entry_Table_Delete(db *sql.DB) (err error) {
	_, err = db.Exec(`
		DROP TABLE IF EXISTS log;
		DROP TABLE IF EXISTS entry;
	`)
	return err
}

func Db_Entry_Create(db *sql.DB, tracker_id int, entry_notes string, logs []struct {
	Field_Id int
	Value    int
}) (entry_id int, err error) {
	fmt.Printf("SQL: INSERT INTO entry (tracker_id, entry_notes) VALUES (%d,'%s');", tracker_id, entry_notes)
	result, err := db.Exec("INSERT INTO entry (tracker_id, entry_notes) VALUES (?,?);", tracker_id, entry_notes)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	entry_id = int(id)

	for _, log := range logs {
		fmt.Printf("SQL: INSERT INTO log (entry_id, field_id, log_value) VALUES (%d,%d,%d);", entry_id, log.Field_Id, log.Value)
		_, err = db.Exec("INSERT INTO log (entry_id, field_id, log_value) VALUES (?,?,?);", entry_id, log.Field_Id, log.Value)
		if err != nil {
			return 0, err
		}
	}

	return entry_id, nil
}

func Db_Entry_Create_Timestamp(db *sql.DB, tracker_id int, entry_notes string, logs []struct {
	Field_Id int
	Value    int
}, timestamp string) (entry_id int, err error) {
	fmt.Printf("SQL: INSERT INTO entry (tracker_id, entry_notes, timestamp) VALUES (%d,'%s','%s');", tracker_id, entry_notes, timestamp)
	result, err := db.Exec("INSERT INTO entry (tracker_id, entry_notes, timestamp) VALUES (?,?,?);", tracker_id, entry_notes, timestamp)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	entry_id = int(id)

	for _, log := range logs {
		fmt.Printf("SQL: INSERT INTO log (entry_id, field_id, log_value) VALUES (%d,%d,%d);", entry_id, log.Field_Id, log.Value)
		_, err = db.Exec("INSERT INTO log (entry_id, field_id, log_value) VALUES (?,?,?);", entry_id, log.Field_Id, log.Value)
		if err != nil {
			return 0, err
		}
	}

	return entry_id, nil
}

func Db_Entry_Get_By_Entry_Id(db *sql.DB, entry_id int) (Db_Entry, error) {
	entries := make([]Db_Entry, 0)

	sql_string := fmt.Sprintf(
		`SELECT
			entry.entry_id,
			entry.tracker_id,
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
				printf(("%%." || number.decimal_places || "f"), log.log_value / power(10, number.decimal_places))
			ELSE
				option.option_name
			END), "") AS present
		FROM entry
		LEFT JOIN log USING (entry_id)
		LEFT JOIN field USING (field_id)
		LEFT JOIN number USING (field_id)
		LEFT JOIN option ON log.field_id = option.field_id AND log.log_value = option.option_value
		WHERE entry.entry_id = ?;`)

	rows, err := db.Query(sql_string, entry_id)
	if err != nil {
		return entries[0], err
	}
	defer rows.Close()

	var log_scan_last_id = 0
	var entry_scan_last_id = 0

	for rows.Next() {
		var entry_scan Db_Entry
		var log_scan Db_Log
		err = rows.Scan(
			&entry_scan.Id, &entry_scan.Tracker_Id, &entry_scan.Timestamp, &entry_scan.Notes,
			&log_scan.Id, &log_scan.Value, &log_scan.Field_Id, &log_scan.Field_Type, &log_scan.Field_Name,
			&log_scan.Decimal_Places, &log_scan.Option_Value, &log_scan.Option_Name, &log_scan.Present,
		)
		if err != nil {
			return entries[0], err
		}

		log.Println(entry_scan, log_scan)

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
		return entries[0], err
	}

	return entries[0], nil
}

func Db_Entry_Get(db *sql.DB, tracker_id int) (entries []Db_Entry, err error) {
	sql_string := fmt.Sprintf(
		`SELECT
			entry.entry_id,
			entry.tracker_id,
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
				printf(("%%." || number.decimal_places || "f"), log.log_value / power(10, number.decimal_places))
			ELSE
				option.option_name
			END), "") AS present
		FROM entry
		LEFT JOIN log USING (entry_id)
		LEFT JOIN field USING (field_id)
		LEFT JOIN number USING (field_id)
		LEFT JOIN option ON log.field_id = option.field_id AND log.log_value = option.option_value
		WHERE entry.tracker_id = %d
		ORDER BY entry.entry_id DESC, field.field_id;`,
		tracker_id)

	rows, err := db.Query(sql_string)
	if err != nil {
		return entries, err
	}
	defer rows.Close()

	var log_scan_last_id = 0
	var entry_scan_last_id = 0

	for rows.Next() {
		var entry_scan Db_Entry
		var log_scan Db_Log
		err = rows.Scan(
			&entry_scan.Id, &entry_scan.Tracker_Id, &entry_scan.Timestamp, &entry_scan.Notes,
			&log_scan.Id, &log_scan.Value, &log_scan.Field_Id, &log_scan.Field_Type, &log_scan.Field_Name,
			&log_scan.Decimal_Places, &log_scan.Option_Value, &log_scan.Option_Name, &log_scan.Present,
		)
		if err != nil {
			return entries, err
		}

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
		return entries, err
	}

	return entries, nil
}

func Db_Entry_Filter_Notes_Get(db *sql.DB, tracker_id int, search string) (entries []Db_Entry, err error) {
	search_pattern := "%" + search + "%"
	sql_string := fmt.Sprintf(
		`SELECT
			entry.entry_id,
			entry.tracker_id,
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
				printf(("%%." || number.decimal_places || "f"), log.log_value / power(10, number.decimal_places))
			ELSE
				option.option_name
			END), "") AS present
		FROM entry
		LEFT JOIN log USING (entry_id)
		LEFT JOIN field USING (field_id)
		LEFT JOIN number USING (field_id)
		LEFT JOIN option ON log.field_id = option.field_id AND log.log_value = option.option_value
		WHERE entry.tracker_id = %d AND entry.entry_notes LIKE '%s'
		ORDER BY entry.entry_id DESC, field.field_id;`,
		tracker_id, search_pattern)

	rows, err := db.Query(sql_string)
	if err != nil {
		return entries, err
	}
	defer rows.Close()

	var log_scan_last_id = 0
	var entry_scan_last_id = 0

	for rows.Next() {
		var entry_scan Db_Entry
		var log_scan Db_Log
		err = rows.Scan(
			&entry_scan.Id, &entry_scan.Tracker_Id, &entry_scan.Timestamp, &entry_scan.Notes,
			&log_scan.Id, &log_scan.Value, &log_scan.Field_Id, &log_scan.Field_Type, &log_scan.Field_Name,
			&log_scan.Decimal_Places, &log_scan.Option_Value, &log_scan.Option_Name, &log_scan.Present,
		)
		if err != nil {
			return entries, err
		}

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
		return entries, err
	}

	return entries, nil
}

func Db_Entry_All_Get(db *sql.DB) (entries []Db_Entry, err error) {
	sql_string := `
		SELECT
			entry.entry_id,
			entry.tracker_id,
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
				printf(("%." || number.decimal_places || "f"), log.log_value / power(10, number.decimal_places))
			ELSE
				option.option_name
			END), "") AS present
		FROM entry
		LEFT JOIN log USING (entry_id)
		LEFT JOIN field USING (field_id)
		LEFT JOIN number USING (field_id)
		LEFT JOIN option ON log.field_id = option.field_id AND log.log_value = option.option_value
		ORDER BY entry.entry_id DESC, field.field_id;`

	rows, err := db.Query(sql_string)
	if err != nil {
		return entries, err
	}
	defer rows.Close()

	var log_scan_last_id = 0
	var entry_scan_last_id = 0

	for rows.Next() {
		var entry_scan Db_Entry
		var log_scan Db_Log
		err = rows.Scan(
			&entry_scan.Id, &entry_scan.Tracker_Id, &entry_scan.Timestamp, &entry_scan.Notes,
			&log_scan.Id, &log_scan.Value, &log_scan.Field_Id, &log_scan.Field_Type, &log_scan.Field_Name,
			&log_scan.Decimal_Places, &log_scan.Option_Value, &log_scan.Option_Name, &log_scan.Present,
		)
		if err != nil {
			return entries, err
		}

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
		return entries, err
	}

	return entries, nil
}

func Db_Entry_All_Filter_Notes_Get(db *sql.DB, search string) (entries []Db_Entry, err error) {
	search_pattern := "%" + search + "%"
	sql_string := fmt.Sprintf(
		`SELECT
			entry.entry_id,
			entry.tracker_id,
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
				printf(("%%." || number.decimal_places || "f"), log.log_value / power(10, number.decimal_places))
			ELSE
				option.option_name
			END), "") AS present
		FROM entry
		LEFT JOIN log USING (entry_id)
		LEFT JOIN field USING (field_id)
		LEFT JOIN number USING (field_id)
		LEFT JOIN option ON log.field_id = option.field_id AND log.log_value = option.option_value
		WHERE entry.entry_notes LIKE '%s'
		ORDER BY entry.entry_id DESC, field.field_id;`,
		search_pattern)

	rows, err := db.Query(sql_string)
	if err != nil {
		return entries, err
	}
	defer rows.Close()

	var log_scan_last_id = 0
	var entry_scan_last_id = 0

	for rows.Next() {
		var entry_scan Db_Entry
		var log_scan Db_Log
		err = rows.Scan(
			&entry_scan.Id, &entry_scan.Tracker_Id, &entry_scan.Timestamp, &entry_scan.Notes,
			&log_scan.Id, &log_scan.Value, &log_scan.Field_Id, &log_scan.Field_Type, &log_scan.Field_Name,
			&log_scan.Decimal_Places, &log_scan.Option_Value, &log_scan.Option_Name, &log_scan.Present,
		)
		if err != nil {
			return entries, err
		}

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
		return entries, err
	}

	return entries, nil
}

func Db_Entry_Timestamp_Update(db *sql.DB, entry_id int, timestamp string) (err error) {
	sql_string := fmt.Sprintf(
		`UPDATE entry SET timestamp = '%s' WHERE entry_id = %d;`,
		timestamp, entry_id)

	fmt.Println("SQL:", sql_string)
	_, err = db.Exec(sql_string)

	return err
}

func Db_Entry_Notes_Update(db *sql.DB, entry_id int, entry_notes string) (err error) {
	fmt.Printf("SQL: UPDATE entry SET entry_notes = '%s' WHERE entry_id = %d;", entry_notes, entry_id)
	_, err = db.Exec("UPDATE entry SET entry_notes = ? WHERE entry_id = ?;", entry_notes, entry_id)
	return err
}

func Db_Entry_Log_Update(db *sql.DB, log_id int, log_value int) (err error) {
	sql_string := fmt.Sprintf(
		`UPDATE log SET log_value = %d WHERE log_id = %d;`,
		log_value, log_id)

	fmt.Println("SQL:", sql_string)
	_, err = db.Exec(sql_string)

	return err
}

func Db_Entry_Delete(db *sql.DB, entry_id int) (err error) {
	sql_string := fmt.Sprintf(
		`DELETE FROM entry WHERE entry_id = %d; DELETE FROM log WHERE entry_id = %d;`,
		entry_id, entry_id)

	fmt.Println("SQL:", sql_string)
	_, err = db.Exec(sql_string)

	return err
}
