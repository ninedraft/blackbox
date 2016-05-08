// utility.go
package main

import (
	"log"
	"os"
	"path"
)

func Ok(err error) bool {
	if err != nil {
		log.Println(err)
		return false
	} else {
		return true
	}
}

func Fatal(err error) {
	if err != nil {
		log.Println(err)
		//panic(err)
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
