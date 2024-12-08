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

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/event", templ.Handler(event.Page("Regncon 2025")))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
