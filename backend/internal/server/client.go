package server

import (
	"encoding/json"
	"log"
	"daifugo-backend/models"
	"daifugo-backend/logic"

	"github.com/gorilla/websocket"
)

type Client struct {
	manager *GameManager
	conn    *websocket.Conn
	send    chan []byte
	room    *models.Room
	player  *models.Player
}

func NewClient(manager *GameManager, conn *websocket.Conn) *Client {
	return &Client{
		manager: manager,
		conn:    conn,
		send:    make(chan []byte, 256),
	}
}

type Message struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

func (c *Client) ReadPump() {
	defer func() {
		if c.room != nil && c.player != nil {
			c.room.RemovePlayer(c.player.ID)
			c.BroadcastRoomState()
		}
		c.manager.UnregisterClient(c)
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.handleMessage(message)
	}
}

func (c *Client) WritePump() {
	defer c.conn.Close()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(msgBytes []byte) {
	var msg Message
	if err := json.Unmarshal(msgBytes, &msg); err != nil {
		log.Printf("json err: %v", err)
		return
	}

	switch msg.Type {
	case "join_game":
		roomID, _ := msg.Data["room_id"].(string)
		playerName, _ := msg.Data["player_name"].(string)
		if roomID == "" {
			roomID = "lobby"
		}
		if playerName == "" {
			playerName = "Guest"
		}

		c.room = c.manager.GetOrCreateRoom(roomID)
		
		playerID := playerName + "_" + c.conn.RemoteAddr().String()
		
		c.room.Mu.Lock()
		isHost := len(c.room.Players) == 0
		c.room.Mu.Unlock()
		
		c.player = models.NewPlayer(playerID, playerName, isHost)
		
		c.room.AddPlayer(c.player)
		c.manager.RegisterClient(c, c.room.ID)
		
		log.Printf("Player %s joined %s", playerName, roomID)
		c.BroadcastRoomState()

	case "start_game":
		if c.room != nil && c.player != nil {
			if !c.player.IsHost {
				log.Println("StartGame blocked: not host")
				return
			}
			err := logic.StartGame(c.room)
			if err != nil {
				log.Println("StartGame error:", err)
				return
			}
			c.BroadcastRoomState()
		}

	case "play_cards":
		if c.room != nil && c.player != nil && msg.Data != nil {
			if cardsInterface, ok := msg.Data["cards"].([]interface{}); ok {
				var playCards []models.Card
				for _, ci := range cardsInterface {
					if cMap, ok := ci.(map[string]interface{}); ok {
						playCards = append(playCards, models.Card{
							Suit: models.Suit(cMap["suit"].(string)),
							Rank: int(cMap["rank"].(float64)),
						})
					}
				}
				
				err := logic.PlayTurn(c.room, c.player.ID, playCards)
				if err != nil {
					c.conn.WriteJSON(map[string]interface{}{"type": "error", "message": err.Error()})
					return
				}
				c.manager.BroadcastToRoom(c.room.ID)
			}
		}

	case "resolve_interaction":
		if c.room != nil && c.player != nil && msg.Data != nil {
			err := logic.ResolveInteraction(c.room, c.player.ID, msg.Data)
			if err != nil {
				c.conn.WriteJSON(map[string]interface{}{"type": "error", "message": err.Error()})
				return
			}
			c.manager.BroadcastToRoom(c.room.ID)
		}

	case "pass_turn":
		if c.room != nil && c.player != nil {
			err := logic.PlayTurn(c.room, c.player.ID, []models.Card{})
			if err != nil {
				log.Println("Pass error:", err)
				return
			}
			c.BroadcastRoomState()
		}
	}
}

func (c *Client) BroadcastRoomState() {
	if c.room == nil {
		return
	}
	c.manager.BroadcastToRoom(c.room.ID)
}
