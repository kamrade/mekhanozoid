package game

import (
	"math/rand"
)

func ShuffleCards(cards []CardInstance, seed int64) {
	r := rand.New(rand.NewSource(seed))

	r.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
}
