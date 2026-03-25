package models

type Suit string

const (
	Spade   Suit = "spade"
	Heart   Suit = "heart"
	Diamond Suit = "diamond"
	Clover  Suit = "clover"
	Joker   Suit = "joker"
)

type Card struct {
	Suit Suit `json:"suit"`
	Rank int  `json:"rank"`
}

func (c Card) Strength(isRevolution bool) int {
	if c.Suit == Joker {
		return 99
	}
	strength := c.Rank
	if c.Rank == 1 {
		strength = 14
	} else if c.Rank == 2 {
		strength = 15
	}

	if isRevolution {
		return 18 - strength
	}
	return strength
}
