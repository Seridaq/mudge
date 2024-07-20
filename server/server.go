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
	"github.com/google/uuid"
)

type Server struct {
	port   string
	name   string
	banner string
	logger *logger.Logger
	users  *sync.Map
}

func New(port, name, banner string) (Server, error) {
	// We will check if the folder exist for our log file and then open it as append only.
	// TODO: We should create a new server file based off of the time. (This is super trivial and will take less time than this comment!)
	_, err := os.Stat("data")

	if os.IsNotExist(err) {
		// Create the directory and any necessary parent directories
		errDir := os.Mkdir("data", 0775)
		if errDir != nil {
			log.Fatal(errDir)
		}
	}

	logFile, err := os.OpenFile(fmt.Sprintf("data/server_%d.log", time.Now().Unix()), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	logger := logger.New(logFile, logger.LogLevelWarn)
	logger.Run()

	return Server{
		port:   ":" + port,
		name:   name,
		banner: banner,
		logger: logger,
		users:  &sync.Map{},
	}, nil
}

func (server *Server) RunServer() error {
	// Make sure we stop and wait for all logging messages when leaving
	//defer time.Sleep(time.Second) // Adds some delay for logging crashes
	defer server.logger.Wait() // Wait stops the logger and waits for it to complete.
	//defer server.logger.Stop()

	listener, err := net.Listen("tcp", server.port)

	if err != nil {
		server.LogError("Starting server failed: %s", err)
		return FormatError("server failed", err)
	}

	//listener4, err := net.ListenTCP("tcp4", &address)

	//if err != nil {
	//	server.LogError("Starting server failed: %s", err)
	//	return FormatError("server failed", err)
	//}
	//defer listener4.Close()
	//defer listener6.Close()

	server.LogInfo("Server started: %v", listener.Addr().String())

	for {
		connection, err := listener.Accept()
		//server.LogInfo(connection.LocalAddr().String())
		//server.LogInfo(err.Error())

		//connection, err = listener6.Accept()
		server.LogInfo(connection.LocalAddr().String())
		//server.LogInfo(err.Error())

		if err != nil {
			server.LogError("Accepting connection failed: %s", err)
			return FormatError("connection failed", err)
		}

		id := uuid.New().String()
		server.users.Store(id, connection)

		conn, ok := server.users.Load(id)
		if ok {
			server.LogInfo("Connection Accepted: %v (%v)", id, conn)
		} else {
			server.LogWarning("Getting the user failed.")
			continue
		}

		go server.handleUserConnection(conn.(net.Conn))
	}
}

// handleUserConnection is the client handler. Here we wait for input
// sent by the client and interact with the data to udate the state.
// The server updates the data the player receives also.
func (server *Server) handleUserConnection(connection net.Conn) {
	var buffer = make([]byte, 1024)
	defer connection.Close()

	go writeConn([]byte(server.banner), connection)

	// Let us just start with doing something simple. We will work on the
	// other adanced things later.
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
		// Write only the number of bytes recieved
		go writeConn(buffer[:n], connection)
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
