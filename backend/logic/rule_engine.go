package logic

import (
	"daifugo-backend/models"
)

func IsValidPlay(playCards []models.Card, board []models.Play, isRevolution bool, boundSuits []models.Suit, isSequenceBound bool, is11Back bool) bool {
	if len(playCards) == 0 {
		return false
	}

	baseRank := -1
	for _, c := range playCards {
		if c.Suit != models.Joker {
			if baseRank == -1 {
				baseRank = c.Rank
			} else if baseRank != c.Rank {
				return false
			}
		}
	}

	if len(board) == 0 {
		return true
	}

	lastPlay := board[len(board)-1]
	if len(playCards) != len(lastPlay.Cards) {
		return false
	}

	effectiveRevolution := isRevolution != is11Back
	playStrength := getPlayStrength(playCards, effectiveRevolution)
	lastStrength := getPlayStrength(lastPlay.Cards, effectiveRevolution)

	if playStrength <= lastStrength {
		return false
	}

	if isSequenceBound {
		if playStrength != lastStrength+1 {
			return false
		}
	}

	if len(boundSuits) > 0 {
		neededSuits := make(map[models.Suit]bool)
		for _, s := range boundSuits {
			neededSuits[s] = true
		}
		
		jokerCount := 0
		providedSuits := make(map[models.Suit]bool)
		for _, c := range playCards {
			if c.Suit == models.Joker {
				jokerCount++
			} else {
				providedSuits[c.Suit] = true
			}
		}

		for s := range providedSuits {
			if !neededSuits[s] {
				return false
			}
		}

		matched := 0
		for s := range neededSuits {
			if providedSuits[s] {
				matched++
			}
		}
		if matched + jokerCount < len(neededSuits) {
			return false
		}
	}

	return true
}

func getPlayStrength(cards []models.Card, isRevolution bool) int {
	maxStrength := -1
	for _, c := range cards {
		s := c.Strength(isRevolution)
		if c.Suit != models.Joker {
			return s
		}
		if s > maxStrength {
			maxStrength = s
		}
	}
	return maxStrength
}

func GetSuits(cards []models.Card) []models.Suit {
	suits := []models.Suit{}
	for _, c := range cards {
		if c.Suit != models.Joker {
			suits = append(suits, c.Suit)
		}
	}
	return suits
}

func SuitsMatch(s1, s2 []models.Suit) bool {
	if len(s1) != len(s2) {
		return false
	}
	m := make(map[models.Suit]int)
	for _, s := range s1 { m[s]++ }
	for _, s := range s2 { m[s]-- }
	for _, count := range m {
		if count != 0 { return false }
	}
	return true
}
