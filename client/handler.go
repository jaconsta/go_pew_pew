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
	if c.gameData.room.selected == "" || len(c.gameData.room.players) < 2 {
		c.gameData.selectARoom()
		// Now request the users.
		outSelectRoom(c)
	}
	log.Println("Rooms updated. Total rooms available: ", len(c.gameData.room.list))
	return nil
}
func RoomUsersHander(event Event, c *Client) error {
	var data UsersInRoom
	if err := json.Unmarshal(event.Payload, &data); err != nil {
		return err
	}
	shouldChangeTarget := c.gameData.room.target == ""
	c.gameData.updatePlayers(data.Users)
	if shouldChangeTarget {
		if selected, err := c.gameData.selectATarget(); err != nil {
			log.Println("RoomUsersHandler, error selecting target")
		} else {
			log.Println("New target: ", selected)
		}
	}

	log.Println("Users updated. Total players in room (incl.): ", len(c.gameData.room.players))
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

// ----

func outSelectRoom(c *Client) {
	roomChange := ChangeRoom{
		Room: c.gameData.room.selected,
	}
	payload, err := json.Marshal(roomChange)
	if err != nil {
		log.Println("outSelectRoom error ")
		log.Println(err)
	}
	event := Event{
		Type:    ChangeRoomEvent,
		Payload: payload,
	}
	c.sendEvent(event)
}

func outShoot(c *Client) {
	impactEvent := ImpactSent{
		Target: c.gameData.room.target,
	}
	payload, err := json.Marshal(impactEvent)
	if err != nil {
		log.Println("outTarget error ")
		log.Println(err)
	}
	event := Event{
		Type:    ImpactSentEvent,
		Payload: payload,
	}
	c.sendEvent(event)
}
