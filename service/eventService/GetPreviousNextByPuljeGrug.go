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
	fmt.Println("Current Pulje ID:", currentPunjeId)
	fmt.Println("Events By Pulje:", eventsByPulje)
	/*

				type PuljeBlock struct {
					Pulje  models.PuljeRow
					Events []models.EventCardModel
				}

				type EventsByPulje map[models.Pulje]*PuljeBlock

		Events By Pulje: map[FredagKveld:0xc0002ba3f0 LordagKveld:0xc0002ba310 LordagMorgen:0xc0002ba380 SondagMorgen:0xc0002ba150]
	*/

	var previousId, previousTitle, nextId, nextTitl string

	var previousPuljeID, nextPuljeID models.Pulje
	for _, ebp := range eventsByPulje {
		if string(ebp.Pulje.ID) == currentPuljeIdString {
			fmt.Println("Found matching Pulje:", ebp.Pulje.ID)
			for i, event := range ebp.Events {
				if event.Id == currentID {
					fmt.Println("Found current event at index:", i)
					if i > 0 {
						previousId = ebp.Events[i-1].Id
						previousTitle = ebp.Events[i-1].Title
						previousPuljeID = currentPunjeId
					}
					if i < len(ebp.Events)-1 {
						nextId = ebp.Events[i+1].Id
						nextTitl = ebp.Events[i+1].Title
						nextPuljeID = currentPunjeId
					}
				}
			}
		}
	}

	prevBanner := ""
	nextBanner := ""
	if previousId != "" {
		prevBanner = eventimage.GetEventImageUrl(previousId, "banner", eventImageDir)
	}
	if nextId != "" {
		nextBanner = eventimage.GetEventImageUrl(nextId, "banner", eventImageDir)
	}

	fmt.Println("previous puljeID:", previousPuljeID)
	var previousUrl = fmt.Sprintf("%s?pulje=%s", previousId, previousPuljeID)
	var nextUrl = fmt.Sprintf("%s?pulje=%s", nextId, nextPuljeID)
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
