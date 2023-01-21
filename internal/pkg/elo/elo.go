package elo

import (
	"math"

	pkggame "github.com/renju24/backend/pkg/game"
)

// Calculate ...
func Calculate(blackUserRating, whiteUserRating int, winner pkggame.Color) (int, int) {
	var score float64
	switch winner {
	case pkggame.Black:
		score = 1
	case pkggame.White:
		score = 0
	default:
		score = 0.5
	}
	expectedScore := 1 / (1 + math.Pow(10, float64(whiteUserRating-blackUserRating)/float64(400)))
	delta := int(float64(32) * (score - expectedScore))
	return blackUserRating + delta, whiteUserRating - delta
}
