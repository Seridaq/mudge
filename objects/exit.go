package objects

import (
	"github.com/Infinite-X-Studios/mudge/commands"
	"github.com/google/uuid"
)

type Exit struct {
	Id          uuid.UUID
	description []byte
	command     []commands.Command
}
