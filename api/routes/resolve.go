package routes

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/phanorcoll/urlshortener/database"
)

func ResolveURL(c echo.Context) error {
	url := c.Param("url")
	r := database.CreateClient(0)
	defer r.Close()

	value, err := r.Get(database.Ctx, url).Result()
	if err != redis.Nil {
		return c.JSON(http.StatusNotFound, "short not found in the database")
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, "cannot connect to db")
	}

	rInr := database.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")
	return c.Redirect(301, value)

}
