package rooms

var VenueData = []FloorSVG{
	{
		Floor: 0,
		Label: "Kjeller",
		Shape: "M0 545 H630 V545 H2000 V980 H0 Z",
	},
	{
		Floor: 1,
		Label: "Inngang",
		Shape: "M0 0 H700 V600 H2000 V1125 H0 Z",
		Rooms: []RoomSVG{
			{
				ID:       "101",
				RoomType: Room,
				Origin: Coordinate{
					X: 0, Y: 0,
				},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 0, Y: 430},
					{X: 235, Y: 430},
					{X: 235, Y: 0},
				},
			},
			{
				ID:       "102",
				RoomType: Room,
				Origin: Coordinate{
					X: 250, Y: 0,
				},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 0, Y: 430},
					{X: 235, Y: 430},
					{X: 235, Y: 0},
				},
			},
		},
	},
}
