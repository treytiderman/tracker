package main

import (
	// "database/sql"
	// "encoding/json"
	// "fmt"
	// "os"
	// "testing"

	_ "modernc.org/sqlite"
)

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
