package server

import (
	"daifugo-backend/models"
	"encoding/json"
	"sync"
)

type GameManager struct {
	Rooms map[string]*models.Room
	Clients map[*Client]string // Client -> roomID mapping
	mu    sync.Mutex
}

func NewGameManager() *GameManager {
	return &GameManager{
		Rooms: make(map[string]*models.Room),
		Clients: make(map[*Client]string),
	}
}

func (gm *GameManager) GetOrCreateRoom(roomID string) *models.Room {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	room, exists := gm.Rooms[roomID]
	if !exists {
		room = models.NewRoom(roomID)
		gm.Rooms[roomID] = room
	}
	return room
}

func (gm *GameManager) RegisterClient(c *Client, roomID string) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	gm.Clients[c] = roomID
}

func (gm *GameManager) UnregisterClient(c *Client) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	delete(gm.Clients, c)
}

func (gm *GameManager) BroadcastToRoom(roomID string) {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	room, exists := gm.Rooms[roomID]
	if !exists { return }
	for client, cRoomID := range gm.Clients {
		if cRoomID == roomID {
			playerID := ""
			if client.player != nil {
				playerID = client.player.ID
			}
			stateBytes, _ := json.Marshal(Message{
				Type: "state_update",
				Data: room.ToMap(playerID),
			})
			select {
			case client.send <- stateBytes:
			default:
				close(client.send)
				delete(gm.Clients, client)
			}
		}
	}
}
