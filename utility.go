// utility.go
package main

import (
	"errors"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	ErrIsNotDir = errors.New("is not dir")
)

func Ok(f string, err error) {
	if err != nil {
		Log.Printf(f, err)
		os.Exit(1)
	}
}
func GetFilesIfDir(path string) ([]os.FileInfo, error) {
	dir, err := os.Open(path)
	defer dir.Close()
	if err != nil {
		return nil, err
	}
	info, err := dir.Stat()
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, ErrIsNotDir
	}
	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func FileNameNoExtension(filename string) string {
	return filename[0 : len(filename)-len(path.Ext(filename))]
}

func ExctractTags(text string) []string {
	tags := regexp.MustCompile("#([^#<>\\s\\t\\r\\n]+)").FindAllString(text, -1)
	for i, t := range tags {
		tags[i] = strings.ToLower(t)
	}
	return tags
}
