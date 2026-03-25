package models

type Play struct {
	PlayerID string `json:"player_id"`
	Cards    []Card `json:"cards"`
}

type GameState struct {
	IsRevolution  bool     `json:"is_revolution"`
	CurrentTurn   string   `json:"current_turn"`    // PlayerID whose turn it is
	Board         []Play   `json:"board"`           // History of plays in the current trick
	PassedPlayers []string `json:"passed_players"`  // PlayerIDs who passed in the current trick
	PlayerOrder   []string `json:"player_order"`    // PlayerIDs in turn order
}

func NewGameState() *GameState {
	return &GameState{
		IsRevolution:  false,
		Board:         make([]Play, 0),
		PassedPlayers: make([]string, 0),
		PlayerOrder:   make([]string, 0),
	}
}
