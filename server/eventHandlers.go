package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func ChangeRoomHandler(event Event, c *Client) error {
	var changeRoomEvent ChangeRoom
	if err := json.Unmarshal(event.Payload, &changeRoomEvent); err != nil {
		return fmt.Errorf("ChangeRoomHandler: bad payload:\n%v\n", err)
	}

	log.Println(fmt.Sprintf("User %s, changed room from %s, to %s", c.gameData.user.name, c.gameData.room, changeRoomEvent.Room))
	c.gameData.room = changeRoomEvent.Room

	RoomUsersSender(c.manager.clients, c.gameData.room)
	return nil
}

func ImpactSentHandler(event Event, c *Client) error {
	var impactSentEvent ImpactSent
	if err := json.Unmarshal(event.Payload, &impactSentEvent); err != nil {
		return fmt.Errorf("ImpactSendHandler: bad payload:\n%v\n", err)
	}

	if c.gameData.life <= 0 {
		ImpactFailedSender(c, "No lifes left")
		return fmt.Errorf("ImpactSendHandler: A zombie is shooting %s.", c.gameData.user.name)
	}

	targetFound := false
	var impactSend ImpactSend
	clientsInRoom := make([]*Client, 0)
	for targetClient := range c.manager.clients {
		if targetClient.gameData.user.name == impactSentEvent.Target {
			targetFound = true
			if targetClient.gameData.room != c.gameData.room {
				log.Println("ImpactSentHandler. Target", impactSentEvent.Target, "is in another room.", targetClient.gameData.room, ". Shooter ", c.gameData.room)
				break
			}
			if targetClient.gameData.life <= 0 {
				log.Println("ImpactSentHandler. Target is zombie.", impactSentEvent.Target, targetClient.gameData.life)
				break
			}

			targetClient.gameData.life -= 1

			impactSend = ImpactSend{
				Attacker: c.gameData.user.name,
				Target:   targetClient.gameData.user.name,
				NewLife:  targetClient.gameData.life,
			}
			ImpactSendSender(impactSend, targetClient)
		} else if targetClient.gameData.room == c.gameData.room {
			clientsInRoom = append(clientsInRoom, targetClient)
		}
	}

	if !targetFound {
		return fmt.Errorf("ImpactSentHandler. Target not found. %s\n.", impactSentEvent.Target)
	}

	if err := NotifiyImpactSender(impactSentEvent.Target, &impactSend, clientsInRoom); err != nil {
		return err
	}

	return nil
}

func ImpactSendSender(impactSend ImpactSend, c *Client) error {
	payload, err := json.Marshal(impactSend)
	if err != nil {
		return fmt.Errorf("Could not marshal impactSend. \n%v\n", err)
	}
	outgoingEvent := Event{
		Type:    ImpactSendEvent,
		Payload: payload,
	}

	c.egress <- outgoingEvent
	return nil
}

func NotifiyImpactSender(target string, impactSend *ImpactSend, clientsInRoom []*Client) error {
	payload, err := json.Marshal(impactSend)
	if err != nil {
		return fmt.Errorf("Could not marshal impactSend. \n%v\n", err)
	}
	outgoingEvent := Event{
		Type:    ImpactNotifyEvent,
		Payload: payload,
	}
	for c := range clientsInRoom {
		clientsInRoom[c].egress <- outgoingEvent
	}

	return nil
}

func ImpactFailedSender(c *Client, reason string) error {
	errMessage := ErrorMessage{
		Error: reason,
	}
	payload, err := json.Marshal(errMessage)
	if err != nil {
		return fmt.Errorf("Could not marshal impactSend. \n%v\n", err)
	}
	outgoingEvent := Event{
		Type:    ImpactFailedEvent,
		Payload: payload,
	}
	c.egress <- outgoingEvent

	return nil
}

func RoomUsersSender(c ClientLists, roomName string) {
	clientsInRoom := make([]*Client, 0)
	usernames := make([]string, 0)
	for client := range c {
		if client.gameData.room == roomName {
			usernames = append(usernames, client.gameData.user.name)
			clientsInRoom = append(clientsInRoom, client)
		}
	}

	usersInRoom := UsersInRoom{
		Room:  roomName,
		Users: usernames,
	}
	payload, err := json.Marshal(usersInRoom)
	if err != nil {
		log.Printf("Could not marshal impactSend. \n%v\n", err)
		return
	}
	outgoingEvent := Event{
		Type:    RoomUsersEvent,
		Payload: payload,
	}
	for c := range clientsInRoom {
		clientsInRoom[c].egress <- outgoingEvent
	}
}

func RoomListSender(c *Client) {
	availableRooms := AvailableRooms{
		Rooms: openRooms,
	}
	payload, err := json.Marshal(availableRooms)
	if err != nil {
		log.Printf("Could not marshal impactSend. \n%v\n", err)
		return
	}
	outgoingEvent := Event{
		Type:    RoomListEvent,
		Payload: payload,
	}
	c.egress <- outgoingEvent
}
