package root

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func SetupRootRoute(router chi.Router, logger *slog.Logger, liveManager *live.Manager, db *sql.DB, eventImageDir *string) error {
	logger = logger.With("component", "root")
	rootLayoutRoute(router, db, logger, eventImageDir)

	router.Route("/root", func(rootRouter chi.Router) {
		rootRouter.Route("/api", func(rootApiRouter chi.Router) {
			rootApiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				liveManager.Stream(w, r, live.Page{
					Buckets: []live.Bucket{live.BucketEvents},
					Render: func(ctx context.Context, r *http.Request) templ.Component {
						isAdmin := authctx.GetAdminFromUserToken(ctx)
                        userInfo := userctx.GetUserRequestInfo(ctx)
						return rootPage(userInfo, db, isAdmin, eventImageDir, logger)
					},
				})
			})
		})
	})

	return nil
}

func MustJSONMarshal(v any) string {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
