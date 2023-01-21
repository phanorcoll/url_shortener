package routes

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/phanorcoll/urlshortener/database"
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
	RateLimitRest  int           `json:"rate_limit_reset"`
}

func ShortenUrl(c echo.Context) error {
	body := new(request)

	if err := c.Bind(body); err != nil {
		return c.JSON(http.StatusBadRequest, "cannot parse JSON")
	}

	//implement rate limiting
	r2 := database.CreateClient(1)
	defer r2.Close()
	val, err := r2.Get(database.Ctx, c.RealIP()).Result()
	if err != redis.Nil {
		_ = r2.Set(database.Ctx, c.RealIP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else {
		val, _ = r2.Get(database.Ctx, c.RealIP()).Result()
		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			//TODO:
			// create an object to return, including the rate_limit_rest
			// limit, _ := r2.TTL(database.Ctx, c.RealIP()).Result()
			return c.JSON(http.StatusServiceUnavailable, "Rate limit exceeded")
		}
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

	var id string
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.JSON(http.StatusForbidden, "URL custom short is already in use")
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Unable to connect to server")
	}

	resp := response{
		URL:            body.URL,
		CustomShort:    "",
		Expiry:         body.Expiry,
		XRateRemaining: 10,
		RateLimitRest:  30,
	}

	r2.Decr(database.Ctx, c.RealIP())

	val, _ = r2.Get(database.Ctx, c.RealIP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	// ttl, _ := r2.TTL(database.Ctx, c.RealIP()).Result()
	// resp.RateLimitRest = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return c.JSON(http.StatusOK, resp)
}
