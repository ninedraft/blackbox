package main

import (
	"errors"
	"html/template"
)

var (
	ErrTemplateNotFound = errors.New("template not found")
)

func LoadTemplates(path string) error {
	files, err := GetFilesIfDir(path)
	if err != nil {
		return err
	}
	var t *template.Template
	var tname string
	for _, f := range files {
		tname = FileNameNoExtension(f.Name())
		t, err = template.ParseFiles(path + f.Name())
		Templates[tname] = t
	}
	if Templates["index"] == nil && Templates["home"] == nil {

	}
	return nil
}
