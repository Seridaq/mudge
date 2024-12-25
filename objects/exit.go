package objects

import (
	"github.com/Infinite-X-Studios/mudge/objects/common"
	"github.com/google/uuid"
)

type Exit struct {
	id          uuid.UUID
	description string
	command     common.Command
	targetRoom  *Room
}

func NewExit() Exit {
	return Exit{}
}

func (exit *Exit) GetTargetRoom() *Room {
	return exit.targetRoom
}
