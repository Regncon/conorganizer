package echoServer

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
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
	Interest           string `json:"interest"`
}

var poolTitlesWithTime = map[string]string{
	"0": "Fredag Kveld Kl 18 - 23",
	"1": "Lørdag Morgen Kl 10 - 15",
	"2": "Lørdag Kveld Kl 18 - 23",
	"3": "Søndag Morgen Kl 10 - 15",
}

func EchoServer() {
	e := echo.New()
	client := supabaseSetup.Client

	// Middleware for housekeeping
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(
		rate.Limit(20),
	)))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Call the next handler in the chain
			err := next(c)

			// If an error occurred
			if err != nil {
				// Log the actual error message
				fmt.Printf("Error: %v\n", err)

				// Return JSON response with the actual error message under "message"
				errorMessage := map[string]interface{}{
					"message": err.Error(), // Use "message" as the key
				}

				return c.JSON(http.StatusInternalServerError, errorMessage)
			}

			return nil
		}
	})
	e.Static("/echo/dist", "echoServer/dist")

	// Template renderer
	util.NewTemplateRenderer(e, "echoServer/public/*.html")

	e.GET("/echo", func(e echo.Context) error {
		c := echo.Context(e)

		// Define the fields that are needed in index.html, including Id
		var eventList []struct {
			Id               string `json:"id"`
			Title            string `json:"title"`
			GameMaster       string `json:"gameMaster"`
			System           string `json:"system"`
			ShortDescription string `json:"shortDescription"`
			BigImageURL      string `json:"bigImageURL"` // Corrected field name to match "bigImageURL"
		}

		// Query only the necessary fields from the database
		count, err := client.From("pool-events").
			Select("id, title, gameMaster, system, shortDescription, bigImageURL", "estimated", false).
			ExecuteTo(&eventList)

		// Check for errors
		if count == 0 && err != nil {
			fmt.Println("Error retrieving pool events:", err)
		} else {
			fmt.Println("Retrieved pool events:", eventList)
		}

		// Prepare the data to render in the template
		res := map[string]interface{}{
			"EventList": eventList,
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

	e.GET("/echo/:id", func(e echo.Context) error {
		c := echo.Context(e)
		id := c.Param("id")
		fmt.Println("id", id)
		var testRes PoolEvents
		count, err := client.From("pool-events").Select("*", "estimated", false).Eq("id", id).Single().ExecuteTo(&testRes)

		if count == 0 && err != nil {
			fmt.Println("err", err)
		}

		// Merge PoolEvent fields and PoolTitlesWithTime into one map
		mergedRes := map[string]interface{}{
			"PoolTitleWithTime": poolTitlesWithTime[testRes.Interest],
		}

		// Use reflection to add PoolEvent fields to mergedRes
		poolEventValue := reflect.ValueOf(testRes)
		poolEventType := reflect.TypeOf(testRes)

		for i := 0; i < poolEventType.NumField(); i++ {
			field := poolEventType.Field(i)
			fieldValue := poolEventValue.Field(i).Interface()
			mergedRes[field.Name] = fieldValue
		}

		return c.Render(http.StatusOK, "event", mergedRes)
	})

	// New endpoint to update the interest value
	e.POST("/echo/update-interest", func(e echo.Context) error {
		c := echo.Context(e)
		id := c.FormValue("id")
		interest := c.FormValue("interest")

		// Log received values for debugging
		fmt.Printf("Received Id: %s, Interest: %s\n", id, interest)

		// Run the database update asynchronously
		go func(id, interest string) {
			interestToUpdate := map[string]interface{}{"interest": interest}

			// Update the interest value in the database
			result, count, err := client.From("pool-events").
				Update(interestToUpdate, "representation", "estimated").
				Eq("id", id).Single().
				Execute()

			resultText := string(result)
			fmt.Printf("Async Update Result: %s\n", resultText)
			fmt.Printf("Rows Affected (Async): %d\n", count)

			if err != nil {
				fmt.Printf("Database update error (Async): %v\n", err)
			}
		}(id, interest) // Pass values to the goroutine

		// Determine PoolTitleWithTime based on Interest level
		poolTitleWithTime, ok := poolTitlesWithTime[interest]
		if !ok {
			poolTitleWithTime = "Unknown Pool Time"
		}

		// Render the interest slider with the known values
		res := map[string]interface{}{
			"Id":                id,
			"PoolTitleWithTime": poolTitleWithTime,
			"Interest":          interest,
		}
		return c.Render(http.StatusOK, "interest_slider", res)
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
