package main

import (
	"database/sql"
	"log/slog"
	"math"
	"strconv"

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

// Helper

func Parse_String_To_Number(str string, decimal_places int) (int, error) {
	log_value_float, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, err
	}

	log_value_adjusted := float64(log_value_float) * float64(math.Pow10(decimal_places))
	log_value_int := int(math.Floor(log_value_adjusted))
	return log_value_int, nil
}

// Create

func Create_Entry_Tables(db *sql.DB) error {
	_, err := db.Exec(`
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

func Create_Entry(db *sql.DB, tracker_id int, entry_notes string) (int, error) {
	slog.Debug("database create entry", "tracker_id", tracker_id, "entry_notes", entry_notes)

	result, err := db.Exec(
		"INSERT INTO entry (tracker_id, entry_notes) VALUES (?,?);",
		tracker_id, entry_notes)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func Add_Log_To_Entry(db *sql.DB, entry_id int, field_id int, log_value int) (int, error) {
	slog.Debug("database add log to entry", "entry_id", entry_id, "field_id", field_id, "log_value", log_value)

	result, err := db.Exec(
		"INSERT INTO log (entry_id, field_id, log_value) VALUES (?,?,?);",
		entry_id, field_id, log_value)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func Create_Entry_With_Logs(db *sql.DB, tracker_id int, entry_notes string, logs []struct {
	Value    int
	Field_Id int
}) (int, error) {
	entry_id, err := Create_Entry(db, tracker_id, entry_notes)
	if err != nil {
		return 0, err
	}
	
	for _, log := range logs {
		_, err := Add_Log_To_Entry(db, entry_id, log.Field_Id, log.Value)
		if err != nil {
			return 0, err
		}
	}

	return entry_id, nil
}

// Read

func Get_Entry(db *sql.DB, entry_id int) (Db_Entry, error) {
	entries := make([]Db_Entry, 0)

	rows, err := db.Query(`SELECT
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
				printf(("%." || (
					CASE WHEN number.decimal_places < 0 THEN 0 ELSE number.decimal_places END
				) || "f"), log.log_value / power(10, number.decimal_places))
			ELSE
				option.option_name
			END), "") AS present
		FROM entry
		LEFT JOIN log USING (entry_id)
		LEFT JOIN field USING (field_id)
		LEFT JOIN number USING (field_id)
		LEFT JOIN option ON log.field_id = option.field_id AND log.log_value = option.option_value
		WHERE entry.entry_id = ?;`, entry_id)
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

// In reverse entry_id order
func Get_Entries(db *sql.DB, tracker_id int) ([]Db_Entry, error) {
	entries := make([]Db_Entry, 0)

	rows, err := db.Query(
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
				printf(("%." || (
					CASE WHEN number.decimal_places < 0 THEN 0 ELSE number.decimal_places END
				) || "f"), log.log_value / power(10, number.decimal_places))
			ELSE
				option.option_name
			END), "") AS present
		FROM entry
		LEFT JOIN log USING (entry_id)
		LEFT JOIN field USING (field_id)
		LEFT JOIN number USING (field_id)
		LEFT JOIN option ON log.field_id = option.field_id AND log.log_value = option.option_value
		WHERE entry.tracker_id = ?
		ORDER BY entry.entry_id DESC, field.field_id;`,
		tracker_id)
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

