package main

import (
	"github.com/Regncon/conorganizer/pages/event/add"
	"github.com/Regncon/conorganizer/pages/root"
	"github.com/Regncon/conorganizer/service"
	"github.com/a-h/templ"
	"log"
	"net/http"
)

func main() {
	log.Println("Staring Regncon 2025")
	db, err := service.InitDB("events.db")
	if err != nil {
		log.Fatalf("Could not initialize DB: %v", err)
	}
	defer db.Close()

	http.Handle("/", templ.Handler(root.Page(db)))
	http.Handle("/event/add/", templ.Handler(add.Page(db)))

	http.HandleFunc("/event/add/new/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("EventNew handler")
		templ.Handler(add.EventNew(w, r, db)).Component.Render(r.Context(), w)
	})

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
