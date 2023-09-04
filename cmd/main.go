package main

import (
	"cs425_mp1/internal"
	"cs425_mp1/internal/grep"
	"fmt"
	"log"
	"net"
)

const (
	MACHINE_NAME_FORMAT = "fa23-cs425-19%02d.cs.illinois.edu"
	NUM_MACHINES        = 10
	PORT_FORMAT         = "80%02d" // 8001, 8002, ... 8010 - based on the
)

// Constructs and returns a socket request message for sending a grep query
func createQueryRequest(query string) {

}

func createOutputRequest(filename string, outputData string) {

}

// Create socket endpoint on the port passed in, and listen to any connections
// For every new accepted TCP connection, call a new goroutine to handle the
// connection --> which should end up doing an infinite loop as long as the connection
// is up and it should wait for any messages being received on that connection, process,
// send back the output, and then repeat
// TODO: change name from initializeServer() to something else b/c its not just initializing it's also
// TODO: ^ accepting new connections
func initializeServer(port string) {
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("net.Listen(): %v", err)
	}
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			log.Fatalf("listen.Close(): %v", err)
		}
	}(listen)

	// Listen for connections, accept, and spawn goroutine to handle that connection
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalf("listen.Accept(): %v", err)
		}

		// TODO: do i need to defer conn.Close()? and do i need to use any waitGroups...?
		fmt.Println("Connected to: ", conn.RemoteAddr())
		go handleServerConnection(conn)
	}
}

// Handler for a connection that the server establishes with a foreign client
func handleServerConnection(clientConn net.Conn) {
	// keep looping through waiting for data to be received
	// read incoming data (grep query)
	// process and execute query
	// send back the output
}

// Initialize all clients by connecting to all the remote servers (peers) and return
// a slice of the connection objects to all the servers/peers
func initializeClients(peerAddresses []string) []net.Conn {
	peerConns := make([]net.Conn, 0) // connection objects of all the connected servers (peers)

	// connect to each server's ipAddress (acting as client - connecting to the servers)
	for _, peerServerAddr := range peerAddresses {
		conn, err := net.Dial("tcp", peerServerAddr) // conn = client connection object
		if err != nil {
			fmt.Printf("Error connecting to %s: %v\n", peerServerAddr, err)
			continue
		}
		defer conn.Close()

		// probably need to call a goroutine to handle this client
		// for now, keep a record of all the client connections in a slice
		peerConns = append(peerConns, conn)
	}

	return peerConns
}

func main() {
	peerServerAddresses := internal.GetPeerServerAddresses(MACHINE_NAME_FORMAT, PORT_FORMAT, NUM_MACHINES)
	serverPort := internal.GetLocalhostPort(MACHINE_NAME_FORMAT, PORT_FORMAT, NUM_MACHINES)

	go initializeServer(serverPort) // create server and listen to connections & accept them

	// TODO: maybe initialize clients once all the computer's servers are booted up
	peerConns := initializeClients(peerServerAddresses)

	// Enter infinite loop and display prompts to user, while getting the query
	for {
		internal.DisplayGrepPrompt()
		userInput := internal.ReadUserInput()

		gq, err := grep.CreateGrepQueryFromInput(userInput)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		// Also execute the grep query locally
		// and then wait to receive back the outputs from all the peers
		// parse the outputs received from all peers and display to user
	}

}
