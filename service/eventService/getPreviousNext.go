package eventservice

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Regncon/conorganizer/components"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/eventimage"
)

func GetPreviousNextInnsendtGodkjent(ctx context.Context, db *sql.DB, currentID string, eventImageDir *string) (components.PreviousNext, error) {
	const q = `
        WITH ordered AS (
            SELECT
                id,
                title,
                LAG(id)     OVER (ORDER BY created_at DESC, id DESC) AS previous_id,
                LAG(title)  OVER (ORDER BY created_at DESC, id DESC) AS previous_title,
                LEAD(id)    OVER (ORDER BY created_at DESC, id DESC) AS next_id,
                LEAD(title) OVER (ORDER BY created_at DESC, id DESC) AS next_title
            FROM events
            WHERE status IN (?, ?)
        )
        SELECT previous_id, previous_title,
               next_id,     next_title
        FROM ordered
        WHERE id = ?;`

	var (
		prevID, prevTitle sql.NullString
		nextID, nextTitle sql.NullString
	)

	err := db.QueryRowContext(ctx, q, models.EventStatusSubmitted, models.EventStatusApproved, currentID).
		Scan(&prevID, &prevTitle, &nextID, &nextTitle)
	if err != nil {
		if err == sql.ErrNoRows {
			// currentID isn't in the filtered set; return empty neighbors.
			return components.PreviousNext{}, nil
		}
		return components.PreviousNext{}, fmt.Errorf("get previous/next scan failed for event %q: %w", currentID, err)
	}

	PrevImageBannerUrl := eventimage.GetEventImageUrl(nstr(prevID), "banner", eventImageDir)
	NextImageBannerUrl := eventimage.GetEventImageUrl(nstr(nextID), "banner", eventImageDir)

	if strings.Contains(PrevImageBannerUrl, "placeholder") {
		PrevImageBannerUrl = ""
	}
	if strings.Contains(NextImageBannerUrl, "placeholder") {
		NextImageBannerUrl = ""
	}

	return components.PreviousNext{
		PreviousUrl:      nstr(prevID),
		PreviousTitle:    nstr(prevTitle),
		PreviousImageURL: PrevImageBannerUrl,
		NextUrl:          nstr(nextID),
		NextTitle:        nstr(nextTitle),
		NextImageURL:     NextImageBannerUrl,
	}, nil
}

func nstr(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
