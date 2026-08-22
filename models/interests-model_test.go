package models

import "testing"

// Score and InterestLevelFromScore read the same table, so every defined level
// must survive a round trip through its score and back.
func TestInterestLevelScoreRoundTrip(t *testing.T) {
	for _, level := range []InterestLevel{InterestLevelHigh, InterestLevelMedium, InterestLevelLow} {
		score := level.Score()
		if score == 0 {
			t.Errorf("%q should map to a non-zero score", level)
		}
		if got := InterestLevelFromScore(score); got != level {
			t.Errorf("round trip: %q -> %d -> %q", level, score, got)
		}
	}
}

func TestInterestLevelFromScoreUnknown(t *testing.T) {
	for _, score := range []int{0, 2, 4, 6, -1} {
		if got := InterestLevelFromScore(score); got != InterestLevelNone {
			t.Errorf("score %d should map to None, got %q", score, got)
		}
	}
	if InterestLevelNone.Score() != 0 {
		t.Errorf("None should score 0, got %d", InterestLevelNone.Score())
	}
}
