package distributed_engine

import (
	"bufio"
	"cs425_mp1/internal/grep"
	"cs425_mp1/internal/network"
	"cs425_mp1/internal/utils"
	"encoding/json"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

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

	clientConns      []net.Conn      // client connections to the peers
	activeClients    map[string]bool // key = addr of client, value = True if connection is active. False if disconnected
	numActiveClients int

	serverPort               string
	peerAddresses            []string
	localLogFile             string
	testOutputFileNameFormat string

	lruCache                *lru.Cache
	cacheInitalizationError error

	verbose            bool
	currentTestFileIdx int
}

type JSONOutput struct {
	Query   string            // packaged string of the grep query
	Outputs []grep.GrepOutput // list of grep_outputs, each grep_output in this list is from a different vm
}

// you want to essentially create a list of these to store in JSON file

/*
Creates a DistributedGrepEngine struct and initializes with default values
*/
func CreateEngine(localLogFile string, serverPort string, peerAddresses []string, cacheSize int, verbose bool, testOutputFileNameFormat string) *DistributedGrepEngine {
	// initialize server and client connections here

	// initialize cache
	dpe := &DistributedGrepEngine{}
	dpe.localLogFile = localLogFile
	dpe.serverPort = serverPort
	dpe.peerAddresses = peerAddresses
	dpe.activeClients = make(map[string]bool)
	dpe.verbose = verbose
	dpe.testOutputFileNameFormat = testOutputFileNameFormat
	dpe.currentTestFileIdx = 1

	dpe.lruCache, dpe.cacheInitalizationError = lru.New(cacheSize)
	if dpe.cacheInitalizationError != nil {
		log.Fatalf("Error in initializing LRU Cache")
	}

	return dpe
}

// Initialize all clients by connecting to all the remote servers (peers)
// This function assumes that the Peers are already setup with their server running. That is,
// It will only connect to the machines that have their servers setup
func (dpe *DistributedGrepEngine) ConnectToPeers() {
	dpe.clientConns = make([]net.Conn, 0) // connection objects of all the connected servers (peers)

	// connect to each server's ipAddress (acting as client - connecting to the servers)
	for _, peerServerAddr := range dpe.peerAddresses {
		didConnect := false
		for !didConnect {
			conn, err := net.Dial("tcp", peerServerAddr) // conn = client connection object
			if err == nil {                              // successfully connected
				//fmt.Printf("Error connecting to %s: %v\n", peerServerAddr, err)
				//continue
				dpe.clientConns = append(dpe.clientConns, conn)
				dpe.activeClients[generateClientConnKey(conn)] = true
				dpe.numActiveClients += 1
				didConnect = true
			} else {
				time.Sleep(125 * time.Millisecond) // wait 0.125 seconds before trying again to connect
			}
		}

	}
}

// Initialize Server on a separate goroutine and engine now actively listens to new connections
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
			connectionMsg := fmt.Sprintf("Server connected to: %s", conn.RemoteAddr())
			utils.PrintMessage(connectionMsg, dpe.verbose)
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
			utils.PrintMessage(msg, dpe.verbose)
			return
		} else if read_err != nil {
			log.Fatalf("Error while performing network.ReadRequest()!")
		}

		gQuery, err1 := grep.DeserializeGrepQuery(gQueryData)
		if err1 != nil {
			log.Fatalf("Failed to Deserialize Grep Query: %v", err1)
		}

		gOut := dpe.checkCacheOrExecute(gQuery) // retrieve output from cache or execute function if not in there

		gOutData, err2 := grep.SerializeGrepOutput(gOut)
		if err2 != nil {
			log.Fatalf("Failed to Serialize Grep Output: %v", err2)
		}

		err := network.SendRequest(gOutData, conn)
		if err != nil {
			log.Fatalf("SendRequest: Failed to send Grep Output Data to %s", conn.RemoteAddr().String())
		}
	}
}

// Helper function that first checks if the query is present in the cache.
// If it is, it returns the output from the cache as well as updating the LRU position of the cache
// Otherwise, it executes the grep query and stores output in the cache, and then returns the output
func (dpe *DistributedGrepEngine) checkCacheOrExecute(gQuery *grep.GrepQuery) *grep.GrepOutput {
	var cacheValue interface{}

	var gOut *grep.GrepOutput
	var ok bool

	cacheKey := gQuery.PackagedString
	if dpe.lruCache.Contains(cacheKey) {
		start := time.Now()
		cacheValue, ok = dpe.lruCache.Get(cacheKey)
		if !ok {
			log.Fatalf("Error in getting cache value: lruCache.Get(%s)", cacheKey)
		}
		gOut = cacheValue.(*grep.GrepOutput)
		end := time.Now()
		elapsed := end.Sub(start)
		gOut.ExecutionTime = elapsed // update exec time since we now got it from cache
	} else {
		gOut = gQuery.Execute(dpe.localLogFile)
		dpe.lruCache.Add(cacheKey, gOut)
	}

	return gOut
}

