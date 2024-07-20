package GMCP

import (
	"fmt"
	"net"
	"time"
)

const (
	SE   byte = 0xF0
	SB   byte = 0xFA
	WILL byte = 0xFB
	WONT byte = 0xFC
	DO   byte = 0xFD
	DONT byte = 0xFE
	IAC  byte = 0xFF
	GMCP byte = 0xC9
)

func Handshake(conn net.Conn) error {
	fmt.Println("Inside Handshake...")
	response := make([]byte, 1024)

	_, err := conn.Write([]byte{IAC, WILL, GMCP})

	if err != nil {
		fmt.Println("Error: ", err)
	}

	var i int
L:
	for {
		fmt.Println("Read: ")
		conn.SetReadDeadline(time.Now().Add(time.Second))
		n, err := conn.Read(response)
		if err != nil {
			fmt.Println("Error: ", err, "(", n, ")")
			continue
		}

		if response[0] == IAC {
			if response[1] == SB {
				if response[2] == GMCP {
					fmt.Println("Sub-negotiation started: ")
				}
			}
		} else {
			fmt.Println("Not Sub-negotiation.")
			return fmt.Errorf("Not Sub-negotiation.")
		}

		var byt byte

		for i, byt = range response[3:n] {
			if byt == IAC {
				if response[i] == SE {
					fmt.Println("Sub-negotiation ended: ")
					fmt.Println(i == n-5)
					break L
				}
			}
		}
		if i <= 1 {
			fmt.Println("This is really bad, and should not happen.")
		}
	}

	//fmt.Println("Received: ", response[:n])
	fmt.Println(string(response[3 : i-1]))

	_, err = conn.Write([]byte{IAC, SB, GMCP})
	return nil
}
