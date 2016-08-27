// logger.go
package main

import (
	"log"
)

type Logger interface {
	Println(...interface{})
	Printf(string, ...interface{})
}

type LoggerDebug struct{}

func (dl *LoggerDebug) Println(v ...interface{}) {
	log.Println(v...)
}

func (dl *LoggerDebug) Printf(f string, v ...interface{}) {
	log.Printf(f, v...)
}

type LoggerVoid struct{}

func (_ *LoggerVoid) Println(...interface{}) {}

func (_ *LoggerVoid) Printf(string, ...interface{}) {}
