package rooms

var roomMaps = map[string]string{
	"705": "/static/rooms/terminus-7-etasje-705.svg",
	"706": "/static/rooms/terminus-7-etasje-706.svg",
	"707": "/static/rooms/terminus-7-etasje-707.svg",
	"709": "/static/rooms/terminus-7-etasje-709.svg",
	"710": "/static/rooms/terminus-7-etasje-710.svg",
	"711": "/static/rooms/terminus-7-etasje-711.svg",
	"712": "/static/rooms/terminus-7-etasje-712.svg",
	"713": "/static/rooms/terminus-7-etasje-713.svg",
	"714": "/static/rooms/terminus-7-etasje-714.svg",
}

// MapPathForRoom returns the Terminus map URL for an exact room number match.
// Rooms without a known map return an empty path and false.
func MapPathForRoom(roomNumber string) (string, bool) {
	path, ok := roomMaps[roomNumber]
	return path, ok
}
