package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type ClientLists map[*Client]bool

type GameData struct {
	user *User
	room string
	life int
}

// Web sockets client
type Client struct {
	connection *websocket.Conn
	manager    *Manager
	// Ensure single write to websockets
	egress   chan Event
	gameData *GameData
}

// Ping pong config
var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

var (
	openRooms = []string{"lobby", "jungle", "tundra"}
)

func NewWsClient(c *websocket.Conn, m *Manager, u *User) *Client {
	client := Client{
		connection: c,
		manager:    m,
		egress:     make(chan Event),
		gameData: &GameData{
			user: u,
			room: openRooms[0],
			life: 5,
		},
	}
	return &client
}

func (c *Client) readMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	c.connection.SetReadLimit(wsMsgMaxSizeBytes)
	// Conf pong automatic receive
	c.connection.SetReadDeadline(time.Now().Add(pongWait)) // Should not error since it is only time
	c.connection.SetPongHandler(c.pongHandler)

	for {
		mtype, message, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure) {
				log.Println("error read:", err)
			}
			break
		}

		// Other types not supported (Ping/Pong is above)
		if mtype != websocket.TextMessage {
			log.Println("Received not supported messageType:", mtype)
			continue
		}
		c.callEventRouter(message)
	}
}

func (c *Client) writeMessages() {
	// Ping ticker. Pong is in readMessages
	pingTicker := time.NewTicker(pingInterval)
	defer func() {
		pingTicker.Stop()
		// Graceful close connection
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			// Send a message
			if !ok {
				// Tell client that c.manager closed the connection
				c.connection.WriteMessage(websocket.CloseMessage, nil) // Can error but should not be critical
				return
			}
			data, err := json.Marshal(message)
			if err != nil {
				log.Println("Write message, error marshalling: ", err)
				continue
			}
			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("Error sending message: ", err)
			}

		case <-pingTicker.C:
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("Ping error: ", err)
				return
			}
		}
	}
}

// Resets the timer every time a PongMessage is received
func (c *Client) pongHandler(pong string) error {
	// Current time + Pong wait time
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

func (c *Client) callEventRouter(data []byte) {
	var eventData Event
	if err := json.Unmarshal(data, &eventData); err != nil {
		log.Println("Read message. Received bad message: ", err)
	}
	if err := c.manager.routeEvent(eventData, c); err != nil {
		log.Println("Read message, error handling: ", err)
	}

}
