package eventservice

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/components"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/pages/root"
	"github.com/Regncon/conorganizer/service/eventimage"
)

func getPreviousPuljeId(currentPuljeId models.Pulje) models.Pulje {
	if currentPuljeId == models.PuljeLordagMorgen {
		return models.PuljeFredagKveld
	}
	if currentPuljeId == models.PuljeLordagKveld {
		return models.PuljeLordagMorgen
	}
	if currentPuljeId == models.PuljeSondagMorgen {
		return models.PuljeLordagKveld
	}
	return models.PuljeFredagKveld
}

func getNextPuljeId(currentPuljeId models.Pulje) models.Pulje {
	if currentPuljeId == models.PuljeFredagKveld {
		return models.PuljeLordagMorgen
	}
	if currentPuljeId == models.PuljeLordagMorgen {
		return models.PuljeLordagKveld
	}
	if currentPuljeId == models.PuljeLordagKveld {
		return models.PuljeSondagMorgen
	}
	return models.PuljeSondagMorgen
}

func getLastEventInPriviousPulje(eventsByPulje root.EventsByPulje, currentPuljeId models.Pulje) *models.EventCardModel {
	var puljeId = getPreviousPuljeId(currentPuljeId)

	if ebp, ok := eventsByPulje[puljeId]; ok {
		if len(ebp.Events) > 0 {
			return &ebp.Events[len(ebp.Events)-1]
		}
	}
	emptyEvent := models.EventCardModel{}
	return &emptyEvent
}

func getFirstEventInPulje(eventsByPulje root.EventsByPulje, currentPuljeId models.Pulje) *models.EventCardModel {
	var puljeId = getNextPuljeId(currentPuljeId)
	if ebp, ok := eventsByPulje[puljeId]; ok {
		if len(ebp.Events) > 0 {
			return &ebp.Events[0]
		}
	}
	emptyEvent := models.EventCardModel{}
	return &emptyEvent
}

func GetPreviousNextByPuljeSimple(
	ctx context.Context,
	eventsByPulje root.EventsByPulje,
	logger *slog.Logger,
	currentID string,
	isAdmin bool,
	r *http.Request,
	eventImageDir *string,
) (components.PreviousNext, error) {
	currentPuljeIdString := r.URL.Query().Get("pulje")

	if currentPuljeIdString == "" {
		return components.PreviousNext{}, nil
	}

	currentPunjeId, ok := models.ParsePulje(currentPuljeIdString)
	if !ok {
		return components.PreviousNext{}, nil
	}

	var previousId, previousTitle, nextId, nextTitl string
	var previousPuljeID, nextPuljeID models.Pulje
	for _, ebp := range eventsByPulje {
		if string(ebp.Pulje.ID) == currentPuljeIdString {
			for i, event := range ebp.Events {
				if event.Id == currentID {
					if i > 0 {
						previousId = ebp.Events[i-1].Id
						previousTitle = ebp.Events[i-1].Title
						previousPuljeID = currentPunjeId
					} else if ebp.Pulje.ID != models.PuljeFredagKveld {
						lastEvent := getLastEventInPriviousPulje(eventsByPulje, currentPunjeId)
						previousId = lastEvent.Id
						previousTitle = lastEvent.Title
						previousPuljeID = getPreviousPuljeId(currentPunjeId)
					}

					if i < len(ebp.Events)-1 {
						nextId = ebp.Events[i+1].Id
						nextTitl = ebp.Events[i+1].Title
						nextPuljeID = currentPunjeId
					} else if ebp.Pulje.ID != models.PuljeSondagMorgen {
						firstEvent := getFirstEventInPulje(eventsByPulje, currentPunjeId)
						nextId = firstEvent.Id
						nextTitl = firstEvent.Title
						nextPuljeID = getNextPuljeId(currentPunjeId)
					}
				}
			}
		}
	}

	prevBanner := ""
	if previousId != "" {
		prevBanner = eventimage.GetEventImageUrl(previousId, "banner", eventImageDir)
	}
	nextBanner := ""
	if nextId != "" {
		nextBanner = eventimage.GetEventImageUrl(nextId, "banner", eventImageDir)
	}
	previousUrl := ""
	if previousId != "" {
		previousUrl = fmt.Sprintf("/event/%s?pulje=%s", previousId, previousPuljeID)
	}
	nextUrl := ""
	if nextId != "" {
		nextUrl = fmt.Sprintf("/event/%s?pulje=%s", nextId, nextPuljeID)
	}
	fmt.Println("previousId:", previousId, "previousTitle:", previousTitle, "nextId:", nextId, "nextTitle:", nextTitl)

	fmt.Println("previous puljeID:", previousPuljeID)
	var result = components.PreviousNext{
		PreviousUrl:      previousUrl,
		PreviousTitle:    previousTitle,
		PreviousImageURL: prevBanner,
		NextUrl:          nextUrl,
		NextTitle:        nextTitl,
		NextImageURL:     nextBanner,
	}
	fmt.Println("PreviousNext Result:", result)
	return result, nil
}
