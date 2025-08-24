package web

import (
	"net/http"

	"github.com/Regncon/conorganizer/backup-service/web/pages"
	"github.com/a-h/templ"

	components "github.com/Regncon/conorganizer/backup-service/web/components/sections/header"
	layouts "github.com/Regncon/conorganizer/backup-service/web/layout"
)

func IndexHandler(res http.ResponseWriter, req *http.Request) {
	templ.Handler(layouts.Base(components.Header(), pages.Index())).ServeHTTP(res, req)
}

func IntervalHandler(res http.ResponseWriter, req *http.Request) {
	templ.Handler(layouts.Base(components.Header(), pages.Index())).ServeHTTP(res, req)
}

func LogSearch(res http.ResponseWriter, req *http.Request) {
	templ.Handler(layouts.Base(components.Header(), pages.Index())).ServeHTTP(res, req)
}
