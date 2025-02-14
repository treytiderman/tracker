package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	_ "modernc.org/sqlite"
)

var db_test *sql.DB

func _test_Reset_Tracker_Database(t *testing.T) {
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
}

func _test_Create_Tracker_Journal(t *testing.T) {
	_, err := Create_Tracker(db_test, "Journal", "Daily journal and notes")
	if err != nil {
		t.Error(err)
	}
}

func _test_Create_Tracker_Weight(t *testing.T) {
	tracker_id, err := Create_Tracker(db_test, "Weight", "Body weight over time")
	if err != nil {
		t.Error(err)
	}

	_, err = Add_Number_Field(db_test, tracker_id, "Weight", "Body weight in pounds (lbs)", 0)
	if err != nil {
		t.Error(err)
	}
}

func _test_Create_Tracker_Money(t *testing.T) {
	tracker_id, err := Create_Tracker(db_test, "Money", "Transactions")
	if err != nil {
		t.Error(err)
	}

	_, err = Add_Number_Field(db_test, tracker_id, "Amount", "Amount of money in dollars", 2)
	if err != nil {
		t.Error(err)
	}

	field_id, err := Add_Option_Field(db_test, tracker_id, "Card", "Payment Method")
	if err != nil {
		t.Error(err)
	}

	_, err = Add_Option_to_Field(db_test, field_id, 1, "Discover")
	if err != nil {
		t.Error(err)
	}

	_, err = Add_Option_to_Field(db_test, field_id, 2, "Visa")
	if err != nil {
		t.Error(err)
	}

	_, err = Add_Option_to_Field(db_test, field_id, 3, "American Express")
	if err != nil {
		t.Error(err)
	}
}

// Create

func Test_Create_Tracker(t *testing.T) {
	_test_Reset_Tracker_Database(t)

	var tests = []struct {
		expected_id   int
		tracker_name  string
		tracker_notes string
	}{
		{1, "Journal", "Daily journal and notes"},
		{2, "Weight", "How much do I weight in pounds"},
		{3, "Money", "Transactions"},
		{4, "DELETE_ME_4", "Tracker made just to delete"},
		{5, "Brush Teeth", "Brush Teeth for 2 Minutes"},
		{6, "AC Filter", "Replace every 3 months, Size: 14x25x1"},
		{7, "Bathroom 💩", ""},
		{8, "Workout", "Complete various exercises"},
	}

	for _, tt := range tests {
		t.Run(tt.tracker_name, func(t *testing.T) {
			id, err := Create_Tracker(db_test, tt.tracker_name, tt.tracker_notes)
			if err != nil {
				t.Error("got error", err)
			}
			if id != tt.expected_id {
				t.Errorf("got %d, expected %d", id, tt.expected_id)
			}
		})
	}

	_, err := Create_Tracker(db_test, tests[0].tracker_name, tests[0].tracker_notes)
	if err == nil {
		t.Error("duplicate names should error")
	}
}

func Test_Add_Number_Field(t *testing.T) {
	_test_Reset_Tracker_Database(t)
	_test_Create_Tracker_Journal(t)
	_test_Create_Tracker_Weight(t)

	var tests = []struct {
		expected_id    int
		tracker_id     int
		field_name     string
		field_notes    string
		decimal_places int
	}{
		{2, 1, "Weight", "Body Weight in pounds (lbs)", 0},
		{3, 1, "Transaction", "Money Spent", 2},
		{4, 1, "num", "xxx", -1},
		{5, 2, "Reps", "", 0},
	}
	for _, tt := range tests {
		t.Run(tt.field_name, func(t *testing.T) {
			id, err := Add_Number_Field(db_test, tt.tracker_id, tt.field_name, tt.field_notes, tt.decimal_places)
			if err != nil {
				t.Error(err)
			}
			if id != tt.expected_id {
				t.Errorf("got %d, expected %d", id, tt.expected_id)
			}
		})
	}
}

