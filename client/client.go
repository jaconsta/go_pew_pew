package main

import (
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)

type Client struct {
	conn     *websocket.Conn
	handlers map[string]EventHandler
	gameData *Game
}

func newClient() (*Client, error) {
	origin := "http://localhost:8000"
	user := NewUser()
	url := fmt.Sprintf("ws://localhost:8000/ws?username=%s", user.username)
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		return nil, err
	}

	c := Client{
		conn:     ws,
		handlers: newEventHandlers(),
		gameData: NewGameData(&user),
	}

	return &c, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) routeEvent(eventData Event) error {
	op, ok := c.handlers[eventData.Type]
	if !ok {
		return fmt.Errorf("No handler registered for event %s", eventData.Type)
	}
	return op(eventData, c)
}

func (c *Client) sendEvent(event Event) {
	payload, err := json.Marshal(event)
	if err != nil {
		log.Println("outSelectRoom error ")
		log.Println(err)
	}

	c.conn.Write(payload)
}
