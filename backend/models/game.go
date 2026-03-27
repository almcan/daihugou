package models

type Play struct {
	PlayerID string `json:"player_id"`
	Cards    []Card `json:"cards"`
}

type InteractionState struct {
	PlayerID string `json:"player_id"`
	Type     string `json:"type"`      // "7_pass", "9_nominate", "10_discard", "12_bomber"
	Count    int    `json:"count"`     // Number of cards to select if applicable
}

type GameState struct {
	IsRevolution  bool     `json:"is_revolution"`
	CurrentTurn   string   `json:"current_turn"`    // PlayerID whose turn it is
	Board         []Play   `json:"board"`           // History of plays in the current trick
	PassedPlayers   []string     `json:"passed_players"`
	PlayerOrder     []string          `json:"player_order"`
	BoundSuits      []Suit            `json:"bound_suits"`
	IsSequenceBound bool              `json:"is_sequence_bound"`
	Is11Back        bool              `json:"is_11_back"`
	MemeEvent       string            `json:"meme_event"`
	Interaction     *InteractionState `json:"interaction"`
	Cemetery        []Card            `json:"cemetery"`
}

func NewGameState() *GameState {
	return &GameState{
		IsRevolution:    false,
		Board:           make([]Play, 0),
		PassedPlayers:   make([]string, 0),
		PlayerOrder:     make([]string, 0),
		BoundSuits:      make([]Suit, 0),
		IsSequenceBound: false,
		Is11Back:        false,
		Cemetery:        make([]Card, 0),
	}
}
