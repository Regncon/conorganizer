package event

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/components/event_components"
	"github.com/Regncon/conorganizer/models"
	eventservice "github.com/Regncon/conorganizer/service/eventService"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar-go/datastar"
)

func selectedInterestHandler(db *sql.DB, logger *slog.Logger, eventImageDir *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := chi.URLParam(r, "idx")
		var signals struct {
			BillettHolderId int    `json:"billettHolderId"`
			PuljeId         string `json:"puljeId"`
		}
		if err := datastar.ReadSignals(r, &signals); err != nil {
			http.Error(w, "Ugyldig billettholdervalg.", http.StatusBadRequest)
			return
		}
		if signals.BillettHolderId <= 0 || signals.PuljeId == "" {
			http.Error(w, "Velg billettholder og pulje.", http.StatusBadRequest)
			return
		}

		userInfo := userctx.GetUserRequestInfo(r.Context())
		logger := logger.With("event_id", eventID, "pulje_id", signals.PuljeId, "billettholder_id", signals.BillettHolderId)
		var hasAccess bool
		err := db.QueryRowContext(r.Context(), `
			SELECT EXISTS (
				SELECT 1 FROM relation_billettholdere_users bu
				JOIN users u ON u.id = bu.user_id
				WHERE bu.billettholder_id = ?1 AND u.external_id = ?2
			)
		`, signals.BillettHolderId, userInfo.Id).Scan(&hasAccess)
		if err != nil {
			logger.Error(fmt.Errorf("failed to check billettholder access: %w", err).Error())
			http.Error(w, "Kunne ikke hente billettholder.", http.StatusInternalServerError)
			return
		}
		if !hasAccess {
			http.Error(w, "Du har ikke tilgang til denne billettholderen.", http.StatusForbidden)
			return
		}

		event, err := eventservice.GetEventById(eventID, db)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Kunne ikke hente arrangementet.", http.StatusInternalServerError)
			return
		}
		viewDecision, err := decideEventView(event, userInfo, db)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Kunne ikke sjekke tilgang til arrangementet.", http.StatusInternalServerError)
			return
		}
		if !viewDecision.CanView {
			http.Error(w, "Arrangementet er ikke tilgjengelig.", viewDecision.HiddenResponseStatusCode)
			return
		}

		notices, err := event_components.LoadInterestNoticeSignals(signals.BillettHolderId, signals.PuljeId, event.AgeGroup, eventImageDir, db)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Kunne ikke hente interessevarsler.", http.StatusInternalServerError)
			return
		}
		interest, err := getSelectedInterest(eventID, signals.BillettHolderId, signals.PuljeId, db)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Kunne ikke hente valgt interesse.", http.StatusInternalServerError)
			return
		}
		response := struct {
			event_components.InterestNoticeSignals
			SelectedInterestLevel      models.InterestLevel `json:"selectedInterestLevel"`
			CurrentInterestLevelChoice string               `json:"currentInterestLevelChoice"`
		}{notices, interest, "Pending choice"}
		sse := datastar.NewSSE(w, r)
		if err := sse.MarshalAndPatchSignals(response); err != nil {
			logger.Error(fmt.Errorf("failed to patch selected interest signals: %w", err).Error())
		}
	}
}
