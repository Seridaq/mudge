package objects

//import "github.com/google/uuid"

type Room struct {
	//id          uuid.UUID
	title       string
	description string
	//shortDesc   string
	//exits       []Exit
	//players     []*Player
}

func NewRoom(title string) *Room {
	return &Room{
		title: title,
		//shortDesc:   "",
		description: "This is a blank room. The empty void surrounds you.",
	}
}

func (room *Room) GetDescription() string {
	return room.title + "\r\n" + room.description
}

// Tunnel connects one room to another via an exit
//func (room *Room) Tunnel(connectTo Room) {}

// Broadcast message to room
//func (room *Room) BroadcastMessage(msg string) {}

// Give a list of players in the room
//func (room *Room) ListPlayers() []*Player {
//return room.players
//}
