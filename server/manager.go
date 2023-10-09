package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	// 	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var websocketUpgrader = websocket.Upgrader{
	CheckOrigin:     checkOrigin,
	ReadBufferSize:  128,
	WriteBufferSize: 128,
}

func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	// Should change to something like allowed_hosts
	return origin == host
}

type Manager struct {
	clients ClientLists
	sync.RWMutex

	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := Manager{
		clients:  make(ClientLists),
		handlers: map[string]EventHandler{},
	}

	m.NewEventHandlers()
	return &m
}

// Setups the Incomming events
func (m *Manager) NewEventHandlers() {
	m.handlers[ChangeRoomEvent] = ChangeRoomHandler
	m.handlers[ImpactSentEvent] = ImpactSentHandler
}

func (m *Manager) serveWebsocket(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil || usernameIsTaken(user, &m.clients) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing username or taken."))
		return
	}

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Could not upgrade."))
		return
	}

	wsClient := NewWsClient(conn, m, user)
	m.addClient(wsClient)

	go wsClient.readMessages()
	go wsClient.writeMessages()
}

func (m *Manager) addClient(c *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[c] = true
}

func (m *Manager) removeClient(c *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[c]; ok {
		c.connection.Close()
		delete(m.clients, c)
	}
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		return handler(event, c)
	}

	return fmt.Errorf("RouteEvent: eventType does not have matching handler: %s", event.Type)
}

// Move me out
type User struct {
	// id uuid.UUID
	name string
}

func getUser(r *http.Request) (*User, error) {
	username := r.URL.Query().Get("username")
	if username == "" {
		return nil, fmt.Errorf("Username not valid")
	}
	user := User{name: username}
	return &user, nil
}

func usernameIsTaken(u *User, clients *ClientLists) bool {
	for c := range *clients {
		if c.gameData.user.name == u.name {
			return true
		}
	}
	return false
}
