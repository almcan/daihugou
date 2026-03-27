package logic

import (
	"daifugo-backend/models"
	"errors"
)

func ResolveInteraction(room *models.Room, playerID string, payload map[string]interface{}) error {
	room.Mu.Lock()
	defer room.Mu.Unlock()

	game := room.Game
	if game == nil || game.Interaction == nil || game.Interaction.PlayerID != playerID {
		return errors.New("no interaction pending for player")
	}

	player := room.Players[playerID]

	switch game.Interaction.Type {
	case "7_pass":
		cardsRaw := extractCardsRaw(payload["cards"])
		if len(cardsRaw) == 0 { return errors.New("missing cards") }
		
		targetID := getNextPlayer(game, room)
		targetPlayer, exists := room.Players[targetID]

		for _, cardMap := range cardsRaw {
			suit := models.Suit(cardMap["suit"].(string))
			rank := int(cardMap["rank"].(float64))
			removeCard(player, suit, rank)
			if exists {
				targetPlayer.Hand = append(targetPlayer.Hand, models.Card{Suit: suit, Rank: rank})
			}
		}
		AssignRankIfWon(room, player)
		if exists { SortHand(targetPlayer.Hand, game.IsRevolution) }

	case "9_pick":
		cardsRaw := extractCardsRaw(payload["cards"])
		for _, cardMap := range cardsRaw {
			suit := models.Suit(cardMap["suit"].(string))
			rank := int(cardMap["rank"].(float64))
			// remove from cemetery and give to player
			game.Cemetery = removeCardFromSlice(game.Cemetery, suit, rank)
			player.Hand = append(player.Hand, models.Card{Suit: suit, Rank: rank})
		}
		SortHand(player.Hand, game.IsRevolution)

	case "10_discard":
		cardsRaw := extractCardsRaw(payload["cards"])
		for _, cardMap := range cardsRaw {
			suit := models.Suit(cardMap["suit"].(string))
			rank := int(cardMap["rank"].(float64))
			removeCard(player, suit, rank)
			game.Cemetery = append(game.Cemetery, models.Card{Suit: suit, Rank: rank})
		}
		AssignRankIfWon(room, player)

	case "12_bomber":
		ranksRaw, ok := payload["ranks"].([]interface{})
		if !ok { return errors.New("invalid ranks payload") }
		
		for _, rRaw := range ranksRaw {
			targetRank := int(rRaw.(float64))
			for _, p := range room.Players {
				newHand := make([]models.Card, 0)
				for _, c := range p.Hand {
					if c.Rank != targetRank {
						newHand = append(newHand, c)
					} else {
						game.Cemetery = append(game.Cemetery, c)
					}
				}
				p.Hand = newHand
			}
		}
		for _, p := range room.Players {
			AssignRankIfWon(room, p)
		}
	}

	game.Interaction = nil
	advanceTurn(game, room)
	return nil
}

func removeCard(player *models.Player, suit models.Suit, rank int) {
	player.Hand = removeCardFromSlice(player.Hand, suit, rank)
}

func removeCardFromSlice(cards []models.Card, suit models.Suit, rank int) []models.Card {
	newCards := make([]models.Card, 0)
	removed := false
	for _, c := range cards {
		if !removed && c.Suit == suit && c.Rank == rank {
			removed = true
			continue
		}
		newCards = append(newCards, c)
	}
	return newCards
}

func extractCardsRaw(raw interface{}) []map[string]interface{} {
	ret := make([]map[string]interface{}, 0)
	if slice, ok := raw.([]interface{}); ok {
		for _, item := range slice {
			if m, ok := item.(map[string]interface{}); ok {
				ret = append(ret, m)
			}
		}
	}
	return ret
}

func getNextPlayer(game *models.GameState, room *models.Room) string {
	currentIdx := -1
	for i, pid := range game.PlayerOrder {
		if pid == game.CurrentTurn {
			currentIdx = i
			break
		}
	}
	if currentIdx != -1 {
		nextIdx := (currentIdx + 1) % len(game.PlayerOrder)
		for nextIdx != currentIdx {
			pid := game.PlayerOrder[nextIdx]
			
			hasPassed := false
			for _, passedPID := range game.PassedPlayers {
				if passedPID == pid {
					hasPassed = true; break
				}
			}
			hasCards := len(room.Players[pid].Hand) > 0

			if !hasPassed && hasCards {
				return pid
			}
			nextIdx = (nextIdx + 1) % len(game.PlayerOrder)
		}
	}
	return ""
}
