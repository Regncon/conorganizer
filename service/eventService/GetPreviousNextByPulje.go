package eventservice

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Regncon/conorganizer/components"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/eventimage"
)

func GetPreviousNextByPulje(
	ctx context.Context,
	db *sql.DB,
	logger *slog.Logger,
	currentID string,
	isAdmin bool,
	eventImageDir *string,
) (components.PreviousNext, error) {
	// Filter toggles: admins can see unpublished; users cannot.
	pubFilter := "AND ep.is_published = 1"
	if isAdmin {
		pubFilter = ""
	}

	// 1) Get neighbors within the same pulje.
	qWithin := fmt.Sprintf(`
WITH filtered AS (
	SELECT
		e.id,
		e.title,
		e.image_url,
		e.inserted_time,
		p.id         AS pulje_id,
		p.start_time AS pulje_start_time
	FROM events e
	JOIN event_puljer ep ON ep.event_id = e.id
	JOIN puljer p       ON p.id = ep.pulje_id
	WHERE e.status = 'Godkjent'
	  AND ep.is_active = 1
	  %s
),
ranked AS (
	SELECT
		id,
		title,
		image_url,
		pulje_id,
		pulje_start_time,
		LAG(id)        OVER (PARTITION BY pulje_id ORDER BY inserted_time DESC, id DESC) AS prev_id,
		LAG(title)     OVER (PARTITION BY pulje_id ORDER BY inserted_time DESC, id DESC) AS prev_title,
		LAG(image_url) OVER (PARTITION BY pulje_id ORDER BY inserted_time DESC, id DESC) AS prev_img,
		LEAD(id)       OVER (PARTITION BY pulje_id ORDER BY inserted_time DESC, id DESC) AS next_id,
		LEAD(title)    OVER (PARTITION BY pulje_id ORDER BY inserted_time DESC, id DESC) AS next_title,
		LEAD(image_url)OVER (PARTITION BY pulje_id ORDER BY inserted_time DESC, id DESC) AS next_img
	FROM filtered
)
SELECT
	pulje_id,
	pulje_start_time,
	prev_id,  prev_title,  prev_img,
	next_id,  next_title,  next_img
FROM ranked
WHERE id = ?
LIMIT 1;
`, pubFilter)

	var (
		curPuljeID, curPuljeStart sql.NullString
		prevID, prevTit, prevImg  sql.NullString
		nextID, nextTit, nextImg  sql.NullString
	)

	err := db.QueryRowContext(ctx, qWithin, currentID).Scan(
		&curPuljeID, &curPuljeStart,
		&prevID, &prevTit, &prevImg,
		&nextID, &nextTit, &nextImg,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return components.PreviousNext{}, nil
		}
		logger.Error("GetPreviousNextByPulje: scan within failed", "error", err)
		return components.PreviousNext{}, err
	}

	// 2) If there's no "next" in the same pulje, jump to the first event of the next pulje.
	var nextPuljeID sql.NullString
	if !nextID.Valid {
		qNextPulje := fmt.Sprintf(`
WITH filtered AS (
	SELECT
		e.id,
		e.title,
		e.image_url,
		p.id         AS pulje_id,
		p.start_time AS pulje_start_time,
		e.inserted_time
	FROM events e
	JOIN event_puljer ep ON ep.event_id = e.id
	JOIN puljer p       ON p.id = ep.pulje_id
	WHERE e.status = 'Godkjent'
	  AND ep.is_active = 1
	  %s
)
SELECT
	id, title, image_url, pulje_id
FROM filtered
WHERE pulje_start_time > ?
ORDER BY pulje_start_time ASC, inserted_time DESC, id DESC
LIMIT 1;
`, pubFilter)

		var nid, ntit, nimg, npul sql.NullString
		if err := db.QueryRowContext(ctx, qNextPulje, curPuljeStart.String).
			Scan(&nid, &ntit, &nimg, &npul); err != nil && err != sql.ErrNoRows {
			logger.Error("GetPreviousNextByPulje: next pulje lookup failed", "error", err)
			return components.PreviousNext{}, err
		} else if err == nil {
			nextID, nextTit, nextImg, nextPuljeID = nid, ntit, nimg, npul
		}
	} else {
		// Next is within the same pulje.
		nextPuljeID = curPuljeID
	}

	// 3) Previous pulje should be set ONLY if a previous event exists.
	var prevPuljeID sql.NullString
	if prevID.Valid {
		prevPuljeID = curPuljeID
	} // else leave as NULL (empty)

	// 4) Resolve images and blank placeholders.
	prevBanner := ""
	if prevID.Valid {
		prevBanner = eventimage.GetEventImageUrl(prevID.String, "banner", eventImageDir)
	}
	nextBanner := ""
	if nextID.Valid {
		nextBanner = eventimage.GetEventImageUrl(nextID.String, "banner", eventImageDir)
	}
	if strings.Contains(prevBanner, "placeholder") {
		prevBanner = ""
	}
	if strings.Contains(nextBanner, "placeholder") {
		nextBanner = ""
	}

	// 5) Build response.
	return components.PreviousNext{
		PreviousUrl:      nstr(prevID),
		PreviousTitle:    nstr(prevTit),
		PreviousImageURL: prevBanner,
		PreviousPulje:    models.Pulje(nstr(prevPuljeID)),
		NextUrl:          nstr(nextID),
		NextTitle:        nstr(nextTit),
		NextImageURL:     nextBanner,
		NextPulje:        models.Pulje(nstr(nextPuljeID)),
		// IsRemoved left default false
	}, nil
}
