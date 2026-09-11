package varsler

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/Regncon/conorganizer/models"
)

var ErrPuljeNotFound = errors.New("pulje not found")

type notificationPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url"`
	Tag   string `json:"tag"`
}

type resultSnapshot struct {
	PuljeID       models.Pulje      `json:"puljeId"`
	PuljeName     string            `json:"puljeName"`
	StartAt       string            `json:"startAt"`
	EndAt         string            `json:"endAt"`
	Assignments   []assignmentState `json:"assignments"`
	Billettholder int               `json:"billettholderId"`
}

type assignmentState struct {
	EventID    string `json:"eventId"`
	EventTitle string `json:"eventTitle"`
	Role       string `json:"role"`
	RoomID     *int64 `json:"roomId"`
	RoomName   string `json:"roomName"`
}

func UpdatePuljeStatus(ctx context.Context, db *sql.DB, pulje models.Pulje, status models.PuljeStatus) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin pulje status transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var currentStatus models.PuljeStatus
	if err := tx.QueryRowContext(ctx, `SELECT status FROM puljer WHERE id = ?`, pulje).Scan(&currentStatus); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPuljeNotFound
		}
		return fmt.Errorf("get pulje %s status: %w", pulje, err)
	}
	if currentStatus == status {
		return tx.Commit()
	}
	if _, err := tx.ExecContext(ctx, `UPDATE puljer SET status = ? WHERE id = ?`, status, pulje); err != nil {
		return fmt.Errorf("update pulje %s status: %w", pulje, err)
	}

	if status != models.PuljeStatusCompleted {
		if _, err := tx.ExecContext(ctx, `
			UPDATE web_push_jobs
			SET status = 'canceled', lease_until = NULL, lease_token = NULL
			WHERE status IN ('pending','retrying','processing')
			  AND publication_id IN (SELECT id FROM pulje_varsel_publications WHERE pulje_id = ?)
		`, pulje); err != nil {
			return fmt.Errorf("cancel pending varsler for pulje %s: %w", pulje, err)
		}
		return tx.Commit()
	}

	var programPublished bool
	if err := tx.QueryRowContext(ctx, `SELECT is_published FROM program_publishing_state WHERE id = 1`).Scan(&programPublished); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tx.Commit()
		}
		return fmt.Errorf("get program publishing state: %w", err)
	}
	if !programPublished {
		return tx.Commit()
	}
	if err := publishPulje(ctx, tx, pulje); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit pulje %s publication: %w", pulje, err)
	}
	return nil
}

