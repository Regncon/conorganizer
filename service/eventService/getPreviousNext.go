package eventservice

import (
	"context"
	"database/sql"
	"github.com/Regncon/conorganizer/components"
	"log/slog"
)

func GetPreviousNext(ctx context.Context, db *sql.DB, logger *slog.Logger, currentID string) (components.PreviousNext, error) {
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
  WHERE status = ?
)
SELECT previous_id, previous_title, previous_image_url,
       next_id,     next_title,     next_image_url
FROM ordered
WHERE id = ?;`

	var (
		prevID, prevTitle, prevImg sql.NullString
		nextID, nextTitle, nextImg sql.NullString
	)

	err := db.QueryRowContext(ctx, q, "Innsendt", currentID).
		Scan(&prevID, &prevTitle, &prevImg, &nextID, &nextTitle, &nextImg)
	if err != nil {
		if err == sql.ErrNoRows {
			// currentID isn’t in the filtered set; return empty neighbors.
			return components.PreviousNext{}, nil
		}
		logger.Error("GetPreviousNext scan failed", "error", err)
		return components.PreviousNext{}, err
	}

	return components.PreviousNext{
		PreviousUrl:      nstr(prevID),
		PreviousTitle:    nstr(prevTitle),
		PreviousImageURL: nstr(prevImg),
		NextUrl:          nstr(nextID),
		NextTitle:        nstr(nextTitle),
		NextImageURL:     nstr(nextImg),
	}, nil
}

func nstr(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
