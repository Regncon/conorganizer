package testutil

import (
	"math/rand"
)

type Person struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func GenerateFakePerson() Person {
	var firstNames = []string{
		"Nora", "Emma", "Olivia", "Sofie", "Ada", "Ella", "Frida", "Leah", "Isabella", "Maja",
		"Elias", "Lucas", "Noah", "Oliver", "Filip", "William", "Jakob", "Emil", "Oskar", "Mathias",
		"Aksel", "Thea", "Emilia", "Hanna", "Ingrid", "Astrid", "Sara", "Julie", "Anna", "Karoline",
		"Mia", "Linnea", "Victoria", "Ida", "Oline", "Amalie", "Alma", "Victoria", "Sigrid", "Eline",
		"Jonas", "Eirik", "Magnus", "Kristian", "Henrik", "Marius", "Andreas", "Simen", "Alexander", "Markus",
		"Sander", "Sebastian", "Isak", "Daniel", "Benjamin", "Nikolai", "Iselin", "Tiril", "Tuva", "Sigurd",
		"Leander", "Sigve", "Peder", "Håkon", "Even", "Fredrik", "Vetle", "Mathea", "Signe", "Julie",
		"Ella", "Olav", "Sondre", "Julius", "Kasper", "Leo", "Aloysius", "Hartvig", "Trym", "Stian",
		"Elise", "Karen", "Vilde", "Stella", "Alva", "Eira", "Hedda", "Nelly", "Ylva", "Sylvia",
		"Inga", "Hilde", "Sofie", "Victoria", "Maren", "Lea", "Liv", "Maren", "Birgitte", "Åshild",
	}

	var lastNames = []string{
		"Hansen", "Johansen", "Olsen", "Larsen", "Andersen", "Nielsen", "Pedersen", "Kristiansen", "Jensen", "Karlsen",
		"Jacobsen", "Iversen", "Haugen", "Moe", "Hagen", "Nygård", "Foss", "Lund", "Berg", "Solberg",
		"Dahl", "Lie", "Bakken", "Karlsen", "Eide", "Halvorsen", "Aas", "Lien", "Amundsen", "Sand",
		"Moen", "Rasmussen", "Holm", "Hansen", "Solheim", "Lundberg", "Rønning", "Knudsen", "Skogen", "Vik",
		"Bergen", "Svendsen", "Myhre", "Sørensen", "Lorentzen", "Sætre", "Bråthen", "Eriksen", "Thomassen", "Wold",
		"Lyng", "Gundersen", "Østby", "Valen", "Bråten", "Huse", "Folden", "Lian", "Sund", "Kleven",
		"Opdal", "Bjørnstad", "Meland", "Strand", "Viken", "Stensrud", "Bø", "Løken", "Flåten", "Stenberg",
		"Skjælaaen", "Roald", "Vik", "Holmgren", "Holt", "Kolstad", "Langli", "Vestby", "Grønn", "Karset",
		"Rindal", "Skålvik", "Bilstad", "Øverland", "Løland", "Brevik", "Korsvik", "Fjeld", "Åsheim", "Skogen",
		"Haavik", "Grimsrud", "Vatne", "Kjølberg", "Sæther", "Breivik", "Lunde", "Ødegård", "Dalsgaard", "Bakke",
	}

	var domains = []string{
		"regncon.no", "gmail.com", "hotmail.com", "msn.com", "hotsinglesinyour.area",
	}

	// Start generating something random
	firstName := firstNames[rand.Intn(len(firstNames))]
	lastName := lastNames[rand.Intn(len(lastNames))]
	email := firstName + lastName + "@" + domains[rand.Intn(len(domains))]

	return Person{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
}
