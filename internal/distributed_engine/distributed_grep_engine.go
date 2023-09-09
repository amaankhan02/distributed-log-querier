package distributed_engine

import (
	"bufio"
	"cs425_mp1/internal/grep"
	"cs425_mp1/internal/network"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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
	// serverConns []net.Conn // server connections		// -> don't need this...
	// server      *network.Server
	listener   net.Listener // listener for the server
	serverQuit chan interface{}
	serverWg   sync.WaitGroup

	clientConns []net.Conn // client connections to the peers

	serverPort    string
	peerAddresses []string
	localLogFile  string
	//isRunning     bool
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
	//dpe.isRunning = false
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
		//defer conn.Close() 	// --> don't want to close the connection here itself... do it at shutdown

		dpe.clientConns = append(dpe.clientConns, conn)
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

	// dpe.isRunning = true
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
			dpe.serverWg.Add(1)
			go func() {
				dpe.handleServerConnection(conn)
				dpe.serverWg.Done()
			}()
		}
	}
}

// Handler for a connection that the server establishes with a foreign client
func (dpe *DistributedGrepEngine) handleServerConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		// TODO: make sure ReadRequest() blocks
		gQueryData, read_err := network.ReadRequest(reader)
		if read_err == io.EOF {
			log.Println("GOT EOF IN HANDLESERVERCONNECTION()... WHAT TO DO HERE!!!??") // TODO: HANDLE THIS
			return
		} else if read_err != nil {
			log.Println("GOT OTHER ERROR IN HANDLESERVERCONNECTION()... WHAT TO DO HERE!!??") // TODO: HANDE THIS
			return
		}

		gQuery, err1 := grep.DeserializeGrepQuery(gQueryData)		// ! RETURNING A NULL OBJECT
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
			_, _ = fmt.Fprintf(os.Stderr, "*FAILED* to send Grep Output Data to %s", conn.RemoteAddr().String())
			continue
		}
	}
	// TODO: how do i exit this function? Should probably have an exit feature in my program...
}

/*
Execute the grep query on local machine and all peer machines by sending grep query
to all peer machines and receive back output from them

Prints the output from each machine to stdout in a nice formatted manner
Additionally prints the total number of lines at the end
*/
func (dpe *DistributedGrepEngine) Execute(gquery *grep.GrepQuery) {
	// TODO: change all NUM_MACHINES to be just the active connected machine not NUM_MACHINES since we don't know how many are connected
	numPeerConnections := len(dpe.clientConns)
	peerChannels := make([]chan *grep.GrepOutput, numPeerConnections)
	localChannel := make(chan *grep.GrepOutput)
	var wg sync.WaitGroup
	var totalNumLines int

	for i := 0; i < numPeerConnections; i++ {
		peerChannels[i] = make(chan *grep.GrepOutput)
	}

	// launch remote executions
	for i := 0; i < numPeerConnections; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			dpe.remoteExecute(gquery, dpe.clientConns[idx], peerChannels[idx])
		}(i) // pass in i so that it accesses the correct values from peerConnections & channels since i will change
	}

	// launch goroutine for local execution
	wg.Add(1)
	go func() {
		defer wg.Done()
		dpe.localExecute(gquery, localChannel)
	}()

	// wait for all goroutines to finish  (similar to pthread_join)
	wg.Wait()

	// Print local grep output to stdout
	grepOut := <-localChannel
	totalNumLines += grepOut.NumLines
	fmt.Println(grepOut.ToString())

	// Print peer grep outputs to stdout
	for i := 0; i < len(peerChannels); i++ {
		grepOut := <-peerChannels[i] // read from channel into grep outputs array
		totalNumLines += grepOut.NumLines
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
func (dpe *DistributedGrepEngine) remoteExecute(gquery *grep.GrepQuery, conn net.Conn, outputChannel chan *grep.GrepOutput) {
	// fmt.Printf("gquery.cmdArgs: %s\n", gquery.CmdArgs)

	gquery_data, ser_err := grep.SerializeGrepQuery(gquery)
	if ser_err != nil {
		log.Fatalf("Failed to serialized gquery data")
	}
	// fmt.Printf("Serialized gquery: %v\n", gquery_data)
	err := network.SendRequest(gquery_data, conn)
	if err != nil {
		fmt.Printf("Failed to send gquery_data to %s", conn.RemoteAddr()) // TODO: how to handle this error?!
		return
	}

	// wait to recv data back
	reader := bufio.NewReader(conn)
	byte_data, err2 := network.ReadRequest(reader)
	if err2 != nil {
		fmt.Printf("Failed to read gquery_data from %s", conn.RemoteAddr()) // TODO: how to handle this error?!
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
	outputChannel <- grepOutput // TODO: is it fine to get the memory address of this var? is it stored on heap??
}

func (dpe *DistributedGrepEngine) Shutdown() {

	/*
		TODO:
			* close all peer client connections
			* close server connections (or is that already handled)
	*/
	// stop client connection first before closing server connection
	dpe.StopServer()
}

func (dpe *DistributedGrepEngine) StopServer() {
	close(dpe.serverQuit)
	err := dpe.listener.Close()
	if err != nil {
		log.Fatal("Failed to close server's listener object")
	}
	dpe.serverWg.Wait()
}
