package root

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

const (
	rootPageLoadErrorMessage   = "Vi klarte ikke å laste forsiden akkurat nå. Prøv igjen om litt."
	rootEventsLoadErrorMessage = "Vi klarte ikke å laste arrangementene akkurat nå. Prøv igjen om litt."
)

func SetupRootRoute(router chi.Router, logger *slog.Logger, liveManager *live.Manager, db *sql.DB, eventImageDir *string) error {
	logger = logger.With("component", "root")
	rootLayoutRoute(router, db, logger, eventImageDir)

	router.Route("/root", func(rootRouter chi.Router) {
		rootRouter.Route("/api", func(rootApiRouter chi.Router) {
			rootApiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				requestedDate := r.URL.Query().Get(programDateQueryParam)
				liveManager.Stream(w, r, live.Page{
					Buckets: []live.Bucket{live.BucketEvents},
					Render: func(ctx context.Context, r *http.Request) templ.Component {
						interests := loadBillettholderInterestsForPage(ctx, userctx.GetUserRequestInfo(ctx), db, logger)
						return rootPage(db, eventImageDir, requestedDate, interests)
					},
				})
			})
		})
	})

	return nil
}

func loadBillettholderInterestsForPage(ctx context.Context, userInfo requestctx.UserRequestInfo, db *sql.DB, logger *slog.Logger) billettholderInterestsByEvent {
	interests, err := loadBillettholderInterests(userInfo, requestctx.SelectedBillettholderID(ctx), db)
	if err != nil {
		logger.Error(err.Error(), "user_id", userInfo.Id)
		return nil
	}
	return interests
}

func MustJSONMarshal(v any) string {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
