package main

import (
	"fmt"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/phanorcoll/urlshortener/routes"
	"os"
)

var (
	port = os.Getenv("PORT")
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	e := echo.New()
	// e.GET("/:url", routes.ResolveURL)
	e.POST("/:url", routes.ShortenUrl)

	e.Logger.Fatal(e.Start(":" + port))
}