func Test_Add_Option_Field(t *testing.T) {
	_test_Reset_Tracker_Database(t)
	_test_Create_Tracker_Journal(t)

	var tests = []struct {
		expected_id int
		tracker_id  int
		field_name  string
		field_notes string
	}{
		{1, 1, "Card", "Which Credit or Debit Card was used?"},
		{2, 1, "W-L", "Win Lose or Draw"},
		{3, 1, "Exercise", ""},
		{4, 1, "option", "xxx"},
	}
	for _, tt := range tests {
		t.Run(tt.field_name, func(t *testing.T) {
			id, err := Add_Option_Field(db_test, tt.tracker_id, tt.field_name, tt.field_notes)
			if err != nil {
				t.Error(err)
			}
			if id != tt.expected_id {
				t.Errorf("got %d, expected %d", id, tt.expected_id)
			}
		})
	}
}

func Test_Add_Option_To_Field(t *testing.T) {
	_test_Reset_Tracker_Database(t)
	_test_Create_Tracker_Journal(t)

	field_id, err := Add_Option_Field(db_test, 1, "Card", "Which Credit or Debit Card was used?")
	if err != nil {
		t.Error(err)
	}

	var tests = []struct {
		expected_id  int
		option_value int
		option_name  string
	}{
		{1, -1, "Discover"},
		{2, 0, "Visa"},
		{3, 2, "American Express"},
	}
	for _, tt := range tests {
		t.Run(tt.option_name, func(t *testing.T) {
			id, err := Add_Option_to_Field(db_test, field_id, tt.option_value, tt.option_name)
			if err != nil {
				t.Error(err)
			}
			if id != tt.expected_id {
				t.Errorf("got %d, expected %d", id, tt.expected_id)
			}
		})
	}
}

