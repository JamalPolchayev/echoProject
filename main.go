package main

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

func main() {
	fmt.Println("Welcome to the server")

	var e = echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(
			http.StatusOK,
			"hello from the web side",
		)
	})
	e.Start(":8000")

}
