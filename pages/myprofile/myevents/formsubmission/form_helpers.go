package formsubmission

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
)

var validColumns = map[string]struct{}{
	"host_name":   {},
	"description": {},
	"start_time":  {},
	"end_time":    {},
}

func isValidColumn(col string) bool {
	_, ok := validColumns[col]
	return ok
}

func saveEventField(
	db *sql.DB,
	logger *slog.Logger,
	column string,
	value interface{},
	eventID int64,
	w http.ResponseWriter,
	r *http.Request,
) (int64, error) {
	if !isValidColumn(column) {
		logger.Error("Invalid column name", "column", column)
		http.Error(w, "Invalid field", http.StatusBadRequest)
		return 0, fmt.Errorf("invalid column name: %s", column)
	}
	query := fmt.Sprintf(
		"UPDATE events SET %s = ? WHERE id = ?",
		column,
	)

	result, err := db.Exec(query, value)
	if err != nil {
		logger.Error("Error inserting "+column, "err", err)
		http.Error(w, fmt.Sprintf("Error inserting %s: %v", column, err), http.StatusBadRequest)
		return 0, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Error("Error getting rows affected for "+column, "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return 0, err
	}
	if rows == 0 {
		logger.Error("No rows inserted for " + column)
		http.Error(w, "Nothing inserted", http.StatusNotFound)
		return 0, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		// Some drivers don’t support LastInsertId; log a warning but still succeed.
		logger.Warn("Could not get LastInsertId for "+column, "warn", err)
		return 0, nil
	}

	return id, nil
}
