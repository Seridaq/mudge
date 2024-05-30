package mudge\objects

import "github.com/google/uuid"

type Exit struct {
	Id          uuid.UUID
	description []byte
	command     []Command
}
