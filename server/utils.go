package server

import (
	"fmt"
	"net"
)

// A client is an active connection
type Client struct {
	net.Conn
}

// Handshake will take an unsecured connection, perform user login and authentication
// and will upgrade it to a secure connection. This will return a user account if successful
// or an error if something goes wrong so it can be handeled by the server.
func (client Client) handshake() error { //(account.Account, error) {
	// Need to check out GMCP Authentication on mudlets page
	// https://wiki.mudlet.org/w/Special:MyLanguage/Standards:GMCP_Authentication
	return fmt.Errorf("client handshake Not Yet Implemented")
}

// enableMCP will attempt to establish an agreement with
// a client to use MUD Client Protocol.
func (client Client) enableMCP() error {
	return fmt.Errorf("MCP Not Yet Implemented")
}

// enableMCCP will attempt to establish an agreement with
// a client to use MUD Client Compression Protocol
func (client Client) enableMCCP() error {
	return fmt.Errorf("MCCP Not Yet Implemented")
}

// enableGMCP will attempt to establish an agreement with
// a client to use Generic MUD Communication Protocol.
func (client Client) enableGMCP() error {
	return fmt.Errorf("GMCP Not Yet Implemented")
}

// enableMCMP will attempt to establish an agreement with
// a client to use MUD Client Media Protocol.
// https://wiki.mudlet.org/w/Standards:MUD_Client_Media_Protocol
func (client Client) enableMCMP() error {
	return fmt.Errorf("MCMP Not Yet Implemented")
}

func (client Client) enableMMP() error {
	return fmt.Errorf("MMP Not Yet Implemented")
}

func (client Client) enableDiscordGMCP() error {
	return fmt.Errorf("Discord GMCP Not Yet Implemented")
}
