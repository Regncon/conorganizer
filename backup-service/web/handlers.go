package web

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/backup-service/services"
	"github.com/Regncon/conorganizer/backup-service/web/pages"
	"github.com/a-h/templ"

	components "github.com/Regncon/conorganizer/backup-service/web/components/sections/header"
	layouts "github.com/Regncon/conorganizer/backup-service/web/layout"
)

type Handlers struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func (h *Handlers) IndexHandler(res http.ResponseWriter, req *http.Request) {
	stats, err := services.FetchLog(h.DB).Stats()
	if err != nil {
		h.Logger.Error("failed to load backup stats", "err", err)
		stats = services.BackupStats{}
	}

	logs, err := services.FetchLog(h.DB).Logs("", "", 99)
	if err != nil {
		h.Logger.Error("failed to fetch hourly logs", "err", err)
		// optionally return an error page
		return
	}

	templ.Handler(layouts.Base(components.Header(stats), pages.Index(logs))).ServeHTTP(res, req)
}

/* func (h *Handlers) IntervalHandler(res http.ResponseWriter, req *http.Request) {
	templ.Handler(layouts.Base(components.Header(1, 2, 3, 4), pages.Index())).ServeHTTP(res, req)
} */

/* func (h *Handlers) LogSearch(res http.ResponseWriter, req *http.Request) {
	templ.Handler(layouts.Base(components.Header(1, 2, 3, 4), pages.Index())).ServeHTTP(res, req)
} */
