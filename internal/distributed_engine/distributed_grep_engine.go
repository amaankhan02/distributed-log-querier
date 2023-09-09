package distributed_engine

import (
	"bufio"
	"cs425_mp1/internal/grep"
	"cs425_mp1/internal/network"
	"cs425_mp1/internal/utils"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

// TODO: CHANGE ALL PASS BY VALUE STRUCT FUNCTION TO PASS BY POINTER (s *Server) instead

// DistributedEngine: Struct defining the distributed_engine to handle the Distributed Grep execution across
// multiple peer machines.
// Contains a server and client where the server accepts connections from all peer machines, and
// the Client is connected to all peer machines. During the execution of a grep query,
// the client is responsible for sending the grep query to all peer machines and receiving back
// the output. The Server is responsible for listening to all machines for any grep queries.
// When a query is received from a peer, it executes the local query and sends back the output
type DistributedGrepEngine struct {
	listener   net.Listener // listener for the server
	serverQuit chan interface{}
	serverWg   sync.WaitGroup

	clientConns      []net.Conn        // client connections to the peers
	activeClients    map[net.Conn]bool // hash-map where value = True if the connection is active. False o.w.
	numActiveClients int

	serverPort    string
	peerAddresses []string
	localLogFile  string
}

/*
Creates a DistributedGrepEngine struct and initializes with default values
*/
func CreateEngine(localLogFile string, serverPort string, peerAddresses []string) *DistributedGrepEngine {
	// initialize server and client connections here
	dpe := &DistributedGrepEngine{}
	dpe.localLogFile = localLogFile
	dpe.serverPort = serverPort
	dpe.peerAddresses = peerAddresses
	dpe.activeClients = make(map[net.Conn]bool)
	return dpe
}

// Initialize all clients by connecting to all the remote servers (peers)
// This function assumes that the Peers are already setup with their server running. That is,
// It will only connect to the machines that have their servers setup
func (dpe *DistributedGrepEngine) ConnectToPeers() {
	dpe.clientConns = make([]net.Conn, 0) // connection objects of all the connected servers (peers)

	// connect to each server's ipAddress (acting as client - connecting to the servers)
	for _, peerServerAddr := range dpe.peerAddresses {
		conn, err := net.Dial("tcp", peerServerAddr) // conn = client connection object
		if err != nil {
			fmt.Printf("Error connecting to %s: %v\n", peerServerAddr, err)
			continue
		}
		dpe.clientConns = append(dpe.clientConns, conn)
		dpe.activeClients[conn] = true
		dpe.numActiveClients += 1
	}
}

// Initialize Server on a separate goroutine and engine now actively listens to new connections
// TODO: change name to StartServer() later
func (dpe *DistributedGrepEngine) InitializeServer() {
	l, err := net.Listen("tcp", dpe.serverPort)
	if err != nil {
		log.Fatalf("net.Listen(): %v", err)
	}
	dpe.listener = l
	dpe.serverWg.Add(1)
	go dpe.serve()
}

// Helper function
// Create socket endpoint on the port passed in, and listen to new connections to this server
// For every new accepted TCP connection, call a new goroutine to handle that connection
func (dpe *DistributedGrepEngine) serve() {
	defer dpe.serverWg.Done()

	// Listen for connections, accept, and spawn goroutine to handle that connection
	for {
		conn, err := dpe.listener.Accept()
		if err != nil {
			select {
			case <-dpe.serverQuit: // signal server to quit
				return
			default:
				log.Println("Accept() error: ", err)
			}
		} else {
			fmt.Println("Server connected to: ", conn.RemoteAddr())
			connectionMsg := fmt.Sprintf("Server connected to: %s", conn.RemoteAddr())
			utils.PrintMessage(connectionMsg)
			dpe.serverWg.Add(1)
			go func() {
				defer dpe.serverWg.Done()
				dpe.handleServerConnection(conn)
			}()
		}
	}
}

// Handler for a connection that the server establishes with a foreign client
func (dpe *DistributedGrepEngine) handleServerConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		gQueryData, read_err := network.ReadRequest(reader)
		if read_err == io.EOF { // this server-client connection disconnected, so we can remove
			msg := fmt.Sprintf("\n**Client [%s] disconnected**\n", conn.RemoteAddr().String())
			dpe.removeClient(conn)
			utils.PrintMessage(msg)
			return
		} else if read_err != nil {
			log.Fatalf("Error while performing network.ReadRequest()!")
			return
		}

		gQuery, err1 := grep.DeserializeGrepQuery(gQueryData)
		if err1 != nil {
			log.Fatalf("Failed to Deserialize Grep Query: %v", err1)
		}

		gOut := gQuery.Execute(dpe.localLogFile)
		gOutData, err2 := grep.SerializeGrepOutput(gOut)
		if err2 != nil {
			log.Fatalf("Failed to Serialize Grep Output: %v", err2)
		}
		err := network.SendRequest(gOutData, conn)
		if err != nil {
			log.Fatalf("*FAILED* to send Grep Output Data to %s", conn.RemoteAddr().String())
			return // TODO: should I "return" or "continue"
		}
	}
}

