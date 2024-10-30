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
	Note      string    `json:"note"`
}

type PoolEvents struct {
	Id                 string `json:"id"`
	AdditionalComments string `json:"additionalComments"`
	AdultsOnly         bool   `json:"adultsOnly"`
	BeginnerFriendly   bool   `json:"beginnerFriendly"`
	BigImageURL        string `json:"bigImageUrl"`
	ChildFriendly      bool   `json:"childFriendly"`
	CreatedAt          string `json:"createdAt"`
	CreatedBy          string `json:"createdBy"`
	Description        string `json:"description"`
	GameMaster         string `json:"gameMaster"`
	GameType           string `json:"gameType"`
	IsSmallCard        bool   `json:"isSmallCard"`
	LessThanThreeHours bool   `json:"lessThanThreeHours"`
	MoreThanSixHours   bool   `json:"moreThanSixHours"`
	ParentEventId      string `json:"parentEventId"`
	Participants       string `json:"participants"`
	PoolName           string `json:"poolName"`
	PossiblyEnglish    bool   `json:"possiblyEnglish"`
	Published          bool   `json:"published"`
	ShortDescription   string `json:"shortDescription"`
	SmallImageURL      string `json:"smallImageUrl"`
	System             string `json:"system"`
	Title              string `json:"title"`
	UpdateAt           string `json:"updateAt"`
	UpdatedBy          string `json:"updatedBy"`
	VolunteersPossible bool   `json:"volunteersPossible"`
}

func EchoServer() {
	e := echo.New()
	client := supabaseSetup.Client

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
		var testRes []PoolEvents
		count, err := client.From("pool-events").Select("*", "estimated", false).ExecuteTo(&testRes)

		if count == 0 && err != nil {
			fmt.Println("err", err)
		}
		fmt.Println("testRes", testRes)
		res := map[string]interface{}{
			"Name":      "Wyndham",
			"PoolEvent": testRes,
			// "Id":        testRes[0].ID,
			// "CreatedAt": testRes[0].CreatedAt,
			// "Note":      testRes[0].Note,
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
	e.GET("/echo/:id", func(c echo.Context) error {
		id := c.Param("id")
		fmt.Println("id", id)
		var testRes []PoolEvents
		count, err := client.From("pool-events").Select("*", "exact", false).Eq("id", id).ExecuteTo(&testRes)

		if count == 0 && err != nil {
			fmt.Println("err", err)
		}

		fmt.Println("testRes", testRes)

		return c.Render(http.StatusOK, "event", testRes[0])
	})
	e.Any("/*", func(e echo.Context) error {
		c := echo.Context(e)
		return c.Render(http.StatusOK, "404", nil)

	})
	port := os.Getenv("PORT")
	if port == "" {
		port = ":3000"
	}
	e.Logger.Fatal(e.Start(port))
}
