package echoServer

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"regncon.no/htmx/util"
)

func EchoServer() {
	e := echo.New()

	// Little bit of middlewares for housekeeping
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(20),
	)))
	e.Static("/echo/dist", "echoServer/dist")

	// This will initiate our template renderer
	util.NewTemplateRenderer(e, "echoServer/public/*.html")
	e.GET("/echo", func(e echo.Context) error {
		c := echo.Context(e)

		res := map[string]interface{}{
			"Name":  "Wyndham",
			"Phone": "8888888",
			"Email": "skyscraper@gmail.com",
		}
		return c.Render(http.StatusOK, "index", res)
	})

	e.GET("/echo/get-info", func(c echo.Context) error {
		res := map[string]interface{}{
			"Name":  "Wyndham",
			"Phone": "8888888",
			"Email": "skyscraper@gmail.com",
		}
		return c.Render(http.StatusOK, "name_card", res)
	})

	e.Logger.Fatal(e.Start(":3000"))
}
