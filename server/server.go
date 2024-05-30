package mudge\server

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Infinite-X-Studios/mudge/logger"
	"github.com/google/uuid"
)

type Server struct {
	address string
	name    string
	banner  string
	logger  *logger.Logger
	users   *sync.Map
}

func New(address, name, banner string) (Server, error) {
	num := len(strings.Split(string(address), ":"))
	if num < 2 {
		err := fmt.Errorf("Invalid Address. Please use ipv4 (127.0.0.1:23) or ipv6 notation ([::1]:23).")
		return Server{}, err
	}

	logFile, err := os.OpenFile("data/server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	logger := logger.New(logFile)
	logger.Run()

	return Server{
		address: address,
		name:    name,
		banner:  banner,
		logger:  logger,
		users:   &sync.Map{},
	}, nil
}

// A helper function to type server.logger less
func (server *Server) Log(level, message string, any ...any) {
	server.logger.Log(level, message, any...)
}

// LogInfo?
// LogError?
// LogWarning?

func (server *Server) RunServer() error {
	// Make sure we stop and wait for all logging messages when leaving
	defer time.Sleep(time.Second) // Adds some delay for logging crashes
	defer server.logger.Wait()
	defer server.logger.Stop()

	listener, err := net.Listen("tcp", string(server.address))
	if err != nil {
		server.Log("ERROR", "Starting server failed: %s", err)
		return fmt.Errorf("[%s] %s %s", "ERROR", "Starting server failed:", err)
	}
	defer listener.Close()

	server.Log("INFO", "Server started: %s", server.address)

	for {
		connection, err := listener.Accept()
		if err != nil {
			server.Log("ERROR", "Accepting connection failed: %s", err)
			return fmt.Errorf("[%s] %s %s", "ERROR", "Accepting connection failed:", err)
		}

		id := uuid.New().String()
		server.users.Store(id, connection)

		conn, ok := server.users.Load(id)
		if ok {
			server.Log("INFO", "Connection Accepted: %v (%v)", id, conn)
		} else {
			server.Log("WARNING", "Getting the user failed.")
			continue
		}

		go server.handleUserConnection(conn.(net.Conn))
	}
}

func (server *Server) handleUserConnection(connection net.Conn) {
	var buffer = make([]byte, 1024)
	defer connection.Close()

	go writeConn([]byte(server.banner), connection)

	for {
		n, err := connection.Read(buffer)
		if err != nil {
			server.Log("ERROR", "Reading from connection: %v (%v)", connection, err)
			continue
		}
		// Write only the number of bytes recieved
		go writeConn(buffer[:n], connection)
	}
}

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
