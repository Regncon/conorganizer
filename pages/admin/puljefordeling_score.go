package admin

import (
	"fmt"

	"github.com/Regncon/conorganizer/models"
	smodel "github.com/Regncon/conorganizer/service/puljefordeling/solver/model"
)

type puljeScoreLine struct {
	Label  string
	Points int
}

// puljeScoreLines explains a solver score as the parts it adds up from: the
// interest band, then each fairness bump that applied to the seat.
func puljeScoreLines(score smodel.ScoreBreakdown) []puljeScoreLine {
	band := models.InterestLevelFromScore(int(score.Score)).Label()
	if score.Score == smodel.MaxScore {
		if score.Satisfied {
			band += ", har fått førstevalg"
		} else {
			band += ", mangler førstevalg"
		}
	}
	lines := []puljeScoreLine{{Label: band, Points: score.Band}}
	if score.MissBonus > 0 {
		label := fmt.Sprintf("Ikke fått førstevalg i %s", antallTekst(score.Misses, "pulje", "puljer"))
		lines = append(lines, puljeScoreLine{Label: label, Points: score.MissBonus})
	}
	if score.NeverSeatedBump > 0 {
		lines = append(lines, puljeScoreLine{Label: "Ikke fått plass ennå", Points: score.NeverSeatedBump})
	}
	if score.DMBump > 0 {
		lines = append(lines, puljeScoreLine{Label: "Spilleder i helgen", Points: score.DMBump})
	}
	return lines
}
