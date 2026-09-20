package eventservice

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/Regncon/conorganizer/components"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/eventimage"
	"github.com/Regncon/conorganizer/service/program"
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

	events, err := program.GetAnnouncedEvents(db)
	if err != nil {
		return components.PreviousNext{}, err
	}

	items := make([]rootEventNavigationItem, 0, len(events))
	for _, event := range events {
		items = append(items, rootEventNavigationItem{
			eventID: event.Id,
			title:   event.Title,
			url:     program.EventURL(event.Id, "", ""),
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

	programDays, err := program.GetDays(db)
	if err != nil {
		return components.PreviousNext{}, err
	}

	requestedDate := query.Get("date")
	if requestedDate == "" {
		requestedDate = dateForPulje(programDays, currentPuljeID)
	}

	var selectedDay *program.Day
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
			url:     program.EventURL(programEvent.Event.Id, string(programEvent.PuljeID), selectedDay.QueryValue()),
			puljeID: programEvent.PuljeID,
			program: true,
		})
	}
	for _, block := range selectedDay.Blocks {
		for _, event := range block.RaffleEvents() {
			items = append(items, rootEventNavigationItem{
				eventID: event.Id,
				title:   event.Title,
				url:     program.EventURL(event.Id, string(block.Pulje.ID), selectedDay.QueryValue()),
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

func dateForPulje(days []program.Day, puljeID models.Pulje) string {
	for _, day := range days {
		for _, block := range day.Blocks {
			if block.Pulje.ID == puljeID {
				return day.QueryValue()
			}
		}
	}
	return ""
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
