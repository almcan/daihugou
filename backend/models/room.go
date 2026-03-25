package models

import "sync"

type Room struct {
	ID      string
	Players map[string]*Player
	State   string
	Game    *GameState
	Mu      sync.Mutex
}

func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		Players: make(map[string]*Player),
		State:   "waiting",
		Game:    NewGameState(),
	}
}

func (r *Room) AddPlayer(p *Player) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.Players[p.ID] = p
}

func (r *Room) RemovePlayer(playerID string) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	delete(r.Players, playerID)
}

func (r *Room) ToMap(clientPlayerID string) map[string]interface{} {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	playersData := make([]map[string]interface{}, 0, len(r.Players))
	for _, p := range r.Players {
		playersData = append(playersData, p.ToMap(p.ID == clientPlayerID))
	}

	return map[string]interface{}{
		"room_id": r.ID,
		"state":   r.State,
		"players": playersData,
		"game":    r.Game,
	}
}
