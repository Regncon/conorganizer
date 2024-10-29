package echoServer

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
	"regncon.no/htmx/echoServer/util"
	supabaseSetup "regncon.no/htmx/echoServer/util/supabase"
)

type Test struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

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
		var testRes []Test
		client := supabaseSetup.Client
		count, err := client.From("test").Select("*", "estimated", false).ExecuteTo(&testRes)

		if count == 0 && err != nil {
			fmt.Println("err", err)
		}
		fmt.Println("testRes", testRes)
		res := map[string]interface{}{
			"Name":      "Wyndham",
			"Id":        testRes[0].ID,
			"CreatedAt": testRes[0].CreatedAt,
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
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}
	e.Logger.Fatal(e.Start(port))
}
