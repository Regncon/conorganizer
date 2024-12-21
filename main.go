package main

import (
	"database/sql"
	"fmt"
	"github.com/Regncon/conorganizer/pages/root"
	"github.com/Regncon/conorganizer/service"
	"github.com/a-h/templ"
	"log"
	"net/http"
)

func createEventsTable(db *sql.DB) error {
	tableCreationQuery := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL
	)`

	_, err := db.Exec(tableCreationQuery)
	if err != nil {
		return fmt.Errorf("failed to create events table: %w", err)
	}

	return nil
}

func main() {
	fmt.Println("Regncon 2025")
	// Initialize the service layer
	// eventService, err := service.Initialize("regncon.db")
	// if err != nil {
	// 	fmt.Printf("Failed to initialize service: %v\n", err)
	// 	os.Exit(1)
	// }
	//
	// Pass the service to the handler
	db, err := service.InitDB("events.db")
	if err != nil {
		log.Fatalf("Could not initialize DB: %v", err)
	}
	defer db.Close()

	// Ensure the `events` table exists
	if err := createEventsTable(db); err != nil {
		log.Fatalf("Failed to create events table: %v", err)
	}
	query := "INSERT INTO events (name, description) VALUES ('Evnt 1', 'This is the first event')"
	result, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result.LastInsertId())

	eventList, err := root.GetEvents(db)
	fmt.Println(err)
	fmt.Println(eventList)

	http.Handle("/", templ.Handler(root.Page("Regncon 2025", db)))

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
