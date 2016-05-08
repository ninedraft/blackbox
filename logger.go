// logger.go
package main

import (
	"log"
)

type Logger interface {
	Println(...interface{})
	Printf(string, ...interface{})
}

type DebugLogger struct{}

func (dl *DebugLogger) Println(v ...interface{}) {
	log.Println(v...)
}

func (dl *DebugLogger) Printf(f string, v ...interface{}) {
	log.Printf(f, v...)
}
