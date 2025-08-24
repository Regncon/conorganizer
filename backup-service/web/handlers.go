package web

import (
	"net/http"

	layouts "github.com/Regncon/conorganizer/backup-service/web/layout"
	"github.com/Regncon/conorganizer/backup-service/web/pages"
	"github.com/a-h/templ"
)

func IndexHandler(res http.ResponseWriter, req *http.Request) {
	templ.Handler(layouts.Base(pages.Index())).ServeHTTP(res, req)
}

func IntervalHandler(res http.ResponseWriter, req *http.Request) {
	templ.Handler(layouts.Base(pages.Index())).ServeHTTP(res, req)
}

func LogSearch(res http.ResponseWriter, req *http.Request) {
	templ.Handler(layouts.Base(pages.Index())).ServeHTTP(res, req)
}
