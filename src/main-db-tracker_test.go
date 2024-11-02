package main

import (
	"database/sql"
	"os"
	// "encoding/json"
	// "fmt"
	"testing"

	_ "modernc.org/sqlite"
)

var has_reset = false
var db_test *sql.DB

func test_setup_function() {
	if has_reset {
		return
	}
	os.Remove("../data/test.db")
	db_test, _ = sql.Open("sqlite", "../data/test.db")
	Db_Tracker_Table_Create(db_test)
	has_reset = true
}

func Test_Create_Db_Tables(t *testing.T) {
	test_setup_function()
}

// Insert Data

func Test_Db_Tracker_Create(t *testing.T) {
	test_setup_function()
	var tests = []struct {
		expected_id   int
		tracker_name  string
		tracker_notes string
	}{
		{1, "Journal", "Daily journal and notes"},
		{2, "Weight", "How much do I weight in pounds"},
		{3, "Money", "Transactions"},
		{4, "DELETE_ME", "Tracker made just to delete"},
		{5, "Brush Teeth", "Brush Teeth for 2 Minutes"},
		{6, "AC Filter", "Replace every 3 months, Size: 14x25x1"},
		{7, "Bathroom", ""},
		{8, "Workout", "Complete various exercises"},
	}
	for _, tt := range tests {
		t.Run(tt.tracker_name, func(t *testing.T) {
			id, err := Db_Tracker_Create(db_test, tt.tracker_name, tt.tracker_notes)
			if err != nil {
				t.Fatalf("got error, expected %d", tt.expected_id)
			}
			if id != tt.expected_id {
				t.Errorf("got %d, expected %d", id, tt.expected_id)
			}
		})
	}
}

// func Test_Add_Number_Field(t *testing.T) {
// 	var tests = []struct {
// 		expected_id    int
// 		tracker_name   string
// 		field_name     string
// 		field_notes    string
// 		decimal_places int
// 	}{
// 		{1, "Weight", "Weight", "in Pounds (lbs)", 0},
// 		{2, "Money", "Transaction", "Money Spent", 2},
// 		{3, "Workout", "Weight", "", 0},
// 		{4, "Workout", "Reps", "", 0},
// 		{5, "DELETE_ME", "num", "xxx", -1},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.tracker_name, func(t *testing.T) {
// 			id, err := Add_Number_Field(db_test, tt.tracker_name, tt.field_name, tt.field_notes, tt.decimal_places)
// 			if err != nil {
// 				t.Fatalf("Failed to add number field to tracker")
// 			}
// 			if id != tt.expected_id {
// 				t.Errorf("got %d, expected %d", id, tt.expected_id)
// 			}
// 		})
// 	}
// }

// func Test_Add_Option_Field(t *testing.T) {
// 	var tests = []struct {
// 		expected_id  int
// 		tracker_name string
// 		field_name   string
// 		field_notes  string
// 		options      []struct {
// 			Value int
// 			Name  string
// 		}
// 	}{
// 		{6, "Money", "Card", "Which Credit or Debit Card was used?", []struct {
// 			Value int
// 			Name  string
// 		}{
// 			{-1, "Discover"},
// 			{0, "Visa"},
// 			{2, "American Express"},
// 		}},
// 		{7, "Workout", "Exercise", "", []struct {
// 			Value int
// 			Name  string
// 		}{
// 			{2, "Squat"},
// 			{3, "Deadlift"},
// 			{1, "Bench Press"},
// 		}},
// 		{8, "DELETE_ME", "opt", "xxx", []struct {
// 			Value int
// 			Name  string
// 		}{
// 			{1, "Option 1"},
// 			{2, "Option 2"},
// 		}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.tracker_name, func(t *testing.T) {
// 			id, err := Add_Option_Field(db_test, tt.tracker_name, tt.field_name, tt.field_notes, tt.options)
// 			if err != nil {
// 				t.Fatalf("Failed to add option field to tracker")
// 			}
// 			if id != tt.expected_id {
// 				t.Errorf("got %d, expected %d", id, tt.expected_id)
// 			}
// 		})
// 	}
// }

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

// // Get Data

// func Test_Get_Tracker_Id_By_Name(t *testing.T) {
// 	var tests = []struct {
// 		expected_id  int
// 		tracker_name string
// 	}{
// 		{1, "Journal"},
// 		{2, "Weight"},
// 		{8, "Workout"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.tracker_name, func(t *testing.T) {
// 			id, err := Get_Tracker_Id_By_Name(db_test, tt.tracker_name)
// 			if err != nil {
// 				t.Fatalf("Failed to Test_Get_Tracker_Id_By_Name")
// 			}
// 			if id != tt.expected_id {
// 				t.Errorf("got %d, expected %d", id, tt.expected_id)
// 			}
// 		})
// 	}
// }

// func Test_Get_Tracker_By_Id(t *testing.T) {
// 	var tests = []struct {
// 		tracker_id int
// 	}{
// 		{3},
// 	}
// 	for _, tt := range tests {
// 		t.Run(fmt.Sprintf("tracker_id=%d", tt.tracker_id), func(t *testing.T) {
// 			tracker, err := Get_Tracker_By_Id(db_test, tt.tracker_id)
// 			if err != nil {
// 				t.Fatalf("Failed to Get_Tracker_By_Id")
// 			}
// 			s, _ := json.Marshal(tracker)
// 			fmt.Sprintln("JSON:", string(s))
// 			// fmt.Println("JSON:", string(s))
// 		})
// 	}
// }

// func Test_Get_Trackers(t *testing.T) {
// 	trackers, err := Get_Trackers(db_test)
// 	if err != nil {
// 		t.Fatalf("Failed to Get_Trackers")
// 	}
// 	s, _ := json.Marshal(trackers)
// 	fmt.Sprintln("JSON:", string(s))
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

// // Update Data

// func Test_Update_Tracker_Notes_By_Id(t *testing.T) {
// 	var tests = []struct {
// 		tracker_id    int
// 		tracker_notes string
// 	}{
// 		{7, "Bathroom breaks"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(fmt.Sprintf("tracker_id=%d", tt.tracker_id), func(t *testing.T) {
// 			err := Update_Tracker_Notes_By_Id(ddb_testb, tt.tracker_id, tt.tracker_notes)
// 			if err != nil {
// 				t.Fatalf("Failed to update tracker notes")
// 			}
// 		})
// 	}
// }

// // Delete Data

// func Test_Delete_Tracker_By_Id(t *testing.T) {
// 	var tests = []struct {
// 		tracker_id int
// 	}{
// 		{7},
// 	}
// 	for _, tt := range tests {
// 		t.Run(fmt.Sprintf("tracker_id=%d", tt.tracker_id), func(t *testing.T) {
// 			err := Delete_Tracker_By_Id(db_test, tt.tracker_id)
// 			if err != nil {
// 				t.Fatalf("Failed to delete tracker")
// 			}
// 		})
// 	}
// }
