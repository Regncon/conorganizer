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
	// ---------- 1) Find current row in the filtered, pulje-partitioned ordering ----------
	// We compute prev/next *within the same pulje*. If `next` is NULL, we will later
	// jump to the first event of the next pulje (no wrap to the beginning).
	pubFilter := "AND ep.is_published = 1"
	if isAdmin {
		pubFilter = "" // admins can see unpublished
	}

	qWithin := fmt.Sprintf(`
WITH filtered AS (
	SELECT
		e.id,
		e.title,
		e.image_url,
		e.inserted_time,
		p.id          AS pulje_id,
		p.start_time  AS pulje_start_time
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
		curPuljeID      sql.NullString
		curPuljeStart   sql.NullString
		prevID, prevTit sql.NullString
		prevImg         sql.NullString
		nextID, nextTit sql.NullString
		nextImg         sql.NullString
	)

	err := db.QueryRowContext(ctx, qWithin, currentID).Scan(
		&curPuljeID, &curPuljeStart,
		&prevID, &prevTit, &prevImg,
		&nextID, &nextTit, &nextImg,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// currentID not visible under this filter -> empty neighbors
			return components.PreviousNext{}, nil
		}
		logger.Error("GetPreviousNextByPulje: scan within failed", "error", err)
		return components.PreviousNext{}, err
	}

	// ---------- 2) If `next` within pulje is NULL, hop to first event of next pulje ----------
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
		// next is within the same pulje as current
		nextPuljeID = curPuljeID
	}

	// ---------- 3) Previous from prior pulje is NOT allowed (no wrap back) ----------
	prevPuljeID := curPuljeID // previous stays within current pulje (or empty)

	// ---------- 4) Resolve image banners and blank placeholders ----------
	prevBanner := ""
	nextBanner := ""
	if prevID.Valid {
		prevBanner = eventimage.GetEventImageUrl(prevID.String, "banner", eventImageDir)
	}
	if nextID.Valid {
		nextBanner = eventimage.GetEventImageUrl(nextID.String, "banner", eventImageDir)
	}
	if strings.Contains(prevBanner, "placeholder") {
		prevBanner = ""
	}
	if strings.Contains(nextBanner, "placeholder") {
		nextBanner = ""
	}

	// ---------- 5) Build response ----------
	out := components.PreviousNext{
		PreviousUrl:      nstr(prevID),
		PreviousTitle:    nstr(prevTit),
		PreviousImageURL: prevBanner,
		NextUrl:          nstr(nextID),
		NextTitle:        nstr(nextTit),
		NextImageURL:     nextBanner,
		// These fields exist on your expanded struct
		PreviousPulje: models.Pulje(nstr(prevPuljeID)),
		NextPulje:     models.Pulje(nstr(nextPuljeID)),
		// IsRemoved: left as default false — set by caller if needed.
	}

	return out, nil
}
