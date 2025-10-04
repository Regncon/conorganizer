package checkIn

import (
	"errors"
	"log/slog"
)

// AssociateUserWithBillettholder uses userID string from
func AssociateUserWithBillettholder(userID string, logger *slog.Logger) error {
	logger.Info("Associating userID with billettholder", "userID", userID)

	return errors.New("user is not associated")
}
