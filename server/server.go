package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/Infinite-X-Studios/mudge/logger"
	"github.com/Infinite-X-Studios/mudge/objects"
	"github.com/google/uuid"
)

type Server struct {
	port        string
	name        string
	banner      string
	loggerLevel string
	logger      *logger.Logger
	users       *sync.Map
	startingMap *objects.Room
}

func New(port, name, banner string) (Server, error) {
	// We will check if the folder exist for our log file and then open it as append only.
	_, err := os.Stat("data")

	if os.IsNotExist(err) {
		// Create the directory and any necessary parent directories
		errDir := os.Mkdir("data", 0775)
		if errDir != nil {
			log.Fatal(errDir)
		}
	}

	// The log file is now created based off the time the server was started.
	logFile, err := os.OpenFile(fmt.Sprintf("data/server_%d.log", time.Now().Unix()), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	// TODO: I want to be able to set the loggers log level in a config file.
	// We will worry about this later, as we do not have a config file yet.
	logger := logger.New(logFile, logger.LogLevelInfo)
	logger.Run()

	// TODO: Remove this when we can!
	// While we are debugging, let us create a room for testing
	theVoid := objects.NewRoom("The Void")

	return Server{
		port:        ":" + port,
		name:        name,
		banner:      banner,
		logger:      logger,
		users:       &sync.Map{},
		startingMap: theVoid,
	}, nil
}

func (server *Server) RunServer() error {
	// Wait stops the logger and waits for it to complete.
	defer server.logger.Wait()

	// NOTE: I was having trouble with TCP listeners created by net.Listen("tcp", server.port) only
	// binding to IPv6 and not IPv4. A workaround to this is to create both a TCPv4 and TCPv6
	// listener and listen for connections separately.

	// Try and create a TCP v4 listener
	listener4, err4 := net.Listen("tcp4", server.port)
	if err4 != nil {
		server.LogError("Failed to listen on IPv4: %v", err4)
	} else {
		defer listener4.Close()
		server.LogInfo("Server started: %v", listener4.Addr().String())
	}

	// Try to create a TCP v6 listener
	listener6, err6 := net.Listen("tcp6", server.port)
	if err6 != nil {
		server.LogError("Failed to listen on IPv6: %v", err6)
	} else {
		defer listener6.Close()
		server.LogInfo("Server started: %v", listener6.Addr().String())
	}

	// Start the server if at least one of the listeners exist
	if err4 == nil && err6 == nil {
		// Run the IPv4 listener on a separate go routine
		go server.acceptConnections(listener4)
		// Listen for IPv6 connections on the main thread
		server.acceptConnections(listener6)
	} else if err4 == nil {
		// Listen for IPv4 connections on the main thread
		server.acceptConnections(listener4)
	} else if err6 == nil {
		// Listen for IPv6 connections on the main thread
		server.acceptConnections(listener6)
	} else {
		// Neither listener was able to start. The server fails
		return fmt.Errorf("[ERROR] The server failed to start on IPv4 and IPv6: \r\n %v \r\n %v", err4, err6)
	}

	return nil
}

func (server *Server) acceptConnections(listener net.Listener) {
	for {
		connection, err := listener.Accept()
		//server.LogInfo(connection.LocalAddr().String())
		//server.LogInfo(err.Error())

		//connection, err = listener6.Accept()
		server.LogInfo(connection.LocalAddr().String())
		//server.LogInfo(err.Error())

		if err != nil {
			server.LogError("Accepting connection failed: %s", err)
			//return FormatError("connection failed", err)
		}

		id := uuid.New().String()
		player := objects.NewPlayer("")
		server.users.Store(player, connection)

		conn, ok := server.users.Load(player)
		if ok {
			server.LogInfo("Connection Accepted: %v (%v)", id, conn)
		} else {
			server.LogWarning("Getting the user failed.")
			continue
		}

		// The user is connected, let's start handling that connection
		go server.handleUserConnection(conn.(net.Conn))
	}
}

// handleUserConnection is the client handler. Here we wait for input
// sent by the client and interact with the data to udate the state.
// The server updates the data the player receives also.
func (server *Server) handleUserConnection(connection net.Conn) {
	var buffer = make([]byte, 1024)
	defer connection.Close()

	writeConn([]byte(server.banner), connection)

	// Let us just start with doing something simple. We will work on the
	// other advanced things later.
	//fmt.Println("GMCP Handshake...")
	//GMCP.Handshake(connection)

	for {
		n, err := connection.Read(buffer)
		if err != nil {
			if err == io.EOF {
				server.LogInfo("Connection closed: %v", connection)
				break
			}

			server.LogError("Reading from connection: %v (%v)", connection, err)
			fmt.Println("Error: ", err)
			continue
		}

		// Write only the number of bytes received
		//go writeConn(buffer[:n], connection)
	}
}

// We write all data back to the sending client. This is an echo server.
func writeConn(data []byte, connection net.Conn) {
	var start, c int
	var err error
	for {
		c, err = connection.Write(data[start:])
		if err != nil {
			fmt.Println(err)
			return
		}
		start += c
		if c == 0 || start == len(data) {
			break
		}
	}
}
