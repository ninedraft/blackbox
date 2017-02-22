package utils

import (
	"errors"
	"strings"

	"github.com/Sirupsen/logrus"
)

var (
	ErrInvalidLogFormat = errors.New("invalid log format")
	ErrInvalidLogLevel  = errors.New("invalid log level")
)

func ParseLogFormat(lf string) (logrus.Formatter, error) {
	switch strings.TrimSpace(lf) {
	case "json":
		return &logrus.JSONFormatter{}, nil
	case "text":
		return &logrus.TextFormatter{}, nil
	default:
		return nil, ErrInvalidLogFormat
	}
}

func ParseLogLevel(ll string) (logrus.Level, error) {
	switch strings.TrimSpace(ll) {
	case "debug":
		return logrus.DebugLevel, nil
	case "info":
		return logrus.InfoLevel, nil
	case "fatal":
		return logrus.FatalLevel, nil
	case "panic":
		return logrus.PanicLevel, nil
	default:
		invalidLevel := ^uint8(0)
		return logrus.Level(invalidLevel), ErrInvalidLogLevel
	}
}
