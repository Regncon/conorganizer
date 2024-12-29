package edit

import (
	"database/sql"
	"net/http"
	"strconv"
)

func Save(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Parse form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		// Extract form values
		id := r.FormValue("id")
		title := r.FormValue("title")
		description := r.FormValue("description")

		// Convert ID to int64
		eventID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			http.Error(w, "Invalid event ID", http.StatusBadRequest)
			return
		}

		// Update event in the database
		query := `UPDATE events SET name = ?, description = ? WHERE id = ?`
		_, err = db.Exec(query, title, description, eventID)
		if err != nil {
			http.Error(w, "Failed to update event in the database", http.StatusInternalServerError)
			return
		}

		// Redirect or render success message
		http.Redirect(w, r, "/event/"+id, http.StatusSeeOther)
	}
}
