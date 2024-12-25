package objects

import "github.com/google/uuid"

type Player struct {
	id          uuid.UUID
	name        string
	currentRoom *Room
}

func NewPlayer(name string, id uuid.UUID) Player {
	return Player{
		id:          id,
		name:        name,
		currentRoom: nil,
	}
}

func (player *Player) GetName() string {
	return player.name
}

func (player *Player) GetID() string {
	return player.id
}
