package main

import (
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"os"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello Echo")
	})

  e.Logger.Fatal(e.Start(":"+os.Getenv("PORT")))
}
