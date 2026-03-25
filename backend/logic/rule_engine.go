package logic

import (
	"daifugo-backend/models"
)

func IsValidPlay(playCards []models.Card, board []models.Play, isRevolution bool) bool {
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

	playStrength := getPlayStrength(playCards, isRevolution)
	lastStrength := getPlayStrength(lastPlay.Cards, isRevolution)

	return playStrength > lastStrength
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
