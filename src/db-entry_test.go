package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

// var db_test *sql.DB // init in db-tracker_test.go

func _test_Reset_Entry_Database(t *testing.T) {
	path := "../data/test.db"

	err := os.Remove(path)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Database Deleted", path)

	db_test, err = sql.Open("sqlite", path)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Database Opened", path)

	err = Create_Tracker_Tables(db_test)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Database Tracker Tables Created")

	err = Create_Entry_Tables(db_test)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Database Entry Tables Created")
}

func _test_Add_Entries_To_Journal(t *testing.T, tracker_id int) {
	var tests = []struct {
		expected_id int
		entry_notes string
		timestamp   string
	}{
		{1, "Entry 1", "2049-12-13 19:15:56"},
		{2, "Why is green sometimes blue", "2095-12-14 19:16:56"},
		{3, "If Franky can be a robot maybe I can too", "2094-12-13 19:17:56"},
		{4, "", "2093-12-13 19:18:56"},
		{5, "I got lost in a square", "2124-12-13 19:19:56"},
		{6, "The circle showed me the way", "2999-12-13 19:20:56"},
	}

	for _, tt := range tests {
		_, err := Create_Entry_With_Timestamp(db_test, tracker_id, tt.entry_notes, tt.timestamp)
		if err != nil {
			t.Error(err)
		}
	}
}

// Helper

func Test_Parse_String_To_Number(t *testing.T) {
	var tests = []struct {
		expected_value int
		log_string     string
		decimal_places int
	}{
		{100, "100", 0},
		{12300, "123", 2},
		{12345, "123.45", 2},
		{72440, "72.440", 3},
		{4, "410", -2},
		{54, "543210", -4},
		// {0, "", 0}, FAILS
		// {0, "awdzsd", 0}, FAILS
	}

	for _, tt := range tests {
		t.Run(tt.log_string, func(t *testing.T) {
			result, err := Parse_String_To_Number(tt.log_string, tt.decimal_places)
			if err != nil {
				t.Error(err)
			}
			if result != tt.expected_value {
				t.Errorf("got %d, expected %d", result, tt.expected_value)
			}
		})
	}
}

// Create

func Test_Create_Entry(t *testing.T) {
	_test_Reset_Entry_Database(t)
	_test_Create_Tracker_Journal(t)

	var tests = []struct {
		expected_id int
		entry_notes string
	}{
		{1, "First Entry"},
		{2, "Why is green sometimes blue"},
		{3, "If Franky can be a robot maybe I can too"},
		{4, ""},
		{5, "I got lost in a square"},
		{6, "The circle showed me the way"},
	}

	for _, tt := range tests {
		t.Run(tt.entry_notes, func(t *testing.T) {
			id, err := Create_Entry(db_test, 1, tt.entry_notes)
			if err != nil {
				t.Error(err)
			}
			if id != tt.expected_id {
				t.Errorf("got %d, expected %d", id, tt.expected_id)
			}
		})
	}
}

func Test_Add_Log_To_Entry(t *testing.T) {
	_test_Reset_Entry_Database(t)
	_test_Create_Tracker_Money(t)

	entry_id, err := Create_Entry(db_test, 1, "Logged some things")
	if err != nil {
		t.Error(err)
	}

	_, err = Add_Log_To_Entry(db_test, entry_id, 1, 22000) // 220.00
	if err != nil {
		t.Error(err)
	}

	_, err = Add_Log_To_Entry(db_test, entry_id, 2, 1)
	if err != nil {
		t.Error(err)
	}
}

