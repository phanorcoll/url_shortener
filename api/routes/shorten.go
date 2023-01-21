package routes

import (
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/phanorcoll/urlshortener/helpers"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL            string        `json:"url"`
	CustomShort    string        `json:"short"`
	Expiry         time.Duration `json:"expiry"`
	XRateRemaining int           `json:"rate_limit"`
	RateLimitRest  int           `json:"rate_limit_reset`
}

func ShortenUrl(c echo.Context) error {
	body := new(request)

	if err := c.Bind(body); err != nil {
		return err
	}

	//check if the input is an actual Url
	if !govalidator.IsURL(body.URL) {
		return c.JSON(http.StatusBadRequest, "invalid URL")
	}

	//check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return c.JSON(http.StatusServiceUnavailable, "service unavailable 😁")
	}

	//enforce https, SSL
	body.URL = helpers.EnforceHTTP(body.URL)

	return c.JSON(http.StatusCreated, body)
}
