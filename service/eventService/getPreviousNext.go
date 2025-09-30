package eventservice

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"

	"github.com/Regncon/conorganizer/components"
	"github.com/Regncon/conorganizer/service/eventimage"
)

func GetPreviousNext(ctx context.Context, db *sql.DB, logger *slog.Logger, currentID string, eventImageDir *string) (components.PreviousNext, error) {
	const q = `
WITH ordered AS (
  SELECT
    id,
    title,
    image_url,
    LAG(id)        OVER (ORDER BY inserted_time DESC, id DESC) AS previous_id,
    LAG(title)     OVER (ORDER BY inserted_time DESC, id DESC) AS previous_title,
    LAG(image_url) OVER (ORDER BY inserted_time DESC, id DESC) AS previous_image_url,
    LEAD(id)       OVER (ORDER BY inserted_time DESC, id DESC) AS next_id,
    LEAD(title)    OVER (ORDER BY inserted_time DESC, id DESC) AS next_title,
    LEAD(image_url)OVER (ORDER BY inserted_time DESC, id DESC) AS next_image_url
  FROM events
  WHERE status IN ('Innsendt', 'Godkjent')
)
SELECT previous_id, previous_title, previous_image_url,
       next_id,     next_title,     next_image_url
FROM ordered
WHERE id = ?;`

	var (
		prevID, prevTitle, prevImg sql.NullString
		nextID, nextTitle, nextImg sql.NullString
	)

	//eventsQuery := "SELECT id, title, intro, status, system, host_name,beginner_friendly, event_type, age_group, event_runtime, can_be_run_in_english FROM events WHERE status IN ('Innsendt', 'Godkjent') ORDER BY inserted_time DESC"
	err := db.QueryRowContext(ctx, q, currentID).
		Scan(&prevID, &prevTitle, &prevImg, &nextID, &nextTitle, &nextImg)
	if err != nil {
		if err == sql.ErrNoRows {
			// currentID isn’t in the filtered set; return empty neighbors.
			return components.PreviousNext{}, nil
		}
		logger.Error("GetPreviousNext scan failed", "error", err)
		return components.PreviousNext{}, err
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
