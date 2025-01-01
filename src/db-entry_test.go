package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

// var db_test *sql.DB // init in db-tracker_test.go

func _test_Reset_Entry_Database(t *testing.T) {
	path := "../data/test.db"

	os.Remove(path)
	fmt.Println("DATABASE Deleted", path)

	var err error
	db_test, err = sql.Open("sqlite", path)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("DATABASE Opened", path)

	err = Create_Tracker_Tables(db_test)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("DATABASE Tracker Tables Created")

	err = Create_Entry_Tables(db_test)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("DATABASE Entry Tables Created")
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

// Read

func Test_Get_Entry(t *testing.T) {
	_test_Reset_Entry_Database(t)

	journal_id, _ := Create_Tracker(db_test, "Journal", "Daily journal and notes")
	Create_Entry(db_test, journal_id, "Entry 1")
	Create_Entry(db_test, journal_id, "Why is green sometimes blue")
	Create_Entry(db_test, journal_id, "If Franky can be a robot maybe I can too")
	Create_Entry(db_test, journal_id, "")
	Create_Entry(db_test, journal_id, "I got lost in a square")
	Create_Entry(db_test, journal_id, "The circle showed me the way")

	var tests = []struct {
		entry_id       int
		expected_notes string
	}{
		{1, "Entry 1"},
		{2, "Why is green sometimes blue"},
		{3, "If Franky can be a robot maybe I can too"},
		{4, ""},
		{5, "I got lost in a square"},
		{6, "The circle showed me the way"},
	}

	for _, tt := range tests {
		t.Run(tt.expected_notes, func(t *testing.T) {
			entry, err := Get_Entry(db_test, tt.entry_id)
			if err != nil {
				t.Error(err)
			}
			if entry.Notes != tt.expected_notes {
				t.Errorf("got %s, expected %s", entry.Notes, tt.expected_notes)
			}
		})
	}

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

	entry, err := Get_Entry(db_test, money_entry_3)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(entry)

	if entry.Id != 9 {
		t.Errorf("got %d, expected %d", entry.Id, 9)
	}

	if entry.Notes != "big spendin" {
		t.Errorf("got %s, expected %s", entry.Notes, "big spendin")
	}

	if entry.Logs[0].Present != "2000.00" {
		t.Errorf("got %s, expected %s", entry.Logs[0].Present, "2000.00")
	}

	if entry.Logs[1].Present != "Visa" {
		t.Errorf("got %s, expected %s", entry.Logs[1].Present, "Visa")
	}

}

func Test_Get_Entries(t *testing.T) {
	_test_Reset_Entry_Database(t)

	journal_id, _ := Create_Tracker(db_test, "Journal", "Daily journal and notes")
	Create_Entry(db_test, journal_id, "Entry 1")
	Create_Entry(db_test, journal_id, "Why is green sometimes blue")
	Create_Entry(db_test, journal_id, "If Franky can be a robot maybe I can too")
	Create_Entry(db_test, journal_id, "")
	Create_Entry(db_test, journal_id, "I got lost in a square")
	Create_Entry(db_test, journal_id, "The circle showed me the way")

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
	entries, err := Get_Entries(db_test, money_id)
	if err != nil {
		t.Error(err)
	}

	// s, _ := json.MarshalIndent(entries, "", "    ")
	// fmt.Println("JSON:", string(s))

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

func Test_Get_Entries_Filter(t *testing.T) {
	_test_Reset_Entry_Database(t)

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
	entries, err := Get_Entries_Filter(db_test, money_id, "en")
	if err != nil {
		t.Error(err)
	}

	// s, _ := json.MarshalIndent(entries, "", "    ")
	// fmt.Println("JSON:", string(s))

	if entries[0].Id != 3 {
		t.Errorf("got %d, expected %d", entries[0].Id, 3)
	}

	if entries[1].Notes != "9.99 dollars entered as 999" {
		t.Errorf("got %s, expected %s", entries[1].Notes, "9.99 dollars entered as 999")
	}

	if entries[1].Logs[0].Present != "-9.99" {
		t.Errorf("got %s, expected %s", entries[1].Logs[0].Present, "-9.99")
	}

	if entries[1].Logs[1].Present != "Discover" {
		t.Errorf("got %s, expected %s", entries[1].Logs[1].Present, "Discover")
	}
}

// Update

func Test_Update_Entry_Timestamp(t *testing.T) {
	_test_Reset_Entry_Database(t)

	journal_id, _ := Create_Tracker(db_test, "Journal", "Daily journal and notes")
	Create_Entry(db_test, journal_id, "Entry 1")
	entry_id, _ := Create_Entry(db_test, journal_id, "Why is green sometimes blue")
	Create_Entry(db_test, journal_id, "If Franky can be a robot maybe I can too")

	err := Update_Entry_Timestamp(db_test, entry_id, "2345-03-04 05:06:07")
	if err != nil {
		t.Error(err)
	}

	entry, err := Get_Entry(db_test, entry_id)
	if err != nil {
		t.Error(err)
	}

	if entry.Timestamp != "2345-03-04T05:06:07Z" {
		t.Error("timestamp was not changed", entry.Timestamp)
	}
}

func Test_Update_Entry_Notes(t *testing.T) {
	_test_Reset_Entry_Database(t)

	journal_id, _ := Create_Tracker(db_test, "Journal", "Daily journal and notes")
	Create_Entry(db_test, journal_id, "Entry 1")
	entry_id, _ := Create_Entry(db_test, journal_id, "Why is green sometimes blue")
	Create_Entry(db_test, journal_id, "If Franky can be a robot maybe I can too")

	// Test
	err := Update_Entry_Notes(db_test, entry_id, "Why is green sometimes cyan")
	if err != nil {
		t.Error(err)
	}

	entry, err := Get_Entry(db_test, entry_id)
	if err != nil {
		t.Error(err)
	}

	if entry.Notes != "Why is green sometimes cyan" {
		t.Error("notes were not changed", entry.Notes)
	}
}

func Test_Update_Log(t *testing.T) {
	_test_Reset_Entry_Database(t)

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
	log_amount_id, _ := Add_Log_To_Entry(db_test, money_entry_2, money_amount_id, -42069)
	log_card_id, _ := Add_Log_To_Entry(db_test, money_entry_2, money_card_id, 3)
	money_entry_3, _ := Create_Entry(db_test, money_id, "big spendin")
	Add_Log_To_Entry(db_test, money_entry_3, money_amount_id, 2000_00)
	Add_Log_To_Entry(db_test, money_entry_3, money_card_id, 2)

	err := Update_Log(db_test, log_amount_id, -42169)
	if err != nil {
		t.Error(err)
	}

	err = Update_Log(db_test, log_card_id, 2)
	if err != nil {
		t.Error(err)
	}

	entry, err := Get_Entry(db_test, money_entry_2)
	if err != nil {
		t.Error(err)
	}

	if entry.Logs[0].Value != -42169 {
		t.Error("value was not changed", entry.Logs[0].Value)
	}

	if entry.Logs[1].Value != 2 {
		t.Error("value was not changed", entry.Logs[1].Value)
	}
}

// Delete

func Test_Delete_Entry(t *testing.T) {
	_test_Reset_Entry_Database(t)

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

	err := Delete_Entry(db_test, money_entry_2)
	if err != nil {
		t.Error(err)
	}

	entries, err := Get_Entries(db_test, money_id)
	if err != nil {
		t.Error(err)
	}

	if len(entries) != 2 {
		t.Error("entry not deleted", len(entries))
	}

	err = Delete_Entry(db_test, money_entry_3)
	if err != nil {
		t.Error(err)
	}

	money_entry_4, _ := Create_Entry(db_test, money_id, "not for what you think")
	Add_Log_To_Entry(db_test, money_entry_4, money_amount_id, -42069)
	Add_Log_To_Entry(db_test, money_entry_4, money_card_id, 3)

	entries, err = Get_Entries(db_test, money_id)
	if err != nil {
		t.Error(err)
	}

	if len(entries) != 2 {
		t.Error("entry not deleted", len(entries))
	}
}