/*
Execute the grep query on local machine and all peer machines by sending grep query
to all peer machines and receive back output from them

Prints the output from each machine to stdout in a nice formatted manner
Additionally prints the total number of lines at the end
*/
func (dpe *DistributedGrepEngine) Execute(gquery *grep.GrepQuery) {
	var outputsJson = make([]grep.GrepOutput, 0)

	start := time.Now()
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
		currConn := dpe.clientConns[i]
		if dpe.activeClients[generateClientConnKey(currConn)] == true {
			go dpe.remoteExecute(gquery, currConn, peerChannels[peerChannelIdx])
			peerChannelIdx += 1
		}
	}

	// * NOTE: localExecute() and remoteExecute() will not exit until its respective channels are read from since the channels
	// * once written to will block until someone reads from them. Therefore, it will block until it is read from below

	// Print local grep output to stdout
	grepOut := <-localChannel
	totalNumLines += grepOut.NumLines

	outputsJson = append(outputsJson, *grepOut)
	fmt.Print(grepOut.ToString())

	// Print peer grep outputs to stdout
	for i := 0; i < len(peerChannels); i++ {
		grepOut := <-peerChannels[i] // read from channel into grep outputs array
		totalNumLines += grepOut.NumLines

		outputsJson = append(outputsJson, *grepOut)
		fmt.Print(grepOut.ToString())
	}

	end := time.Now()

	if dpe.testOutputFileNameFormat != "" {
		_, err := dpe.CreateJson(gquery.PackagedString, outputsJson)
		if err != nil {
			fmt.Println("Error in creating json file ")
		}
	}

	elapsed := end.Sub(start)
	fmt.Printf("Total Number of Lines: %d\n", totalNumLines)
	fmt.Printf("Elapsed Query Execution Time: %dns\n\n", elapsed.Nanoseconds())
}

func (dpe *DistributedGrepEngine) CreateJson(packagedString string, outputsJson []grep.GrepOutput) ([]byte, error) {
	data := JSONOutput{
		Query:   packagedString,
		Outputs: outputsJson,
	}

	dataBytes, err := json.MarshalIndent(data, "", " ")

	if err != nil {
		fmt.Println("Error writing to file using json.Marshal")
	}
	filename := fmt.Sprintf(dpe.testOutputFileNameFormat, dpe.currentTestFileIdx)
	dpe.currentTestFileIdx += 1
	err = os.WriteFile(filename, dataBytes, os.FileMode(0644))
	if err != nil {
		fmt.Println("Error in Writing Creating Json file")
	}

	return dataBytes, err
}

func DeserializeJson(jsonFileName string) (string, []grep.GrepOutput) {
	jsonFile, err := os.Open(jsonFileName)
	if err != nil {
		log.Fatalf("Failed to open json file")
	}
	defer func(jsonFile *os.File) {
		_ = jsonFile.Close()
	}(jsonFile)

	dataBytes, _ := io.ReadAll(jsonFile)

	var jsonOutput JSONOutput
	err2 := json.Unmarshal(dataBytes, &jsonOutput)

	if err2 != nil {
		log.Fatalf("Error in deserializing json file")
	}

	return jsonOutput.Query, jsonOutput.Outputs
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
	grepOutput := dpe.checkCacheOrExecute(gquery)
	outputChannel <- grepOutput
}

func (dpe *DistributedGrepEngine) Shutdown() {
	dpe.StopServer()
	// TODO: anything else needed here?
}

func (dpe *DistributedGrepEngine) StopServer() {
	if dpe.serverQuit != nil {
		close(dpe.serverQuit)
	}

	err := dpe.listener.Close()
	if err != nil {
		log.Fatal("Failed to close server's listener object")
	}
	dpe.serverWg.Wait()
}

// When a client was disconnected, call this function to remove
// the client information from the DistributedGrepEngine struct
func (dpe *DistributedGrepEngine) removeClient(conn net.Conn) {
	dpe.activeClients[generateClientConnKey(conn)] = false
	dpe.numActiveClients -= 1
}

// Generate a key for a connection object for the client
func generateClientConnKey(conn net.Conn) string {
	remote_addr := conn.RemoteAddr().String()
	ip, _, _ := net.SplitHostPort(remote_addr)
	return ip
}
