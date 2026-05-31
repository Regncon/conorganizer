package eventservice

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
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

	currentPuljeValue := r.URL.Query().Get("pulje")
	if currentPuljeValue == "" {
		return components.PreviousNext{}, nil
	}

	currentPuljeID, ok := models.ParsePulje(currentPuljeValue)
	if !ok {
		return components.PreviousNext{}, nil
	}

	occurrences, err := root.GetPublishedEventOccurrences(db)
	if err != nil {
		return components.PreviousNext{}, err
	}

	items := make([]rootEventNavigationItem, 0, len(occurrences))
	for _, occurrence := range occurrences {
		items = append(items, rootEventNavigationItem{
			eventID: occurrence.Event.Id,
			title:   occurrence.Event.Title,
			url:     fmt.Sprintf("/event/%s?pulje=%s", occurrence.Event.Id, occurrence.PuljeID),
			puljeID: occurrence.PuljeID,
		})
	}

	return previousNextFromRootEventNavigationItems(items, func(item rootEventNavigationItem) bool {
		return item.eventID == currentID && item.puljeID == currentPuljeID
	}, eventImageDir), nil
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
