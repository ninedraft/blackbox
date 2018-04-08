package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// Echo instance
	server := echo.New()

	// Middleware
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	// Routes
	server.GET("/", hello)

	// Start server
	server.Logger.Fatal(server.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