func publishPulje(ctx context.Context, tx *sql.Tx, pulje models.Pulje) error {
	snapshots, err := loadResultSnapshots(ctx, tx, pulje)
	if err != nil {
		return err
	}
	if len(snapshots) == 0 {
		return nil
	}
	previous, err := loadPreviousFingerprints(ctx, tx, pulje)
	if err != nil {
		return err
	}
	hasChanges := len(previous) != len(snapshots)
	if !hasChanges {
		for _, snapshot := range snapshots {
			canonical, err := json.Marshal(snapshot)
			if err != nil {
				return fmt.Errorf("encode result for billettholder %d: %w", snapshot.Billettholder, err)
			}
			hash := sha256.Sum256(canonical)
			if previous[snapshot.Billettholder] != hex.EncodeToString(hash[:]) {
				hasChanges = true
				break
			}
		}
	}
	if !hasChanges {
		return nil
	}

	var revision int
	if err := tx.QueryRowContext(ctx, `SELECT COALESCE(MAX(revision), 0) + 1 FROM pulje_varsel_publications WHERE pulje_id = ?`, pulje).Scan(&revision); err != nil {
		return fmt.Errorf("get next publication revision for pulje %s: %w", pulje, err)
	}
	result, err := tx.ExecContext(ctx, `INSERT INTO pulje_varsel_publications(pulje_id,revision) VALUES (?,?)`, pulje, revision)
	if err != nil {
		return fmt.Errorf("record pulje %s publication: %w", pulje, err)
	}
	publicationID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get pulje %s publication ID: %w", pulje, err)
	}

	for _, snapshot := range snapshots {
		canonical, err := json.Marshal(snapshot)
		if err != nil {
			return fmt.Errorf("encode result for billettholder %d: %w", snapshot.Billettholder, err)
		}
		hash := sha256.Sum256(canonical)
		fingerprint := hex.EncodeToString(hash[:])
		payload := notificationPayload{
			Title: "Programmet ditt er klart",
			Body:  fmt.Sprintf("Programmet ditt for %s er klart.", snapshot.PuljeName),
			URL:   "/profile?b_id=" + url.QueryEscape(fmt.Sprint(snapshot.Billettholder)) + "&pulje=" + url.QueryEscape(string(pulje)) + "#mitt-program",
			Tag:   fmt.Sprintf("pulje-%d-billettholder-%d", publicationID, snapshot.Billettholder),
		}
		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("encode notification for billettholder %d: %w", snapshot.Billettholder, err)
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO pulje_varsel_results(publication_id,billettholder_id,result_fingerprint,result_json,payload_json)
			VALUES (?,?,?,?,?)
		`, publicationID, snapshot.Billettholder, fingerprint, string(canonical), string(payloadJSON)); err != nil {
			return fmt.Errorf("record result for billettholder %d: %w", snapshot.Billettholder, err)
		}
		if previous[snapshot.Billettholder] == fingerprint {
			continue
		}
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO web_push_jobs(publication_id,billettholder_id,user_id,subscription_id)
			SELECT ?, ?, rbu.user_id, subscription.id
			FROM relation_billettholdere_users rbu
			JOIN web_push_subscriptions subscription ON subscription.user_id = rbu.user_id
			WHERE rbu.billettholder_id = ?
		`, publicationID, snapshot.Billettholder, snapshot.Billettholder); err != nil {
			return fmt.Errorf("queue varsler for billettholder %d: %w", snapshot.Billettholder, err)
		}
	}

	return nil
}

