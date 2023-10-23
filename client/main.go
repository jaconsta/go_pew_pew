package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

func main() {
	// Connect
	client, err := newClient()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	var wg sync.WaitGroup
	wg.Add(1)
	go readLoop(client, &wg)
	quit := make(chan struct{})
	go writeLoop(client, &quit)
	wg.Wait()
	close(quit)
}

func readLoop(c *Client, wg *sync.WaitGroup) {
	for {
		var msg = make([]byte, 512)
		n, err := c.conn.Read(msg)
		if err != nil {
			log.Println(err)
			break
		}

		var eventData Event
		if err := json.Unmarshal(msg[:n], &eventData); err != nil {
			log.Println("Read message. Received bad message: ", err)
		}
		if err := c.routeEvent(eventData); err != nil {
			log.Println("Read message, error handling: ", err)
		}
	}
	wg.Done()
}

func writeLoop(c *Client, q *chan struct{}) {
	shootTicker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-shootTicker.C:
			if c.gameData.ready {
				log.Println("Aiming")
				outShoot(c)
			}

		case <-*q:
			shootTicker.Stop()
			return
		}
	}
}