/*
Execute the grep query on local machine and all peer machines by sending grep query
to all peer machines and receive back output from them

Prints the output from each machine to stdout in a nice formatted manner
Additionally prints the total number of lines at the end
*/
func (dpe *DistributedGrepEngine) Execute(gquery *grep.GrepQuery) {
	numTotalPeerConnections := len(dpe.clientConns)
	localChannel := make(chan *grep.GrepOutput)
	var totalNumLines int

	peerChannels := make([]chan *grep.GrepOutput, dpe.numActiveClients)
	for i := 0; i < dpe.numActiveClients; i++ {
		peerChannels[i] = make(chan *grep.GrepOutput)
	}

	// launch goroutines for local and remote executions to all run in parallel
	go dpe.localExecute(gquery, localChannel)
	var peerChannelIdx = 0
	for i := 0; i < numTotalPeerConnections; i++ {
		if dpe.activeClients[dpe.clientConns[i]] == true {
			go dpe.remoteExecute(gquery, dpe.clientConns[i], peerChannels[peerChannelIdx])
			peerChannelIdx += 1
		}
	}

	// * NOTE: localExecute() and remoteExecute() will not exit until its respective channels are read from since the channels
	// * once written to will block until someone reads from them. Therefore it will block until it is read from below

	// Print local grep output to stdout
	grepOut := <-localChannel
	totalNumLines += grepOut.NumLines
	fmt.Print(grepOut.ToString())

	// Print peer grep outputs to stdout
	for i := 0; i < len(peerChannels); i++ {
		grepOut := <-peerChannels[i] // read from channel into grep outputs array
		totalNumLines += grepOut.NumLines
		fmt.Print(grepOut.ToString())
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
func (dpe *DistributedGrepEngine) remoteExecute(gquery *grep.GrepQuery, conn net.Conn, outputChannel chan *grep.GrepOutput) {
	gquery_data, ser_err := grep.SerializeGrepQuery(gquery)
	if ser_err != nil {
		log.Fatalf("Failed to serialized gquery data")
	}

	err := network.SendRequest(gquery_data, conn)
	if err != nil {
		fmt.Printf("Failed to send gquery_data to %s\n", conn.RemoteAddr()) // TODO: how to handle this error?!
		return
	}

	// wait to recv data back
	reader := bufio.NewReader(conn)
	byte_data, err2 := network.ReadRequest(reader)
	if err2 != nil {
		fmt.Printf("Failed to read gquery_data from %s\n", conn.RemoteAddr()) // TODO: how to handle this error?!
		return
	}

	grepOutput, err1 := grep.DeserializeGrepOutput(byte_data)
	if err1 != nil {
		log.Fatalf("Failed to Deserialize Grep Output: %v", err1)
	}

	outputChannel <- grepOutput
}

func (dpe *DistributedGrepEngine) localExecute(gquery *grep.GrepQuery, outputChannel chan *grep.GrepOutput) {
	grepOutput := gquery.Execute(dpe.localLogFile)
	outputChannel <- grepOutput
}

func (dpe *DistributedGrepEngine) Shutdown() {
	dpe.StopServer()
	// TODO: anything else needed here?
}

func (dpe *DistributedGrepEngine) StopServer() {
	close(dpe.serverQuit)
	err := dpe.listener.Close()
	if err != nil {
		log.Fatal("Failed to close server's listener object")
	}
	dpe.serverWg.Wait()
}

// When a client was disconnected, call this function to remove
// the client information from the DistributedGrepEngine struct
func (dpe *DistributedGrepEngine) removeClient(conn net.Conn) {
	dpe.activeClients[conn] = false
	dpe.numActiveClients -= 1
}
