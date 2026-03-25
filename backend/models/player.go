package models

type Player struct {
	ID       string
	Name     string
	Hand     []Card
	IsActive bool
	Rank     int
	IsHost   bool
}

func NewPlayer(id, name string, isHost bool) *Player {
	return &Player{
		ID:       id,
		Name:     name,
		Hand:     make([]Card, 0),
		IsActive: true,
		Rank:     0,
		IsHost:   isHost,
	}
}

func (p *Player) ToMap(includeHand bool) map[string]interface{} {
	data := map[string]interface{}{
		"id":         p.ID,
		"name":       p.Name,
		"card_count": len(p.Hand),
		"is_active":  p.IsActive,
		"rank":       p.Rank,
		"is_host":    p.IsHost,
	}
	if includeHand {
		data["hand"] = p.Hand
	}
	return data
}
