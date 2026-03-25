package logic

import (
	"daifugo-backend/models"
	"math/rand"
	"sort"
)

func GenerateDeck(includeJokers int) []models.Card {
	suits := []models.Suit{models.Spade, models.Heart, models.Diamond, models.Clover}
	var deck []models.Card

	for _, suit := range suits {
		for rank := 1; rank <= 13; rank++ {
			deck = append(deck, models.Card{Suit: suit, Rank: rank})
		}
	}

	for i := 0; i < includeJokers; i++ {
		deck = append(deck, models.Card{Suit: models.Joker, Rank: 0})
	}

	return deck
}

func Shuffle(deck []models.Card) {
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
}

func SortHand(hand []models.Card, isRevolution bool) {
	sort.Slice(hand, func(i, j int) bool {
		strI := hand[i].Strength(isRevolution)
		strJ := hand[j].Strength(isRevolution)
		if strI == strJ {
			return hand[i].Suit < hand[j].Suit
		}
		return strI < strJ
	})
}
