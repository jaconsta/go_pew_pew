package main

import "encoding/json"

const (
	// IN
	ChangeRoomEvent = "CHANGE_ROOM"
	ImpactSentEvent = "SHOOT"
	// OUT
	ImpactSendEvent   = "RECEIVE_IMPACT"
	ImpactNotifyEvent = "NOTIFY_IMPACT"
)

type EventHandler func(event Event, c *Client) error

// Parent structure of the event
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// Specific Event's payload

type ChangeRoom struct {
	Room string `json:"room"`
}

type ImpactSent struct {
	Target string `json:"target"`
}

type ImpactSend struct {
	Attacker string `json:"attacker"`
	Target   string `json:"target"`
	NewLife  int    `json:"newLife"`
}

// type ImpactNotify struct {
// 	// From
// 	Attacker string `json:"attacker"`
// 	// To
// 	Target  string `json:"target"`
// 	NewLife int    `json:"newLife"`
// }