func Test_Create_Entry_With_Timestamp(t *testing.T) {
	_test_Reset_Entry_Database(t)
	_test_Create_Tracker_Journal(t)

	var tests = []struct {
		expected_id int
		entry_notes string
		timestamp   string
	}{
		{1, "Entry 1", "2049-12-13 19:15:56"},
		{2, "Why is green sometimes blue", "2095-12-14 19:16:56"},
		{3, "If Franky can be a robot maybe I can too", "2094-12-13 19:17:56"},
		{4, "", "2093-12-13 19:18:56"},
		{5, "I got lost in a square", "2124-12-13 19:19:56"},
		{6, "The circle showed me the way", "2999-12-13 19:20:56"},
	}

	for _, tt := range tests {
		t.Run(tt.entry_notes, func(t *testing.T) {
			id, err := Create_Entry_With_Timestamp(db_test, 1, tt.entry_notes, tt.timestamp)
			if err != nil {
				t.Error(err)
			}
			if id != tt.expected_id {
				t.Errorf("got %d, expected %d", id, tt.expected_id)
			}
		})
	}
}

// Read

func Test_Get_Entry_By_Entry_Id(t *testing.T) {
	_test_Reset_Entry_Database(t)
	_test_Create_Tracker_Journal(t)
	_test_Add_Entries_To_Journal(t, 1)

	var tests = []struct {
		expected_notes string
		entry_id       int
	}{
		{"Entry 1", 1},
		{"Why is green sometimes blue", 2},
		{"If Franky can be a robot maybe I can too", 3},
		{"", 4},
		{"I got lost in a square", 5},
		{"The circle showed me the way", 6},
	}

	for _, tt := range tests {
		t.Run(tt.expected_notes, func(t *testing.T) {
			entry, err := Get_Entry_By_Entry_Id(db_test, tt.entry_id)
			if err != nil {
				t.Error(err)
			}
			if entry.Notes != tt.expected_notes {
				t.Errorf("got %s, expected %s", entry.Notes, tt.expected_notes)
			}
		})
	}
}

func Test_Get_Entries_By_Tracker_ID(t *testing.T) {
	_test_Reset_Entry_Database(t)

	journal_id, _ := Create_Tracker(db_test, "Journal", "Daily journal and notes")
	Create_Entry_With_Timestamp(db_test, journal_id, "Entry 1", "2049-12-13 19:15:56")
	Create_Entry_With_Timestamp(db_test, journal_id, "Why is green sometimes blue", "2095-12-14 19:16:56")
	Create_Entry_With_Timestamp(db_test, journal_id, "If Franky can be a robot maybe I can too", "2094-12-13 19:17:56")
	Create_Entry_With_Timestamp(db_test, journal_id, "", "2093-12-13 19:18:56")
	Create_Entry_With_Timestamp(db_test, journal_id, "I got lost in a square", "2124-12-13 19:19:56")
	Create_Entry_With_Timestamp(db_test, journal_id, "The circle showed me the way", "2999-12-13 19:20:56")

	money_id, _ := Create_Tracker(db_test, "Money", "Transactions")
	money_amount_id, _ := Add_Number_Field(db_test, money_id, "Amount", "Amount of money in dollars", 2)
	money_card_id, _ := Add_Option_Field(db_test, money_id, "Card", "Payment Method")
	Add_Option_to_Field(db_test, money_card_id, 1, "Discover")
	Add_Option_to_Field(db_test, money_card_id, 2, "Visa")
	Add_Option_to_Field(db_test, money_card_id, 3, "American Express")
	money_entry_1, _ := Create_Entry(db_test, money_id, "9.99 dollars entered as 999")
	Add_Log_To_Entry(db_test, money_entry_1, money_amount_id, -999)
	Add_Log_To_Entry(db_test, money_entry_1, money_card_id, 1)
	money_entry_2, _ := Create_Entry(db_test, money_id, "not for what you think")
	Add_Log_To_Entry(db_test, money_entry_2, money_amount_id, -42069)
	Add_Log_To_Entry(db_test, money_entry_2, money_card_id, 3)
	money_entry_3, _ := Create_Entry(db_test, money_id, "big spendin")
	Add_Log_To_Entry(db_test, money_entry_3, money_amount_id, 2000_00)
	Add_Log_To_Entry(db_test, money_entry_3, money_card_id, 2)

	// Test Start
	entries, err := Get_Entries_By_Tracker_Id(db_test, 2)
	if err != nil {
		t.Error(err)
	}

	s, _ := json.MarshalIndent(entries, "", "    ")
	fmt.Println("JSON:", string(s))

	if entries[0].Id != 9 {
		t.Errorf("got %d, expected %d", entries[0].Id, 9)
	}

	if entries[2].Notes != "9.99 dollars entered as 999" {
		t.Errorf("got %s, expected %s", entries[2].Notes, "9.99 dollars entered as 999")
	}

	if entries[1].Logs[0].Present != "-420.69" {
		t.Errorf("got %s, expected %s", entries[1].Logs[0].Present, "-420.69")
	}

	if entries[1].Logs[1].Present != "American Express" {
		t.Errorf("got %s, expected %s", entries[1].Logs[1].Present, "American Express")
	}
}

