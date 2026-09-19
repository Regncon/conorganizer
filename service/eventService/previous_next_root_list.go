package eventservice

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Regncon/conorganizer/components"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/pages/root"
	"github.com/Regncon/conorganizer/service/eventimage"
)

type rootEventNavigationItem struct {
	eventID string
	title   string
	url     string
	puljeID models.Pulje
	program bool
}

func GetPreviousNextForRootEventList(
	ctx context.Context,
	db *sql.DB,
	currentID string,
	programPublished bool,
	r *http.Request,
	eventImageDir *string,
) (components.PreviousNext, error) {
	if programPublished {
		return getPreviousNextForPublishedRootEventList(ctx, db, currentID, r, eventImageDir)
	}

	return getPreviousNextForAnnouncedRootEventList(ctx, db, currentID, eventImageDir)
}

func getPreviousNextForAnnouncedRootEventList(
	ctx context.Context,
	db *sql.DB,
	currentID string,
	eventImageDir *string,
) (components.PreviousNext, error) {
	_ = ctx

	events, err := root.GetAnnouncedEventsAlphabetically(db)
	if err != nil {
		return components.PreviousNext{}, err
	}

	items := make([]rootEventNavigationItem, 0, len(events))
	for _, event := range events {
		items = append(items, rootEventNavigationItem{
			eventID: event.Id,
			title:   event.Title,
			url:     fmt.Sprintf("/event/%s", event.Id),
		})
	}

	return previousNextFromRootEventNavigationItems(items, func(item rootEventNavigationItem) bool {
		return item.eventID == currentID
	}, eventImageDir), nil
}

func getPreviousNextForPublishedRootEventList(
	ctx context.Context,
	db *sql.DB,
	currentID string,
	r *http.Request,
	eventImageDir *string,
) (components.PreviousNext, error) {
	_ = ctx

	query := r.URL.Query()
	currentPuljeValue := query.Get("pulje")
	if currentPuljeValue == "" {
		return components.PreviousNext{}, nil
	}

	currentPuljeID, ok := models.ParsePulje(currentPuljeValue)
	if !ok {
		return components.PreviousNext{}, nil
	}

	programDays, err := root.GetPublishedProgramDays(db)
	if err != nil {
		return components.PreviousNext{}, err
	}

	requestedDate := query.Get("date")
	if requestedDate == "" {
		requestedDate = dateForPulje(programDays, currentPuljeID)
	}

	var selectedDay *root.ProgramDay
	for index := range programDays {
		if programDays[index].QueryValue() == requestedDate {
			selectedDay = &programDays[index]
			break
		}
	}
	if selectedDay == nil {
		return components.PreviousNext{}, nil
	}

	items := make([]rootEventNavigationItem, 0)
	for _, programEvent := range selectedDay.ProgramEvents {
		items = append(items, rootEventNavigationItem{
			eventID: programEvent.Event.Id,
			title:   programEvent.Event.Title,
			url:     rootEventURL(programEvent.Event.Id, selectedDay.QueryValue(), programEvent.PuljeID),
			puljeID: programEvent.PuljeID,
			program: true,
		})
	}
	for _, block := range selectedDay.Blocks {
		for _, event := range root.RaffleEvents(block) {
			items = append(items, rootEventNavigationItem{
				eventID: event.Id,
				title:   event.Title,
				url:     rootEventURL(event.Id, selectedDay.QueryValue(), block.Pulje.ID),
				puljeID: block.Pulje.ID,
			})
		}
	}

	return previousNextFromRootEventNavigationItems(items, func(item rootEventNavigationItem) bool {
		if item.eventID != currentID {
			return false
		}
		if item.program {
			return true
		}
		return item.puljeID == currentPuljeID
	}, eventImageDir), nil
}

func dateForPulje(days []root.ProgramDay, puljeID models.Pulje) string {
	for _, day := range days {
		for _, block := range day.Blocks {
			if block.Pulje.ID == puljeID {
				return day.QueryValue()
			}
		}
	}
	return ""
}

func rootEventURL(eventID string, date string, puljeID models.Pulje) string {
	query := url.Values{}
	query.Set("date", date)
	query.Set("pulje", string(puljeID))
	return fmt.Sprintf("/event/%s?%s", eventID, query.Encode())
}

func previousNextFromRootEventNavigationItems(
	items []rootEventNavigationItem,
	matchesCurrent func(rootEventNavigationItem) bool,
	eventImageDir *string,
) components.PreviousNext {
	currentIndex := -1
	for i, item := range items {
		if matchesCurrent(item) {
			currentIndex = i
			break
		}
	}
	if currentIndex == -1 {
		return components.PreviousNext{}
	}

	var previousItem, nextItem *rootEventNavigationItem
	if currentIndex > 0 {
		previousItem = &items[currentIndex-1]
	}
	if currentIndex < len(items)-1 {
		nextItem = &items[currentIndex+1]
	}

	var previousURL, previousTitle, previousImageURL string
	if previousItem != nil {
		previousURL = previousItem.url
		previousTitle = previousItem.title
		previousImageURL = rootEventNavigationImageURL(previousItem.eventID, eventImageDir)
	}

	var nextURL, nextTitle, nextImageURL string
	if nextItem != nil {
		nextURL = nextItem.url
		nextTitle = nextItem.title
		nextImageURL = rootEventNavigationImageURL(nextItem.eventID, eventImageDir)
	}

	return components.PreviousNext{
		PreviousUrl:      previousURL,
		PreviousTitle:    previousTitle,
		PreviousImageURL: previousImageURL,
		NextUrl:          nextURL,
		NextTitle:        nextTitle,
		NextImageURL:     nextImageURL,
	}
}

func rootEventNavigationImageURL(eventID string, eventImageDir *string) string {
	imageURL := eventimage.GetEventImageUrl(eventID, "banner", eventImageDir)
	if strings.Contains(imageURL, "placeholder") {
		return ""
	}
	return imageURL
}
