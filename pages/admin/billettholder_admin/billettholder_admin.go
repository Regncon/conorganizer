package billettholderadmin

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	addbillettholder "github.com/Regncon/conorganizer/pages/admin/billettholder_admin/add"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func SetupBillettholderAdminRoute(router chi.Router, liveManager *live.Manager, logger *slog.Logger, db *sql.DB) error {
	baseLogger := logger
	logger = logger.With("component", "billettholder_admin")

	indexRoute(router, db, logger)

	router.Route("/admin/billettholder/api/", func(billettholderAdminRouter chi.Router) {
		billettholderAdminRouter.With(authctx.RequireAdmin(baseLogger)).Get("/", func(w http.ResponseWriter, r *http.Request) {
			liveManager.Stream(w, r, live.Page{
				Buckets: []live.Bucket{live.BucketBillettholders, live.BucketInterests},
				Render: func(ctx context.Context, r *http.Request) templ.Component {
					return BillettholderAdminPage(db, logger)
				},
			})
		})
		billettholdereSearchRoute(billettholderAdminRouter, liveManager)
		addEmailToBilettholderRoute(billettholderAdminRouter, db, logger, liveManager)
		deleteEmailFromBillettholderRoute(billettholderAdminRouter, db, logger, liveManager)
	})

	addbillettholder.AddBillettholderRoute(router, db, logger)

	router.Route("/admin/billettholder/add/api/", func(addBillettholderRouter chi.Router) {
		addBillettholderRouter.With(authctx.RequireAdmin(baseLogger)).Get("/", func(w http.ResponseWriter, r *http.Request) {
			liveManager.Stream(w, r, live.Page{
				Buckets: []live.Bucket{live.BucketBillettholders},
				Render: func(ctx context.Context, r *http.Request) templ.Component {
					return addbillettholder.AddBillettholderAdminPage(db, logger)
				},
			})
		})

		addbillettholder.CheckInTicketsSearchRoute(addBillettholderRouter, db, logger, liveManager)
		addbillettholder.ConvertTicketToBillettholderRoute(addBillettholderRouter, db, liveManager, logger)
	})

	return nil
}
