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
		
		activeCount := 0
		for _, pid := range game.PlayerOrder {
			if len(room.Players[pid].Hand) > 0 {
				activeCount++
			}
		}

		if len(game.PassedPlayers) >= activeCount - 1 && len(game.Board) > 0 {
			lastPlayerID := game.Board[len(game.Board)-1].PlayerID
			clearBoard(game)
			game.CurrentTurn = lastPlayerID
			
			if len(room.Players[lastPlayerID].Hand) == 0 {
				advanceTurn(game, room)
			}
		} else {
			advanceTurn(game, room)
		}
		return nil
	}

	if !IsValidPlay(cards, game.Board, game.IsRevolution, game.BoundSuits, game.IsSequenceBound, game.Is11Back) {
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

	AssignRankIfWon(room, player)

	if len(game.Board) > 0 && len(game.BoundSuits) == 0 {
		lastPlay := game.Board[len(game.Board)-1]
		
		lastStr := getPlayStrength(lastPlay.Cards, game.IsRevolution)
		currStr := getPlayStrength(cards, game.IsRevolution)
		if currStr == lastStr+1 {
			game.IsSequenceBound = true
		}

		lastSuits := GetSuits(lastPlay.Cards)
		currSuits := GetSuits(cards)
		if len(game.BoundSuits) == 0 && len(lastSuits) > 0 && SuitsMatch(lastSuits, currSuits) {
			game.BoundSuits = lastSuits
		}
	}

	game.Board = append(game.Board, models.Play{PlayerID: playerID, Cards: cards})

	var memeEvents []string

	baseRank := -1
	for _, c := range cards {
		if c.Suit != models.Joker {
			baseRank = c.Rank
			break
		}
	}

	if len(cards) >= 4 {
		game.IsRevolution = !game.IsRevolution
		memeEvents = append(memeEvents, "革命発動！")
	}

	if baseRank == 8 {
		memeEvents = append(memeEvents, "8切り発動！")
		game.MemeEvent = joinMemes(memeEvents)
		clearBoard(game)
		game.CurrentTurn = playerID
		return nil
	}

	if baseRank == 11 {
		game.Is11Back = true
		memeEvents = append(memeEvents, "11バック発動！")
	}

	skipCount := 0
	if baseRank == 5 {
		skipCount = len(cards)
		memeEvents = append(memeEvents, "5飛び発動！")
	}
	
	if game.IsSequenceBound && len(game.Board) > 1 {
		memeEvents = append(memeEvents, "連番縛り発動！")
	} else if len(game.BoundSuits) > 0 && len(game.Board) > 1 {
		// Only trigger the meme on the play that *caused* the bind. 
		// If last trick length was 2, it just got bound.
		if len(game.Board) == 2 {
			memeEvents = append(memeEvents, "縛り発動！")
		}
	}

	game.MemeEvent = joinMemes(memeEvents)

	if baseRank == 7 {
		game.Interaction = &models.InteractionState{PlayerID: playerID, Type: "7_pass", Count: len(cards)}
		return nil
	}
	if baseRank == 9 {
		game.Interaction = &models.InteractionState{PlayerID: playerID, Type: "9_pick", Count: len(cards)}
		return nil
	}
	if baseRank == 10 {
		game.Interaction = &models.InteractionState{PlayerID: playerID, Type: "10_discard", Count: len(cards)}
		return nil
	}
	if baseRank == 12 {
		game.Interaction = &models.InteractionState{PlayerID: playerID, Type: "12_bomber", Count: len(cards)}
		return nil
	}

	advanceTurn(game, room, skipCount)
	
	otherActiveCount := 0
	for _, pid := range game.PlayerOrder {
		if pid == playerID {
			continue
		}
		hasPassed := false
		for _, passedPID := range game.PassedPlayers {
			if passedPID == pid {
				hasPassed = true
				break
			}
		}
		if !hasPassed && len(room.Players[pid].Hand) > 0 {
			otherActiveCount++
		}
	}

	if skipCount > 0 && skipCount >= otherActiveCount {
		clearBoard(game)
		game.CurrentTurn = playerID
		if len(room.Players[playerID].Hand) == 0 {
			advanceTurn(game, room)
		}
	}
	return nil
}

func clearBoard(game *models.GameState) {
	for _, play := range game.Board {
		game.Cemetery = append(game.Cemetery, play.Cards...)
	}
	game.Board = make([]models.Play, 0)
	game.PassedPlayers = make([]string, 0)
	game.BoundSuits = make([]models.Suit, 0)
	game.IsSequenceBound = false
	game.Is11Back = false
}

func joinMemes(memes []string) string {
	res := ""
	for i, m := range memes {
		if i > 0 {
			res += " & "
		}
		res += m
	}
	return res
}

func advanceTurn(game *models.GameState, room *models.Room, skipCount ...int) {
	skips := 0
	if len(skipCount) > 0 {
		skips = skipCount[0]
	}

	currentIdx := -1
	for i, pid := range game.PlayerOrder {
		if pid == game.CurrentTurn {
			currentIdx = i
			break
		}
	}

	if currentIdx != -1 {
		nextIdx := (currentIdx + 1) % len(game.PlayerOrder)
		foundValid := 0
		for i := 0; i < len(game.PlayerOrder)*2; i++ {
			pid := game.PlayerOrder[nextIdx]
			
			hasPassed := false
			for _, passedPID := range game.PassedPlayers {
				if passedPID == pid {
					hasPassed = true; break
				}
			}
			hasCards := len(room.Players[pid].Hand) > 0

			if !hasPassed && hasCards {
				if foundValid == skips {
					game.CurrentTurn = pid
					return
				}
				foundValid++
			}
			nextIdx = (nextIdx + 1) % len(game.PlayerOrder)
		}
	}
}

func AssignRankIfWon(room *models.Room, player *models.Player) {
	if len(player.Hand) == 0 && player.Rank == 0 {
		maxRank := 0
		for _, p := range room.Players {
			if p.Rank > maxRank {
				maxRank = p.Rank
			}
		}
		player.Rank = maxRank + 1
		
		activePlayers := 0
		for _, p := range room.Players {
			if len(p.Hand) > 0 {
				activePlayers++
			}
		}
		if activePlayers <= 1 {
			room.State = "waiting"
			for _, p := range room.Players {
				if len(p.Hand) > 0 && p.Rank == 0 {
					p.Rank = maxRank + 2
				}
			}
		}
	}
}


