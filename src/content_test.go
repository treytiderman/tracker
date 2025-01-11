package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_Get_Content_List(t *testing.T) {
	content_list, err := Get_Content_List()
	if err != nil {
		t.Error(err)
	}

	s, _ := json.MarshalIndent(content_list, "", "    ")
	fmt.Println("JSON:", string(s))
}

func Test_Content(t *testing.T) {
	var err error

	err = Upload_Content("test.txt", []byte("Hello File System!"))
	if err != nil {
		t.Error(err)
	}

	content, err := Get_Content("test.txt")
	if err != nil {
		t.Error(err)
	}

	err = Upload_Content("test2.txt", content)
	if err != nil {
		t.Error(err)
	}

	err = Delete_Content("test.txt")
	if err != nil {
		t.Error(err)
	}

	err = Delete_Content("test2.txt")
	if err != nil {
		t.Error(err)
	}
}

