package rooms

import "database/sql"

// ValidateRooms validates that all entries in `room` is valid
func ValidateRooms(db *sql.DB)

// ValidateRoomsByPulje is used for validating that a snapshot of rooms, based on a pulje, is valid (eg. assigned vs max concurrent events per pulje)
func ValidateRoomsByPulje(db *sql.DB, puljeID string)

// ValidateDisabledRoomsCascade Validates that room disabled status has cascaded and no orphans exist in `relation_event_puljer`
func ValidateDisabledRoomsCascade(db *sql.DB)
