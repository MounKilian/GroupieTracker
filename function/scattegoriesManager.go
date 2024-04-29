package groupieTracker

import (
	"math/rand/v2"
)

var Letters = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "T", "U", "V", "W", "X", "Y", "Z"}

func selectRandomLetter() string {
	randomIndex := rand.IntN(len(Letters) - 1)
	return Letters[randomIndex]
}
