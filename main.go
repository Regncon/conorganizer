package main

import (
	"fmt"
	"github.com/Regncon/conorganizer/pages/event"
	"github.com/Regncon/conorganizer/pages/root"
	"github.com/a-h/templ"
	"net/http"
)

func main() {

	http.Handle("/", templ.Handler(root.Page("Regncon 2025")))
	http.Handle("/event", templ.Handler(event.Page("Regncon 2025")))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
