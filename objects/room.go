package mudge\objects

import "github.com/google/uuid"

type Room struct {
	Id               uuid.UUID
	ShortDescription []byte
	Description      []byte
	Exits            []Exit
	LinkedRoom       uuid.UUID
}

func New() {}

// Tunnel connects one room to another via an exit
func (room *Room) Tunnel(connectTo Room) {}
