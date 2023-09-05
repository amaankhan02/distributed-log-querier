package main

import (
	"cs425_mp1/internal/grep"
	"cs425_mp1/internal/network"
	"cs425_mp1/internal/utils"
	"fmt"
	"log"
	"net"
	"sync"
)

/**
TODO: QUESTIONS!!
	* is there a possiblity that the file size may not fit in memory... cuz then i can't keep
	  a GrepOutput object in memory... since it holds the entire output in memory...
	* in my socket message requests, i put the [size] as a 4-Byte number, which can represent up to ~4GB file size.
		is this enough? Or should I use 8-Byte number to represent the size instead?
	* is conn.Write() a bufferred write or do i need to create a loop like in c to send all of it
*/

const (
	MACHINE_NAME_FORMAT = "fa23-cs425-19%02d.cs.illinois.edu"
	NUM_MACHINES        = 10       // num of total machines in the network, although you should be able to use less
	PORT_FORMAT         = "80%02d" // 8001, 8002, ... 8010 - based on the
)

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
		defer conn.Close() // TODO: wrap error handling in a closure

		// probably need to call a goroutine to handle this client
		// for now, keep a record of all the client connections in a slice
		peerConns = append(peerConns, conn)
	}

	return peerConns
}

/*
Execute the grep query on local machine and all peer machines by sending grep query
to all peer machines and receive back output from them

Prints the output from each machine to stdout in a nice formatted manner
Additionally prints the total number of lines at the end
*/
func distributedExecute(gquery grep.GrepQuery, peerConnections []net.Conn) {
	// TODO: change all NUM_MACHINES to be just the active connected machine not NUM_MACHINES since we don't know how many are connected

	peerChannels := make([]chan *grep.GrepOutput, NUM_MACHINES-1)
	localChannel := make(chan *grep.GrepOutput)
	var wg sync.WaitGroup
	var totalNumLines int

	for i := 0; i < NUM_MACHINES-1; i++ {
		peerChannels[i] = make(chan *grep.GrepOutput)
	}

	// launch remote executions
	for i := 0; i < NUM_MACHINES-1; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			remoteExecute(gquery, peerConnections[idx], peerChannels[idx])
		}(i) // pass in i so that it accesses the correct values from peerConnections & channels since i will change
	}

	// launch goroutine for local execution
	wg.Add(1)
	go func() {
		defer wg.Done()
		localExecute(gquery, localChannel)
	}()

	// wait for all goroutines to finish  (similar to pthread_join)
	wg.Wait()

	// Print local grep output to stdout
	grepOut := <-localChannel
	totalNumLines += grepOut.NumLines()
	fmt.Println(grepOut.ToString())

	// Print peer grep outputs to stdout
	for i := 0; i < len(peerChannels); i++ {
		grepOut := <-peerChannels[i] // read from channel into grep outputs array
		totalNumLines += grepOut.NumLines()
		fmt.Println(grepOut.ToString())
	}

	fmt.Printf("Total Number of Lines: %d\n", totalNumLines)
}

/*
Execute a grep query on a remote machine by sending the query to the machine
and waiting to receive the output and then returning it.

Designed to be ran as a goroutine.

Parameters:

	gquery: query to execute
	conn: net.Conn client object to the remote machine
	outputChannel: channel that remoteExecute() will send its grep output to
*/
func remoteExecute(gquery grep.GrepQuery, conn net.Conn, outputChannel chan *grep.GrepOutput) {
	gquery_data := grep.SerializeGrepQuery(gquery)
	err := network.SendRequest(gquery_data, conn)
	if err != nil {
		fmt.Printf("Failed to send gquery_data to %s", conn.RemoteAddr()) // TODO: how to handle this error?!
		return
	}

	// wait to recv data back
	byte_data, err2 := network.ReadRequest(conn)
	if err2 != nil {
		fmt.Printf("Failed to read gquery_data from %s", conn.RemoteAddr()) // TODO: how to handle this error?!
		return
	}

	var grepOutput = grep.DeserializeGrepOutput(byte_data)
	outputChannel <- grepOutput
}

func localExecute(gquery grep.GrepQuery, outputChannel chan *grep.GrepOutput) {
	grepOutput := gquery.Execute(utils.GetLocalLogFile())
	outputChannel <- &grepOutput // TODO: is it fine to get the memory address of this var? is it stored on heap??
}

func main() {
	peerServerAddresses := utils.GetPeerServerAddresses(MACHINE_NAME_FORMAT, PORT_FORMAT, NUM_MACHINES)
	serverPort := utils.GetLocalhostPort(MACHINE_NAME_FORMAT, PORT_FORMAT, NUM_MACHINES)

	go initializeServer(serverPort) // create server and listen to connections & accept them

	// TODO: maybe initialize clients once all the computer's servers are booted up
	peerConns := initializeClients(peerServerAddresses)

	// Enter infinite loop and display prompts to user, while getting the query
	for {
		utils.DisplayGrepPrompt()
		userInput := utils.ReadUserInput()

		gq, err := grep.CreateGrepQueryFromInput(userInput)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		distributedExecute(gq, peerConns)
	}

}
