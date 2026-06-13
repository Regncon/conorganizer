package rooms

var VenueData = []FloorSVG{
	{
		Floor: 0,
		Label: "Kjeller",
		Shape: "M0 545 H630 V545 H2000 V980 H0 Z",
	},
	{
		Floor: 7,
		Label: "Spillerom",
		Shape: "M0 0 H750 V875 H2000 V1500 H0 Z",
		Rooms: []RoomSVG{
			{
				ID:       "S3",
				RoomType: Stairs,
				Origin: Coordinate{
					X: 1875, Y: 875,
				},
				Center: Coordinate{X: 60, Y: 130},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 125, Y: 0},
					{X: 125, Y: 260},
					{X: 0, Y: 260},
				},
			},
			{
				ID:       "T2",
				RoomType: Toilet,
				Origin: Coordinate{
					X: 1725, Y: 875,
				},
				Center: Coordinate{X: 70, Y: 130},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 150, Y: 0},
					{X: 150, Y: 260},
					{X: 0, Y: 260},
				},
			},
			{
				ID:       "714",
				RoomType: Room,
				Origin: Coordinate{
					X: 1210, Y: 875,
				},
				Center: Coordinate{X: 258, Y: 130},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 515, Y: 0},
					{X: 515, Y: 260},
					{X: 0, Y: 260},
				},
			},
			{
				ID:       "715",
				RoomType: Room,
				Origin: Coordinate{
					X: 1060, Y: 875,
				},
				Center: Coordinate{X: 75, Y: 130},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 150, Y: 0},
					{X: 150, Y: 260},
					{X: 50, Y: 260},
					{X: 50, Y: 160},
					{X: 0, Y: 160},
				},
			},
			{
				ID:       "U2",
				RoomType: Utility,
				Origin: Coordinate{
					X: 1010, Y: 1035,
				},
				Center: Coordinate{X: 50, Y: 50},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 100, Y: 0},
					{X: 100, Y: 100},
					{X: 0, Y: 100},
				},
			},
			{
				ID:       "716",
				RoomType: Room,
				Origin: Coordinate{
					X: 910, Y: 875,
				},
				Center: Coordinate{X: 75, Y: 130},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 150, Y: 0},
					{X: 150, Y: 160},
					{X: 100, Y: 160},
					{X: 100, Y: 260},
					{X: 0, Y: 260},
				},
			},
			{
				ID:       "S2",
				RoomType: Stairs,
				Origin: Coordinate{
					X: 450, Y: 875,
				},
				Center: Coordinate{X: 230, Y: 130},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 460, Y: 0},
					{X: 460, Y: 260},
					{X: 0, Y: 260},
				},
			},
			{
				ID:       "E1",
				RoomType: Elevator,
				Origin: Coordinate{
					X: 450, Y: 785,
				},
				Center: Coordinate{X: 45, Y: 45},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 90, Y: 0},
					{X: 90, Y: 90},
					{X: 0, Y: 90},
				},
			},
			{
				ID:       "E2",
				RoomType: Elevator,
				Origin: Coordinate{
					X: 540, Y: 785,
				},
				Center: Coordinate{X: 45, Y: 45},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 90, Y: 0},
					{X: 90, Y: 90},
					{X: 0, Y: 90},
				},
			},
			{
				ID:       "U1",
				RoomType: Utility,
				Origin: Coordinate{
					X: 450, Y: 655,
				},
				Center: Coordinate{X: 150, Y: 85},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 300, Y: 0},
					{X: 300, Y: 220},
					{X: 180, Y: 220},
					{X: 180, Y: 130},
					{X: 0, Y: 130},
				},
			},
			{
				ID:       "705",
				RoomType: Room,
				Origin: Coordinate{
					X: 450, Y: 480,
				},
				Center: Coordinate{X: 150, Y: 85},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 300, Y: 0},
					{X: 300, Y: 175},
					{X: 0, Y: 175},
				},
			},
			{
				ID:       "706",
				RoomType: Room,
				Origin: Coordinate{
					X: 450, Y: 120,
				},
				Center: Coordinate{X: 150, Y: 180},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 300, Y: 0},
					{X: 300, Y: 360},
					{X: 0, Y: 360},
				},
			},
			{
				ID:       "S1",
				RoomType: Stairs,
				Origin: Coordinate{
					X: 450, Y: 0,
				},
				Center: Coordinate{X: 150, Y: 60},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 300, Y: 0},
					{X: 300, Y: 120},
					{X: 0, Y: 120},
				},
			},
			{
				ID:       "707",
				RoomType: Room,
				Origin: Coordinate{
					X: 0, Y: 0,
				},
				Center: Coordinate{X: 170, Y: 300},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 340, Y: 0},
					{X: 340, Y: 600},
					{X: 0, Y: 600},
				},
			},
			{
				ID:       "T1",
				RoomType: Toilet,
				Origin: Coordinate{
					X: 0, Y: 590,
				},
				Center: Coordinate{X: 170, Y: 105},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 340, Y: 0},
					{X: 340, Y: 210},
					{X: 0, Y: 210},
				},
			},
			{
				ID:       "709",
				RoomType: Room,
				Origin: Coordinate{
					X: 0, Y: 800,
				},
				Center: Coordinate{X: 170, Y: 350},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 340, Y: 0},
					{X: 340, Y: 700},
					{X: 0, Y: 700},
				},
			},
			{
				ID:       "710",
				RoomType: Room,
				Origin: Coordinate{
					X: 340, Y: 1240,
				},
				Center: Coordinate{X: 150, Y: 130},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 300, Y: 0},
					{X: 300, Y: 260},
					{X: 0, Y: 260},
				},
			},
			{
				ID:       "711",
				RoomType: Room,
				Origin: Coordinate{
					X: 640, Y: 1240,
				},
				Center: Coordinate{X: 250, Y: 130},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 500, Y: 0},
					{X: 500, Y: 260},
					{X: 0, Y: 260},
				},
			},
			{
				ID:       "712",
				RoomType: Room,
				Origin: Coordinate{
					X: 1140, Y: 1240,
				},
				Center: Coordinate{X: 250, Y: 130},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 500, Y: 0},
					{X: 500, Y: 260},
					{X: 0, Y: 260},
				},
			},
			{
				ID:       "713",
				RoomType: Room,
				Origin: Coordinate{
					X: 1640, Y: 1240,
				},
				Center: Coordinate{X: 180, Y: 130},
				Coordinates: []Coordinate{
					{X: 0, Y: 0},
					{X: 360, Y: 0},
					{X: 360, Y: 260},
					{X: 0, Y: 260},
				},
			},
		},
	},
}