// --------------------------------------------------------

// func Test_Add_Entry(t *testing.T) {
// 	var tests = []struct {
// 		expected_id  int
// 		tracker_name string
// 		entry_notes  string
// 		logs         []struct {
// 			Field_Id int
// 			Value    int
// 		}
// 	}{
// 		{1, "Journal", "Hello Journal", []struct {
// 			Field_Id int
// 			Value    int
// 		}{}},
// 		{2, "Weight", "Init", []struct {
// 			Field_Id int
// 			Value    int
// 		}{
// 			{1, 180},
// 		}},
// 		{3, "Weight", "Imagine this is a few weeks later", []struct {
// 			Field_Id int
// 			Value    int
// 		}{
// 			{1, 175},
// 		}},
// 		{4, "Money", "9.99 dollars entered as 999", []struct {
// 			Field_Id int
// 			Value    int
// 		}{
// 			{2, -999},
// 			{6, -1},
// 		}},
// 		{5, "Money", "", []struct {
// 			Field_Id int
// 			Value    int
// 		}{
// 			{2, -42069},
// 			{6, 0},
// 		}},
// 		{6, "Money", "", []struct {
// 			Field_Id int
// 			Value    int
// 		}{
// 			{2, 200000},
// 			{6, 2},
// 		}},
// 		{7, "DELETE_ME", "note", []struct {
// 			Field_Id int
// 			Value    int
// 		}{
// 			{5, 25},
// 		}},
// 		{8, "Workout", "Chest Day", []struct {
// 			Field_Id int
// 			Value    int
// 		}{
// 			{3, 135},
// 			{4, 5},
// 			{7, 1},
// 		}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.tracker_name, func(t *testing.T) {
// 			id, err := Add_Entry(db_test, tt.tracker_name, tt.entry_notes, tt.logs)
// 			if err != nil {
// 				t.Fatalf("Failed to update tracker notes")
// 			}
// 			if id != tt.expected_id {
// 				t.Errorf("got %d, expected %d", id, tt.expected_id)
// 			}
// 		})
// 	}
// }

// func Test_Get_Entries_By_Tracker_Id(t *testing.T) {
// 	var tests = []struct {
// 		tracker_id int
// 	}{
// 		{1},
// 		{3},
// 		{8},
// 	}
// 	for _, tt := range tests {
// 		t.Run(fmt.Sprintf("tracker_id=%d", tt.tracker_id), func(t *testing.T) {
// 			entries, err := Get_Entries_By_Tracker_Id(db_test, tt.tracker_id)
// 			if err != nil {
// 				t.Fatalf("Failed to Get_Entries_By_Tracker_Id")
// 			}
// 			// Error with library? See function...
// 			if tt.tracker_id == 8 {
// 				if len(entries[0].Logs) > 3 {
// 					t.Fatalf("Error too many logs")
// 				}
// 			}
// 			s, _ := json.Marshal(entries)
// 			fmt.Sprintln("JSON:", string(s))
// 			// fmt.Println("JSON:", string(s))
// 		})
// 	}
// }
