package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/google/uuid"
)

type User struct {
	username string
}

func NewUser() User {
	return User{
		username: uuid.New().String(),
	}
}

type Room struct {
	list     []string
	selected string
	players  []string
	target   string
}

func NewRoom() Room {
	return Room{
		list:     []string{},
		selected: "",
		players:  []string{},
		target:   "",
	}
}

type Game struct {
	room Room
	user *User
	life int

	sync.RWMutex
}

func NewGameData(user *User) *Game {
	return &Game{
		room: NewRoom(),
		user: user,
		life: 5,
	}
}

func (g *Game) updateRoomList(rooms []string) {
	g.Lock()
	defer g.Unlock()

	g.room.list = rooms
}

func (g *Game) updatePlayers(players []string) {
	g.Lock()
	defer g.Unlock()

	g.room.players = players
}

func (g *Game) setLife(newLife int) {
	g.Lock()
	defer g.Unlock()

	g.life = newLife
	if g.life == 0 {
		log.Fatal("I am dead... exit")
	}
}

func (g *Game) selectARoom() (string, error) {
	g.Lock()
	defer g.Unlock()

	if len(g.room.list) == 0 {
		g.room.selected = ""
		return g.room.selected, fmt.Errorf("Room list is empty")
	}
	g.room.selected = g.room.list[len(g.room.list)-1]
	return g.room.selected, nil
}

func (g *Game) selectATarget() (string, error) {
	g.Lock()
	defer g.Unlock()

	if len(g.room.players) <= 1 {
		g.room.target = ""
		return "", fmt.Errorf("Room %s is empty, current players: me, myself and I", g.room.selected)
	}

	for {
		i := rand.Intn(len(g.room.players))
		selected := g.room.players[i]
		if selected != g.user.username {
			return selected, nil
		}
	}
}
