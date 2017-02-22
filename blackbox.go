package main

import (
	"fmt"

	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/ninedraft/blackbox/cmd"
	"github.com/ninedraft/blackbox/utils"
)

var (
	GlobalLog = logrus.New()
	logFormat logrus.Formatter
	logLevel  logrus.Level
)

func init() {
	var err error
	if err = cmd.RootCmd.Execute(); err != nil {
		GlobalLog.Fatalf("error while parsing flags: %v", err)
	}
	logFormat, err = utils.ParseLogFormat(cmd.Configuration.LogFormatParam)

	if err != nil {
		GlobalLog.Fatalf("error while parsing flags: %v", err)
	}
	GlobalLog.Formatter = logFormat

	logLevel, err = utils.ParseLogLevel(cmd.Configuration.LogLevelParam)
	if err != nil {
		GlobalLog.Fatalf("error while parsing flags: %v", err)
	}
	GlobalLog.Level = logLevel
}

func main() {
	configuration := cmd.Configuration
	server := echo.New()
	server.Logger.SetLevel(log.Lvl(logLevel))
	server.Use(middleware.Recover())
	server.GET("/", func(ctx echo.Context) error {
		ctx.Logger().Info("hello!")
		return ctx.String(http.StatusOK, fmt.Sprintf("%+v", cmd.Configuration))
	})
	GlobalLog.Fatal(server.Start(configuration.Addr))
}
