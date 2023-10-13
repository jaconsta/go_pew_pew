package main

import (
	"encoding/json"
	"log"
)

func newEventHandlers() map[string]EventHandler {
	handlers := make(map[string]EventHandler)

	handlers[ImpactSendEvent] = ImpactSendHandler
	handlers[ImpactNotifyEvent] = doNothingHandler
	handlers[ImpactFailedEvent] = ImpactFailedHandler
	handlers[RoomListEvent] = RoomListHander
	handlers[RoomUsersEvent] = RoomUsersHander

	return handlers
}

func undefinedHandler(event Event, c *Client) error {
	log.Println("Undefined handler for: ", event.Type)
	return nil
}
func doNothingHandler(event Event, c *Client) error {
	log.Println("Undefined handler for: ", event.Type)
	return nil
}

func RoomListHander(event Event, c *Client) error {
	var data AvailableRooms
	if err := json.Unmarshal(event.Payload, &data); err != nil {
		return err
	}
	c.gameData.updateRoomList(data.Rooms)
	return nil
}
func RoomUsersHander(event Event, c *Client) error {
	var data UsersInRoom
	if err := json.Unmarshal(event.Payload, &data); err != nil {
		return err
	}
	c.gameData.updatePlayers(data.Users)
	return nil
}

func ImpactSendHandler(event Event, c *Client) error {
	var data ImpactSend
	if err := json.Unmarshal(event.Payload, &data); err != nil {
		return err
	}
	log.Println("Attack received from ", data.Attacker)
	c.gameData.setLife(data.NewLife)

	return nil
}

func ImpactFailedHandler(event Event, c *Client) error {
	for {
		log.Println("Changing target")
		if _, err := c.gameData.selectATarget(); err != nil {
			log.Println("Changing room")
			if _, err := c.gameData.selectARoom(); err == nil {
				break
			}
		} else {
			break
		}
	}
	return nil
}
