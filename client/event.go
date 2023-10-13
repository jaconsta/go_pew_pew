package main

import "encoding/json"

const (
	// IN
	ChangeRoomEvent = "CHANGE_ROOM"
	ImpactSentEvent = "SHOOT"
	// OUT
	ImpactSendEvent   = "RECEIVE_IMPACT"
	ImpactNotifyEvent = "NOTIFY_IMPACT"
	ImpactFailedEvent = "FAILED_IMPACT"
	RoomListEvent     = "ROOM_LIST"
	RoomUsersEvent    = "ROOM_USERS"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

type AvailableRooms struct {
	Rooms []string `json:"rooms"`
}

type ChangeRoom struct {
	Room string `json:"room"`
}

type UsersInRoom struct {
	Room  string   `json:"room"`
	Users []string `json:"users"`
}

type ImpactSent struct {
	Target string `json:"target"`
}

type ImpactSend struct {
	Attacker string `json:"attacker"`
	Target   string `json:"target"`
	NewLife  int    `json:"newLife"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}
