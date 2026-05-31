package root

import (
	"database/sql"
	"fmt"

	"github.com/Regncon/conorganizer/models"
	puljerService "github.com/Regncon/conorganizer/service/puljer"
)

type RootEventOccurrence struct {
	Event   models.EventCardModel
	PuljeID models.Pulje
}

func GetPublishedEventPuljeBlocks(db *sql.DB) ([]PuljeBlock, error) {
	eventsByPulje, err := GetEventsByPulje(db)
	if err != nil {
		return nil, err
	}

	puljer, err := puljerService.GetAllPuljer(db)
	if err != nil {
		return nil, fmt.Errorf("query puljer for published root events: %w", err)
	}

	blocks := make([]PuljeBlock, 0, len(puljer))
	for _, pulje := range puljer {
		block := eventsByPulje[pulje.ID]
		if block == nil {
			continue
		}
		blocks = append(blocks, *block)
	}

	return blocks, nil
}

func GetPublishedEventOccurrences(db *sql.DB) ([]RootEventOccurrence, error) {
	blocks, err := GetPublishedEventPuljeBlocks(db)
	if err != nil {
		return nil, err
	}

	occurrences := make([]RootEventOccurrence, 0)
	for _, block := range blocks {
		for _, event := range block.Events {
			occurrences = append(occurrences, RootEventOccurrence{
				Event:   event,
				PuljeID: block.Pulje.ID,
			})
		}
	}

	return occurrences, nil
}
