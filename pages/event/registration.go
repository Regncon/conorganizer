package event

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	datastar "github.com/starfederation/datastar-go/datastar"
)

var (
	errRegistrationProgramNotPublished = errors.New("program is not published")
	errRegistrationEventUnavailable    = errors.New("event is not available in pulje")
	errRegistrationPuljeNotOpen        = errors.New("pulje is not open")
	errRegistrationAccessDenied        = errors.New("user does not have access to billettholder")
	errRegistrationEventNotOpen        = errors.New("event is not open for registration")
	errRegistrationAdultsOnly          = errors.New("event is adults only")
	errRegistrationGamemaster          = errors.New("billettholder is a gamemaster in pulje")
)

type registrationChange struct {
	UserExternalID  string
	BillettholderID int
	EventID         string
	PuljeID         models.Pulje
	IsRegistered    bool
}

func setRegistration(ctx context.Context, db *sql.DB, change registrationChange) error {
	if change.UserExternalID == "" {
		return errRegistrationAccessDenied
	}
	if change.BillettholderID <= 0 {
		return fmt.Errorf("billettholder id is required")
	}
	if change.EventID == "" {
		return fmt.Errorf("event id is required")
	}
	if _, ok := models.ParsePulje(string(change.PuljeID)); !ok {
		return fmt.Errorf("valid pulje id is required")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin registration change: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var programPublished int
	if err := tx.QueryRowContext(ctx, `
		SELECT is_published
		FROM program_publishing_state
		WHERE id = 1
	`).Scan(&programPublished); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errRegistrationProgramNotPublished
		}
		return fmt.Errorf("read program publishing state: %w", err)
	}
	if programPublished == 0 {
		return errRegistrationProgramNotPublished
	}

	var (
		isOpenRegistration int
		ageGroup           models.AgeGroup
		puljeStatus        models.PuljeStatus
		isOver18           int
		userHasAccess      int
		isGamemaster       int
	)
	if err := tx.QueryRowContext(ctx, `
		SELECT
			e.is_open_registration,
			e.age_group,
			p.status,
			b.is_over_18,
			EXISTS (
				SELECT 1
				FROM relation_billettholdere_users bu
				JOIN users u ON u.id = bu.user_id
				WHERE bu.billettholder_id = b.id
					AND u.external_id = ?
			),
			EXISTS (
				SELECT 1
				FROM relation_events_players assigned
				WHERE assigned.billettholder_id = b.id
					AND assigned.pulje_id = p.id
					AND assigned.role = ?
			)
		FROM events e
		JOIN relation_event_puljer ep ON ep.event_id = e.id
		JOIN puljer p ON p.id = ep.pulje_id
		JOIN billettholdere b ON b.id = ?
		WHERE e.id = ?
			AND p.id = ?
			AND e.status = ?
			AND ep.is_in_pulje = 1
			AND ep.is_published = 1
	`,
		change.UserExternalID,
		models.EventPlayerRoleGM,
		change.BillettholderID,
		change.EventID,
		change.PuljeID,
		models.EventStatusAnnounced,
	).Scan(
		&isOpenRegistration,
		&ageGroup,
		&puljeStatus,
		&isOver18,
		&userHasAccess,
		&isGamemaster,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errRegistrationEventUnavailable
		}
		return fmt.Errorf("read registration eligibility: %w", err)
	}

	if userHasAccess == 0 {
		return errRegistrationAccessDenied
	}
	if puljeStatus != models.PuljeStatusOpen {
		return fmt.Errorf("%w: %s", errRegistrationPuljeNotOpen, puljeStatus)
	}
	if isGamemaster != 0 {
		return errRegistrationGamemaster
	}
	if change.IsRegistered && isOpenRegistration == 0 {
		return errRegistrationEventNotOpen
	}
	if change.IsRegistered && ageGroup == models.AgeGroupAdultsOnly && isOver18 == 0 {
		return errRegistrationAdultsOnly
	}

	if change.IsRegistered {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO relation_events_players (
				event_id,
				pulje_id,
				billettholder_id,
				role,
				source
			)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(billettholder_id, event_id, pulje_id) DO NOTHING
		`,
			change.EventID,
			change.PuljeID,
			change.BillettholderID,
			models.EventPlayerRolePlayer,
			models.EventPlayerSourceRegistration,
		); err != nil {
			return fmt.Errorf("add registration: %w", err)
		}
	} else {
		if _, err := tx.ExecContext(ctx, `
			DELETE FROM relation_events_players
			WHERE event_id = ?
				AND pulje_id = ?
				AND billettholder_id = ?
				AND role = ?
				AND source IN (?, ?)
		`,
			change.EventID,
			change.PuljeID,
			change.BillettholderID,
			models.EventPlayerRolePlayer,
			models.EventPlayerSourceManual,
			models.EventPlayerSourceRegistration,
		); err != nil {
			return fmt.Errorf("remove registration: %w", err)
		}
	}

	if _, err := tx.ExecContext(ctx, `
		DELETE FROM interests
		WHERE event_id = ?
			AND pulje_id = ?
			AND billettholder_id = ?
	`, change.EventID, change.PuljeID, change.BillettholderID); err != nil {
		return fmt.Errorf("remove interest for registration change: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit registration change: %w", err)
	}
	return nil
}

func setupRegistrationRoute(router chi.Router, db *sql.DB, liveManager *live.Manager, logger *slog.Logger) {
	logger = logger.With("component", "registration")

	router.Put("/", func(w http.ResponseWriter, r *http.Request) {
		type signals struct {
			BillettholderID int    `json:"billettHolderId"`
			PuljeID         string `json:"puljeId"`
			IsRegistered    bool   `json:"isRegistered"`
		}
		store := &signals{}
		if err := datastar.ReadSignals(r, store); err != nil {
			http.Error(w, "Klarte ikkje å lese skjemadata", http.StatusBadRequest)
			return
		}

		eventID := chi.URLParam(r, "idx")
		puljeID, validPulje := models.ParsePulje(store.PuljeID)
		if eventID == "" || store.BillettholderID <= 0 || !validPulje {
			http.Error(w, "Arrangement, pulje eller billettheldar manglar", http.StatusBadRequest)
			return
		}

		userInfo := userctx.GetUserRequestInfo(r.Context())
		err := setRegistration(r.Context(), db, registrationChange{
			UserExternalID:  userInfo.Id,
			BillettholderID: store.BillettholderID,
			EventID:         eventID,
			PuljeID:         puljeID,
			IsRegistered:    store.IsRegistered,
		})
		if err != nil {
			statusCode, message, expected := registrationErrorResponse(err)
			if !expected {
				logger.Error(err.Error(),
					"event_id", eventID,
					"pulje_id", puljeID,
					"billettholder_id", store.BillettholderID,
					"user_id", userInfo.Id,
					"request_id", middleware.GetReqID(r.Context()),
				)
			}
			http.Error(w, message, statusCode)
			return
		}

		if err := liveManager.Broadcast(r.Context(), live.BucketEvents, live.BucketInterests); err != nil {
			logger.Error(fmt.Errorf("broadcast registration change: %w", err).Error(),
				"event_id", eventID,
				"pulje_id", puljeID,
				"billettholder_id", store.BillettholderID,
				"request_id", middleware.GetReqID(r.Context()),
			)
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

func registrationErrorResponse(err error) (statusCode int, message string, expected bool) {
	switch {
	case errors.Is(err, errRegistrationAccessDenied):
		return http.StatusForbidden, "Du har ikkje tilgang til denne billettheldaren.", true
	case errors.Is(err, errRegistrationAdultsOnly):
		return http.StatusForbidden, "Arrangementet har 18-årsgrense.", true
	case errors.Is(err, errRegistrationGamemaster):
		return http.StatusConflict, "Du er spilleder i denne pulja og kan ikkje endre påmelding.", true
	case errors.Is(err, errRegistrationProgramNotPublished):
		return http.StatusConflict, "Programmet er ikkje publisert enno.", true
	case errors.Is(err, errRegistrationEventUnavailable):
		return http.StatusConflict, "Arrangementet er ikkje tilgjengeleg i denne pulja.", true
	case errors.Is(err, errRegistrationPuljeNotOpen):
		return http.StatusConflict, "Pulja er ikkje open for endringar.", true
	case errors.Is(err, errRegistrationEventNotOpen):
		return http.StatusConflict, "Arrangementet er ikkje ope for påmelding.", true
	default:
		return http.StatusInternalServerError, "Klarte ikkje å endre påmeldinga.", false
	}
}
