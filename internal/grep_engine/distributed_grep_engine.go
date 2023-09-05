package engine

import (
	"fmt"
	"log"
	"net"
)

// DistributedEngine: Struct defining the engine to handle the Distributed Grep execution across
// multiple peer machines.
// Contains a server and client where the server accepts connections from all peer machines, and
// the Client is connected to all peer machines. During the execution of a grep query,
// the client is responsible for sending the grep query to all peer machines and receiving back
// the output. The Server is responsible for listening to all machines for any grep queries.
// When a query is received from a peer, it executes the local query and sends back the output
type DistributedGrepEngine struct {
	serverConns []net.Conn // server connections
	clientConns []net.Conn // client connections to the peers

	serverPort string
}

// Initialize all clients by connecting to all the remote servers (peers) and return
// a slice of the connection objects to all the servers/peers
// TODO: should this assume that the remote connection servers are already up and running?
func (dpe DistributedGrepEngine) initializeClients(peerAddresses []string) {
	dpe.clientConns = make([]net.Conn, 0) // connection objects of all the connected servers (peers)

	// connect to each server's ipAddress (acting as client - connecting to the servers)
	for _, peerServerAddr := range peerAddresses {
		conn, err := net.Dial("tcp", peerServerAddr) // conn = client connection object
		if err != nil {
			fmt.Printf("Error connecting to %s: %v\n", peerServerAddr, err)
			continue
		}
		defer conn.Close() // TODO: wrap error handling in a closure

		// probably need to call a goroutine to handle this client
		// for now, keep a record of all the client connections in a slice
		dpe.clientConns = append(dpe.clientConns, conn)
	}
}

// Create socket endpoint on the port passed in, and listen to any connections
// For every new accepted TCP connection, call a new goroutine to handle the
// connection --> which should end up doing an infinite loop as long as the connection
// is up and it should wait for any messages being received on that connection, process,
// send back the output, and then repeat
// TODO: ^ accepting new connections
func (dpe DistributedGrepEngine) initializeServer(port string) {
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
		fmt.Println("Server established connection to: ", conn.RemoteAddr())
		go handleServerConnection(conn)
	}
}
