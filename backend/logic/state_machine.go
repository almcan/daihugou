package logic

import (
	"daifugo-backend/models"
	"errors"
)

func StartGame(room *models.Room) error {
	room.Mu.Lock()
	defer room.Mu.Unlock()

	if len(room.Players) < 2 {
		return errors.New("not enough players")
	}

	room.State = "playing"
	room.Game = models.NewGameState()
	
	playerIDs := make([]string, 0, len(room.Players))
	for id := range room.Players {
		playerIDs = append(playerIDs, id)
	}
	room.Game.PlayerOrder = playerIDs
	room.Game.CurrentTurn = playerIDs[0]

	deck := GenerateDeck(1)
	Shuffle(deck)

	pIdx := 0
	for _, card := range deck {
		pid := playerIDs[pIdx]
		room.Players[pid].Hand = append(room.Players[pid].Hand, card)
		pIdx = (pIdx + 1) % len(playerIDs)
	}

	for _, pid := range playerIDs {
		SortHand(room.Players[pid].Hand, room.Game.IsRevolution)
	}

	return nil
}

func PlayTurn(room *models.Room, playerID string, cards []models.Card) error {
	room.Mu.Lock()
	defer room.Mu.Unlock()

	if room.State != "playing" {
		return errors.New("game not in progress")
	}

	game := room.Game
	if game.CurrentTurn != playerID {
		return errors.New("not your turn")
	}

	if len(cards) == 0 {
		game.PassedPlayers = append(game.PassedPlayers, playerID)
		advanceTurn(game)
		return nil
	}

	if !IsValidPlay(cards, game.Board, game.IsRevolution) {
		return errors.New("invalid play")
	}

	player := room.Players[playerID]
	newHand := make([]models.Card, 0)
	for _, cHand := range player.Hand {
		removed := false
		for _, cPlay := range cards {
			if cHand.Suit == cPlay.Suit && cHand.Rank == cPlay.Rank {
				removed = true
				break
			}
		}
		if !removed {
			newHand = append(newHand, cHand)
		}
	}
	player.Hand = newHand

	game.Board = append(game.Board, models.Play{PlayerID: playerID, Cards: cards})
	advanceTurn(game)
	return nil
}

func advanceTurn(game *models.GameState) {
	currentIdx := -1
	for i, pid := range game.PlayerOrder {
		if pid == game.CurrentTurn {
			currentIdx = i
			break
		}
	}

	if currentIdx != -1 {
		game.CurrentTurn = game.PlayerOrder[(currentIdx+1)%len(game.PlayerOrder)]
	}
}