func UpdateProgramPublished(ctx context.Context, db *sql.DB, published bool) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin program publication transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO program_publishing_state(id,is_published)
		VALUES (1,?)
		ON CONFLICT(id) DO UPDATE SET is_published = excluded.is_published
	`, published); err != nil {
		return fmt.Errorf("update program publishing state: %w", err)
	}
	if !published {
		if _, err := tx.ExecContext(ctx, `
			UPDATE web_push_jobs
			SET status = 'canceled', lease_until = NULL, lease_token = NULL
			WHERE status IN ('pending','retrying','processing')
		`); err != nil {
			return fmt.Errorf("cancel pending varsler for unpublished program: %w", err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit program unpublication: %w", err)
		}
		return nil
	}

	rows, err := tx.QueryContext(ctx, `SELECT id FROM puljer WHERE status = ? ORDER BY start_at,id`, models.PuljeStatusCompleted)
	if err != nil {
		return fmt.Errorf("query completed puljer for program publication: %w", err)
	}
	var completedPuljer []models.Pulje
	for rows.Next() {
		var pulje models.Pulje
		if err := rows.Scan(&pulje); err != nil {
			_ = rows.Close()
			return fmt.Errorf("scan completed pulje for program publication: %w", err)
		}
		completedPuljer = append(completedPuljer, pulje)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		return fmt.Errorf("iterate completed puljer for program publication: %w", err)
	}
	if err := rows.Close(); err != nil {
		return fmt.Errorf("close completed puljer for program publication: %w", err)
	}
	for _, pulje := range completedPuljer {
		if err := publishPulje(ctx, tx, pulje); err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit program publication: %w", err)
	}
	return nil
}

func loadResultSnapshots(ctx context.Context, tx *sql.Tx, pulje models.Pulje) ([]resultSnapshot, error) {
	rows, err := tx.QueryContext(ctx, `
		WITH affected AS (
			SELECT player.billettholder_id
			FROM relation_events_players player
			JOIN relation_event_puljer ep ON ep.event_id = player.event_id AND ep.pulje_id = player.pulje_id
			JOIN events event ON event.id = player.event_id
			WHERE player.pulje_id = ? AND ep.is_in_pulje = 1 AND ep.is_published = 1 AND event.status = 'Annonsert'
			UNION
			SELECT interest.billettholder_id
			FROM interests interest
			JOIN relation_event_puljer ep ON ep.event_id = interest.event_id AND ep.pulje_id = interest.pulje_id
			JOIN events event ON event.id = interest.event_id
			WHERE interest.pulje_id = ? AND ep.is_in_pulje = 1 AND ep.is_published = 1 AND event.status = 'Annonsert'
		)
		SELECT affected.billettholder_id, pulje.name, pulje.start_at, pulje.end_at
		FROM affected
		JOIN puljer pulje ON pulje.id = ?
		ORDER BY affected.billettholder_id
	`, pulje, pulje, pulje)
	if err != nil {
		return nil, fmt.Errorf("query affected recipients for pulje %s: %w", pulje, err)
	}
	defer rows.Close()

	var snapshots []resultSnapshot
	for rows.Next() {
		var snapshot resultSnapshot
		snapshot.PuljeID = pulje
		if err := rows.Scan(&snapshot.Billettholder, &snapshot.PuljeName, &snapshot.StartAt, &snapshot.EndAt); err != nil {
			return nil, fmt.Errorf("scan affected recipient for pulje %s: %w", pulje, err)
		}
		snapshot.Assignments, err = loadAssignments(ctx, tx, pulje, snapshot.Billettholder)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, snapshot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate affected recipients for pulje %s: %w", pulje, err)
	}
	return snapshots, nil
}

func loadAssignments(ctx context.Context, tx *sql.Tx, pulje models.Pulje, billettholderID int) ([]assignmentState, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT event.id, event.title, player.role, ep.room_id, COALESCE(room.name, '')
		FROM relation_events_players player
		JOIN relation_event_puljer ep ON ep.event_id = player.event_id AND ep.pulje_id = player.pulje_id
		JOIN events event ON event.id = player.event_id
		LEFT JOIN rooms room ON room.id = ep.room_id
		WHERE player.pulje_id = ? AND player.billettholder_id = ?
		  AND ep.is_in_pulje = 1 AND ep.is_published = 1 AND event.status = 'Annonsert'
		ORDER BY event.id, player.role
	`, pulje, billettholderID)
	if err != nil {
		return nil, fmt.Errorf("query assignments for billettholder %d: %w", billettholderID, err)
	}
	defer rows.Close()

	assignments := make([]assignmentState, 0)
	for rows.Next() {
		var assignment assignmentState
		var roomID sql.NullInt64
		if err := rows.Scan(&assignment.EventID, &assignment.EventTitle, &assignment.Role, &roomID, &assignment.RoomName); err != nil {
			return nil, fmt.Errorf("scan assignment for billettholder %d: %w", billettholderID, err)
		}
		if roomID.Valid {
			assignment.RoomID = &roomID.Int64
		}
		assignments = append(assignments, assignment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate assignments for billettholder %d: %w", billettholderID, err)
	}
	return assignments, nil
}

func loadPreviousFingerprints(ctx context.Context, tx *sql.Tx, pulje models.Pulje) (map[int]string, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT result.billettholder_id, result.result_fingerprint
		FROM pulje_varsel_results result
		JOIN pulje_varsel_publications publication ON publication.id = result.publication_id
		WHERE publication.pulje_id = ?
		  AND publication.revision = (SELECT MAX(revision) FROM pulje_varsel_publications WHERE pulje_id = ?)
	`, pulje, pulje)
	if err != nil {
		return nil, fmt.Errorf("query previous results for pulje %s: %w", pulje, err)
	}
	defer rows.Close()

	previous := make(map[int]string)
	for rows.Next() {
		var billettholderID int
		var fingerprint string
		if err := rows.Scan(&billettholderID, &fingerprint); err != nil {
			return nil, fmt.Errorf("scan previous result for pulje %s: %w", pulje, err)
		}
		previous[billettholderID] = fingerprint
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate previous results for pulje %s: %w", pulje, err)
	}
	return previous, nil
}