func Get_Entries_Filter(db *sql.DB, tracker_id int, search string) ([]Db_Entry, error) {
	entries := make([]Db_Entry, 0)
	search_pattern := "%" + search + "%" // contains

	rows, err := db.Query(
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
				printf(("%." || (
					CASE WHEN number.decimal_places < 0 THEN 0 ELSE number.decimal_places END
				) || "f"), log.log_value / power(10, number.decimal_places))
			ELSE
				option.option_name
			END), "") AS present
		FROM entry
		LEFT JOIN log USING (entry_id)
		LEFT JOIN field USING (field_id)
		LEFT JOIN number USING (field_id)
		LEFT JOIN option ON log.field_id = option.field_id AND log.log_value = option.option_value
		WHERE entry.tracker_id = ? AND entry.entry_notes LIKE ?
		ORDER BY entry.entry_id DESC, field.field_id;`,
		tracker_id, search_pattern)
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

func Get_All_Entries(db *sql.DB) ([]Db_Entry, error) {
	entries := make([]Db_Entry, 0)

	rows, err := db.Query(
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
				printf(("%." || (
					CASE WHEN number.decimal_places < 0 THEN 0 ELSE number.decimal_places END
				) || "f"), log.log_value / power(10, number.decimal_places))
			ELSE
				option.option_name
			END), "") AS present
		FROM entry
		LEFT JOIN log USING (entry_id)
		LEFT JOIN field USING (field_id)
		LEFT JOIN number USING (field_id)
		LEFT JOIN option ON log.field_id = option.field_id AND log.log_value = option.option_value
		ORDER BY entry.entry_id DESC, field.field_id;`)
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

func Get_Log(db *sql.DB, log_id int) (Db_Log, error) {
	row := db.QueryRow(
		`SELECT
			IFNULL(log.log_id, 0) AS log_id,
			IFNULL(log.log_value, 0) AS log_value,
			IFNULL(field.field_id, 0) AS field_id,
			IFNULL(field.field_type, "") AS field_type,
			IFNULL(field.field_name, "") AS field_name,
			IFNULL(number.decimal_places, 0) AS decimal_places,
			IFNULL(option.option_value, 0) AS option_value,
			IFNULL(option.option_name, "") AS option_name,
			IFNULL((CASE WHEN field.field_type == "number" THEN
				printf(("%." || (
					CASE WHEN number.decimal_places < 0 THEN 0 ELSE number.decimal_places END
				) || "f"), log.log_value / power(10, number.decimal_places))
			ELSE
				option.option_name
			END), "") AS present
		FROM entry
		LEFT JOIN log USING (entry_id)
		LEFT JOIN field USING (field_id)
		LEFT JOIN number USING (field_id)
		LEFT JOIN option ON log.field_id = option.field_id AND log.log_value = option.option_value
		WHERE log.log_id = ?
		ORDER BY field.field_id;`,
		log_id)

	var log_scan Db_Log
	err := row.Scan(
		&log_scan.Id, &log_scan.Value, &log_scan.Field_Id, &log_scan.Field_Type, &log_scan.Field_Name,
		&log_scan.Decimal_Places, &log_scan.Option_Value, &log_scan.Option_Name, &log_scan.Present)
	if err != nil {
		return log_scan, err
	}

	return log_scan, nil
}

// Update

func Update_Entry_Timestamp(db *sql.DB, entry_id int, timestamp string) error {
	slog.Debug("database update entry timestamp", "entry_id", entry_id, "timestamp", timestamp)
	_, err := db.Exec("UPDATE entry SET timestamp = ? WHERE entry_id = ?;", timestamp, entry_id)
	return err
}

func Update_Entry_Notes(db *sql.DB, entry_id int, entry_notes string) error {
	slog.Debug("database update entry notes", "entry_id", entry_id, "entry_notes", entry_notes)
	_, err := db.Exec("UPDATE entry SET entry_notes = ? WHERE entry_id = ?;", entry_notes, entry_id)
	return err
}

func Update_Log(db *sql.DB, log_id int, log_value int) error {
	slog.Debug("database update log", "log_id", log_id, "log_value", log_value)
	_, err := db.Exec("UPDATE log SET log_value = ? WHERE log_id = ?;", log_value, log_id)
	return err
}

// Delete

func Delete_Entry(db *sql.DB, entry_id int) error {
	slog.Debug("database delete entry", "entry_id", entry_id)
	_, err := db.Exec("DELETE FROM entry WHERE entry_id = ?; DELETE FROM log WHERE entry_id = ?;", entry_id, entry_id)
	return err
}
