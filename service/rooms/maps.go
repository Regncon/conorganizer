package rooms

var roomMaps = map[string]string{
	"001": "/static/rooms/terminus-0-etasje-001.svg",
	"002": "/static/rooms/terminus-0-etasje-002.svg",
	"003": "/static/rooms/terminus-0-etasje-003.svg",
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

var floorMaps = map[int]string{
	0: "/static/rooms/terminus-0-etasje.svg",
	7: "/static/rooms/terminus-7-etasje.svg",
}

// MapPathForRoom returns the Terminus map URL for an exact room number match.
// Rooms without a known map return an empty path and false.
func MapPathForRoom(roomNumber string) (string, bool) {
	path, ok := roomMaps[roomNumber]
	return path, ok
}

// MapPathForFloor returns the overview map URL for a floor with an interactive map.
func MapPathForFloor(floor int) (string, bool) {
	path, ok := floorMaps[floor]
	return path, ok
}