func Test_Add_Option_Field_With_Options(t *testing.T) {
	_test_Reset_Tracker_Database(t)
	_test_Create_Tracker_Journal(t)

	var tests = []struct {
		expected_id int
		tracker_id  int
		field_name  string
		field_notes string
		options     []struct {
			Value int
			Name  string
		}
	}{
		{1, 1, "Card", "Which Credit or Debit Card was used?", []struct {
			Value int
			Name  string
		}{
			{-1, "Discover"},
			{0, "Visa"},
			{2, "American Express"},
		}},
		{2, 1, "W-L", "Win Lose or Draw", []struct {
			Value int
			Name  string
		}{
			{1, "Win"},
			{0, "Draw"},
			{-1, "Lose"},
		}},
		{3, 1, "Exercise", "", []struct {
			Value int
			Name  string
		}{
			{2, "Squat"},
			{3, "Deadlift"},
			{1, "Bench Press"},
		}},
		{4, 1, "option", "xxx", []struct {
			Value int
			Name  string
		}{
			{1, "option 1"},
			{2, "option 2"},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.field_name, func(t *testing.T) {
			id, err := Add_Option_Field_With_Options(db_test, tt.tracker_id, tt.field_name, tt.field_notes, tt.options)
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

func Test_Get_Trackers(t *testing.T) {
	_test_Reset_Tracker_Database(t)
	_test_Create_Tracker_Journal(t)
	_test_Create_Tracker_Weight(t)
	_test_Create_Tracker_Money(t)

	trackers, err := Get_Trackers(db_test)
	if err != nil {
		t.Error(err)
	}
	if trackers[0].Id != 1 || trackers[0].Name != "Journal" {
		t.Errorf("got %s, expected %s", trackers[0].Name, "Journal")
	}
	if trackers[1].Id != 2 || trackers[1].Name != "Weight" {
		t.Errorf("got %s, expected %s", trackers[1].Name, "Weight")
	}
	if trackers[2].Id != 3 || trackers[2].Name != "Money" {
		t.Errorf("got %s, expected %s", trackers[2].Name, "Money")
	}
	if trackers[2].Fields[1].Options[0].Name != "Discover" {
		t.Errorf("got %s, expected %s", trackers[2].Fields[1].Options[0].Name, "Discover")
		j, _ := json.MarshalIndent(trackers, "", "    ")
		fmt.Println("JSON:", string(j))
	}
}

func Test_Get_Tracker_By_Id(t *testing.T) {
	_test_Reset_Tracker_Database(t)
	_test_Create_Tracker_Money(t)

	test_json := `{
    "Id": 1,
    "Name": "Money",
    "Notes": "Transactions",
    "Fields": [
        {
            "Id": 1,
            "Type": "number",
            "Name": "Amount",
            "Notes": "Amount of money in dollars",
            "Number": {
                "Id": 1,
                "Decimal_Places": 2
            },
            "Options": null
        },
        {
            "Id": 2,
            "Type": "option",
            "Name": "Card",
            "Notes": "Payment Method",
            "Number": {
                "Id": 0,
                "Decimal_Places": 0
            },
            "Options": [
                {
                    "Id": 1,
                    "Value": 1,
                    "Name": "Discover"
                },
                {
                    "Id": 2,
                    "Value": 2,
                    "Name": "Visa"
                },
                {
                    "Id": 3,
                    "Value": 3,
                    "Name": "American Express"
                }
            ]
        }
    ]
}`
	tracker, err := Get_Tracker(db_test, 1)
	if err != nil {
		t.Error(err)
	}

	j, _ := json.MarshalIndent(tracker, "", "    ")
	tracker_json := string(j)

	if test_json != tracker_json {
		t.Error("test_json doesn't match tracker_json")
		fmt.Println("test_json:", test_json)
		fmt.Println("tracker_json:", tracker_json)
	}
}

func Test_Get_Tracker_Id_By_Name(t *testing.T) {
	_test_Reset_Tracker_Database(t)
	_test_Create_Tracker_Journal(t)
	_test_Create_Tracker_Weight(t)
	_test_Create_Tracker_Money(t)

	var tests = []struct {
		expected_id  int
		tracker_name string
	}{
		{1, "Journal"},
		{2, "Weight"},
		{3, "Money"},
	}
	for _, tt := range tests {
		t.Run(tt.tracker_name, func(t *testing.T) {
			id, err := Get_Tracker_Id_By_Name(db_test, tt.tracker_name)
			if err != nil {
				t.Error(err)
			}
			if id != tt.expected_id {
				t.Errorf("got %d, expected %d", id, tt.expected_id)
			}
		})
	}
}

// Update

func Test_Update_Tracker_Name(t *testing.T) {
	_test_Reset_Tracker_Database(t)
	_test_Create_Tracker_Journal(t)

	err := Update_Tracker_Name(db_test, 1, "Notes")
	if err != nil {
		t.Error(err)
	}

	tracker, err := Get_Tracker(db_test, 1)
	if err != nil {
		t.Error(err)
	}
	if tracker.Name != "Notes" {
		t.Errorf("got %s, expected %s", tracker.Name, "Notes")
	}

	err = Update_Tracker_Name(db_test, 1, "Journal")
	if err != nil {
		t.Error(err)
	}

	tracker, err = Get_Tracker(db_test, 1)
	if err != nil {
		t.Error(err)
	}
	if tracker.Name != "Journal" {
		t.Errorf("got %s, expected %s", tracker.Name, "Journal")
	}
}

func Test_Update_Tracker_Notes(t *testing.T) {
	_test_Reset_Tracker_Database(t)
	_test_Create_Tracker_Journal(t)

	err := Update_Tracker_Notes(db_test, 1, "Some notes")
	if err != nil {
		t.Error(err)
	}

	tracker, err := Get_Tracker(db_test, 1)
	if err != nil {
		t.Error(err)
	}
	if tracker.Notes != "Some notes" {
		t.Errorf("got %s, expected %s", tracker.Notes, "Some notes")
	}
}

func Test_Update_Field_Name(t *testing.T) {
	_test_Reset_Tracker_Database(t)
	_test_Create_Tracker_Money(t)

	err := Update_Field_Name(db_test, 1, "Dollar Amount")
	if err != nil {
		t.Error(err)
	}

	err = Update_Field_Name(db_test, 2, "Payment Method")
	if err != nil {
		t.Error(err)
	}

	tracker, err := Get_Tracker(db_test, 1)
	if err != nil {
		t.Error(err)
	}
	if tracker.Fields[0].Name != "Dollar Amount" {
		t.Errorf("got %s, expected %s", tracker.Fields[0].Name, "Dollar Amount")
	}
	if tracker.Fields[1].Name != "Payment Method" {
		t.Errorf("got %s, expected %s", tracker.Fields[0].Name, "Payment Method")
	}
}

func Test_Update_Field_Notes(t *testing.T) {
	_test_Reset_Tracker_Database(t)
	_test_Create_Tracker_Money(t)

	err := Update_Field_Notes(db_test, 1, "monies")
	if err != nil {
		t.Error(err)
	}

	err = Update_Field_Notes(db_test, 2, "card")
	if err != nil {
		t.Error(err)
	}

	tracker, err := Get_Tracker(db_test, 1)
	if err != nil {
		t.Error(err)
	}
	if tracker.Fields[0].Notes != "monies" {
		t.Errorf("got %s, expected %s", tracker.Fields[0].Notes, "monies")
	}
	if tracker.Fields[1].Notes != "card" {
		t.Errorf("got %s, expected %s", tracker.Fields[0].Notes, "card")
	}
}

// Update - Effects Logged Data

func Test_Update_Number_Decimal_Places(t *testing.T) {
	_test_Reset_Entry_Database(t)

	// Create "Money" tracker, add fields, and add entries
	money_id, _ := Create_Tracker(db_test, "Money", "Transactions")
	money_amount_id, _ := Add_Number_Field(db_test, money_id, "Amount", "Amount of money in dollars", 2)
	money_card_id, _ := Add_Option_Field(db_test, money_id, "Card", "Payment Method")
	Add_Option_to_Field(db_test, money_card_id, 1, "Discover")
	Add_Option_to_Field(db_test, money_card_id, 2, "Visa")
	Add_Option_to_Field(db_test, money_card_id, 3, "American Express")
	money_entry_1, _ := Create_Entry(db_test, money_id, "9.99 dollars entered as 999")
	Add_Log_To_Entry(db_test, money_entry_1, money_amount_id, -9_99)
	Add_Log_To_Entry(db_test, money_entry_1, money_card_id, 1)
	money_entry_2, _ := Create_Entry(db_test, money_id, "not for what you think")
	log_id, _ := Add_Log_To_Entry(db_test, money_entry_2, money_amount_id, -420_69)
	Add_Log_To_Entry(db_test, money_entry_2, money_card_id, 3)
	money_entry_3, _ := Create_Entry(db_test, money_id, "big spendin")
	Add_Log_To_Entry(db_test, money_entry_3, money_amount_id, 2000_00)
	Add_Log_To_Entry(db_test, money_entry_3, money_card_id, 2)

	// Test
	err := Update_Number_Decimal_Places(db_test, money_amount_id, 0)
	if err != nil {
		t.Error(err)
	}

	log, _ := Get_Log(db_test, log_id)

	if log.Value != -421 {
		t.Errorf("got %d, expected %d", log.Value, -421)
	}
	if log.Present != "-421" {
		t.Errorf("got %s, expected %s", log.Present, "-421")
	}

	err = Update_Number_Decimal_Places(db_test, money_amount_id, 2)
	if err != nil {
		t.Error(err)
	}

	log, _ = Get_Log(db_test, log_id)

	if log.Value != -42100 {
		t.Errorf("got %d, expected %d", log.Value, -42100)
	}
	if log.Present != "-421.00" {
		t.Errorf("got %s, expected %s", log.Present, "-421.00")
	}

	err = Update_Number_Decimal_Places(db_test, money_amount_id, -2)
	if err != nil {
		t.Error(err)
	}

	log, _ = Get_Log(db_test, log_id)

	if log.Value != -4 {
		t.Errorf("got %d, expected %d", log.Value, -4)
	}
	if log.Present != "-400" {
		t.Errorf("got %s, expected %s", log.Present, "-400")
	}
}

func Test_Update_Option_Name(t *testing.T) {
	_test_Reset_Entry_Database(t)

	// Create "Money" tracker, add fields, and add entries
	money_id, _ := Create_Tracker(db_test, "Money", "Transactions")
	money_amount_id, _ := Add_Number_Field(db_test, money_id, "Amount", "Amount of money in dollars", 2)
	money_card_id, _ := Add_Option_Field(db_test, money_id, "Card", "Payment Method")
	Add_Option_to_Field(db_test, money_card_id, 1, "Discover")
	Add_Option_to_Field(db_test, money_card_id, 2, "Visa")
	option_id_3, _ := Add_Option_to_Field(db_test, money_card_id, 3, "American Express")
	money_entry_1, _ := Create_Entry(db_test, money_id, "9.99 dollars entered as 999")
	Add_Log_To_Entry(db_test, money_entry_1, money_amount_id, -9_99)
	Add_Log_To_Entry(db_test, money_entry_1, money_card_id, 1)
	money_entry_2, _ := Create_Entry(db_test, money_id, "not for what you think")
	Add_Log_To_Entry(db_test, money_entry_2, money_amount_id, -420_69)
	log_id, _ := Add_Log_To_Entry(db_test, money_entry_2, money_card_id, 3)
	money_entry_3, _ := Create_Entry(db_test, money_id, "big spendin")
	Add_Log_To_Entry(db_test, money_entry_3, money_amount_id, 2000_00)
	Add_Log_To_Entry(db_test, money_entry_3, money_card_id, 2)

	// Test
	err := Update_Option_Name(db_test, option_id_3, "Bitcoin")
	if err != nil {
		t.Error(err)
	}

	log, _ := Get_Log(db_test, log_id)

	if log.Value != 3 {
		t.Errorf("got %d, expected %d", log.Value, 3)
	}
	if log.Present != "Bitcoin" {
		t.Errorf("got %s, expected %s", log.Present, "Bitcoin")
	}

}

func Test_Update_Option_Value(t *testing.T) {
	_test_Reset_Entry_Database(t)

	// Create "Money" tracker, add fields, and add entries
	money_id, _ := Create_Tracker(db_test, "Money", "Transactions")
	money_amount_id, _ := Add_Number_Field(db_test, money_id, "Amount", "Amount of money in dollars", 2)
	money_card_id, _ := Add_Option_Field(db_test, money_id, "Card", "Payment Method")
	Add_Option_to_Field(db_test, money_card_id, 1, "Discover")
	Add_Option_to_Field(db_test, money_card_id, 2, "Visa")
	option_id_3, _ := Add_Option_to_Field(db_test, money_card_id, 3, "American Express")
	money_entry_1, _ := Create_Entry(db_test, money_id, "9.99 dollars entered as 999")
	Add_Log_To_Entry(db_test, money_entry_1, money_amount_id, -9_99)
	Add_Log_To_Entry(db_test, money_entry_1, money_card_id, 1)
	money_entry_2, _ := Create_Entry(db_test, money_id, "not for what you think")
	Add_Log_To_Entry(db_test, money_entry_2, money_amount_id, -420_69)
	log_id, _ := Add_Log_To_Entry(db_test, money_entry_2, money_card_id, 3)
	money_entry_3, _ := Create_Entry(db_test, money_id, "big spendin")
	Add_Log_To_Entry(db_test, money_entry_3, money_amount_id, 2000_00)
	Add_Log_To_Entry(db_test, money_entry_3, money_card_id, 2)

	// Test
	err := Update_Option_Value(db_test, option_id_3, 4)
	if err != nil {
		t.Error(err)
	}

	log, _ := Get_Log(db_test, log_id)

	if log.Value != 4 {
		t.Errorf("got %d, expected %d", log.Value, 4)
	}
	if log.Present != "American Express" {
		t.Errorf("got %s, expected %s", log.Present, "American Express")
	}

}

// Delete

func Test_Delete_Tracker(t *testing.T) {
	_test_Reset_Tracker_Database(t)
	_test_Create_Tracker_Weight(t)
	_test_Create_Tracker_Money(t)

	err := Delete_Tracker(db_test, 1)
	if err != nil {
		t.Error(err)
	}

	trackers, err := Get_Trackers(db_test)
	if err != nil {
		t.Error(err)
	}
	if len(trackers) != 1 {
		t.Errorf("got %d, expected %d", len(trackers), 1)
	}
}

func Test_Delete_Field(t *testing.T) {
	_test_Reset_Entry_Database(t)

	// Create "Money" tracker, add fields, and add entries
	money_id, _ := Create_Tracker(db_test, "Money", "Transactions")
	money_amount_id, _ := Add_Number_Field(db_test, money_id, "Amount", "Amount of money in dollars", 2)
	money_card_id, _ := Add_Option_Field(db_test, money_id, "Card", "Payment Method")
	Add_Option_to_Field(db_test, money_card_id, 1, "Discover")
	Add_Option_to_Field(db_test, money_card_id, 2, "Visa")
	Add_Option_to_Field(db_test, money_card_id, 3, "American Express")
	money_entry_1, _ := Create_Entry(db_test, money_id, "9.99 dollars entered as 999")
	Add_Log_To_Entry(db_test, money_entry_1, money_amount_id, -9_99)
	Add_Log_To_Entry(db_test, money_entry_1, money_card_id, 1)
	money_entry_2, _ := Create_Entry(db_test, money_id, "not for what you think")
	Add_Log_To_Entry(db_test, money_entry_2, money_amount_id, -420_69)
	Add_Log_To_Entry(db_test, money_entry_2, money_card_id, 3)
	money_entry_3, _ := Create_Entry(db_test, money_id, "big spendin")
	Add_Log_To_Entry(db_test, money_entry_3, money_amount_id, 2000_00)
	Add_Log_To_Entry(db_test, money_entry_3, money_card_id, 2)

	// Test
    err := Delete_Field(db_test, money_amount_id)
    if err != nil {
        t.Error(err)
    }

    tracker, _ := Get_Tracker(db_test, money_id)
	// s, _ := json.MarshalIndent(tracker, "", "    ")
	// fmt.Println("JSON:", string(s))

    if len(tracker.Fields) != 1 {
        t.Errorf("got %d, expected %d", len(tracker.Fields), 1)
    }
	
    entries, _ := Get_Entries(db_test, money_id)
	// s2, _ := json.MarshalIndent(entries, "", "    ")
	// fmt.Println("JSON:", string(s2))

    if len(entries[0].Logs) != 1 {
        t.Errorf("got %d, expected %d", len(entries[0].Logs), 1)
    }

    err = Delete_Field(db_test, money_card_id)
    if err != nil {
        t.Error(err)
    }

    tracker, _ = Get_Tracker(db_test, money_id)
	// s3, _ := json.MarshalIndent(tracker, "", "    ")
	// fmt.Println("JSON:", string(s3))

    if len(tracker.Fields) != 0 {
        t.Errorf("got %d, expected %d", len(tracker.Fields), 0)
    }
	
    entries, _ = Get_Entries(db_test, money_id)
	// s4, _ := json.MarshalIndent(entries, "", "    ")
	// fmt.Println("JSON:", string(s4))

    if len(entries[0].Logs) != 0 {
        t.Errorf("got %d, expected %d", len(entries[0].Logs), 0)
    }
}

func Test_Delete_Option(t *testing.T) {
	_test_Reset_Entry_Database(t)

	// Create "Money" tracker, add fields, and add entries
	money_id, _ := Create_Tracker(db_test, "Money", "Transactions")
	money_amount_id, _ := Add_Number_Field(db_test, money_id, "Amount", "Amount of money in dollars", 2)
	money_card_id, _ := Add_Option_Field(db_test, money_id, "Card", "Payment Method")
	Add_Option_to_Field(db_test, money_card_id, 1, "Discover")
	option_id_2, _ := Add_Option_to_Field(db_test, money_card_id, 2, "Visa")
	Add_Option_to_Field(db_test, money_card_id, 3, "American Express")
	money_entry_1, _ := Create_Entry(db_test, money_id, "9.99 dollars entered as 999")
	Add_Log_To_Entry(db_test, money_entry_1, money_amount_id, -9_99)
	Add_Log_To_Entry(db_test, money_entry_1, money_card_id, 1)
	money_entry_2, _ := Create_Entry(db_test, money_id, "not for what you think")
	Add_Log_To_Entry(db_test, money_entry_2, money_amount_id, -420_69)
	Add_Log_To_Entry(db_test, money_entry_2, money_card_id, 3)
	money_entry_3, _ := Create_Entry(db_test, money_id, "big spendin")
	Add_Log_To_Entry(db_test, money_entry_3, money_amount_id, 2000_00)
	Add_Log_To_Entry(db_test, money_entry_3, money_card_id, 2)

	// Test
    err := Delete_Option(db_test, option_id_2)
    if err != nil {
        t.Error(err)
    }

    tracker, _ := Get_Tracker(db_test, money_id)
	s, _ := json.MarshalIndent(tracker, "", "    ")
	fmt.Println("JSON:", string(s))

    if len(tracker.Fields[1].Options) != 2 {
        t.Errorf("got %d, expected %d", len(tracker.Fields[1].Options), 2)
    }

    entries, _ := Get_Entries(db_test, money_id)
	s2, _ := json.MarshalIndent(entries, "", "    ")
	fmt.Println("JSON:", string(s2))

    if len(entries) != 2 {
        t.Errorf("got %d, expected %d", len(entries), 2)
    }
}

// Other
