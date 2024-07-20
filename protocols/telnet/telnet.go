package TELNET

const (
	NULL           = byte(0)
	LineFeed       = byte(10)
	CarriageReturn = byte(13)

	BELL          = byte(7)
	BackSpace     = byte(8)
	HorizontalTab = byte(9)
	VerticalTab   = byte(11)
	FormFeed      = byte(12)

	SE       = byte(240)
	NOP      = byte(242)
	DataMark = byte(242)

	Break            = byte(243)
	InterruptProcess = byte(244)
	AbortOutput      = byte(245)
	AreYouThere      = byte(246)
	EraseCharacter   = byte(247)
	EraseLine        = byte(248)
	GoAhead          = byte(249)
	SB               = byte(250)

	WILL = byte(251)
	WONT = byte(252)
	DO   = byte(253)
	DONT = byte(254)

	IAC = byte(255)

	GMCP = byte(201)
)

func Handshake() {

}
