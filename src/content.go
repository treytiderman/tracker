package main

import (
	"log/slog"
	"os"
)

type Content struct {
	Path string
	Name string
	Size int
}

var content_folder_path = "../content/"

func Get_Content_List() ([]Content, error) {
	content_list := make([]Content, 0)

	files, err := os.ReadDir(content_folder_path)
	if err != nil {
		return content_list, err
	}

	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return content_list, err
		}

		content_list = append(content_list, Content{
			Path: content_folder_path + file.Name(),
			Name: file.Name(),
			Size: int(info.Size()),
		})
	}

	return content_list, nil
}

func Get_Content(file_name string) ([]byte, error) {
	file, err := os.ReadFile(content_folder_path + file_name)
	if err != nil {
		return []byte(""), err
	}

	return file, nil
}

func Upload_Content(file_name string, file_bytes []byte) error {
	file, err := os.Create(content_folder_path + file_name)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(file_bytes)
	if err != nil {
		return err
	}
	slog.Info("file uploaded", "path", content_folder_path+file_name)

	// add content path to db?

	return nil
}

func Delete_Content(file_name string) error {
	err := os.Remove(content_folder_path + file_name)
	if err != nil {
		return err
	}
	slog.Info("file deleted", "path", content_folder_path+file_name)

	// remove content path from db?

	return nil
}
