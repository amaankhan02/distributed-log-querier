package main

import (
	"cs425_mp1/internal/distributed_engine"
	"cs425_mp1/internal/grep"
	"cs425_mp1/internal/utils"
	"encoding/gob"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	MACHINE_NAME_FORMAT = "fa23-cs425-19%02d.cs.illinois.edu"
	PORT_FORMAT         = "80%02d" // 8001, 8002, ... 8010 - based on the
)

var flagNumMachines *int
var localLogFile *string // full path of the local log file of this machine
var cacheSize *int
var verbose *bool

func ParseArguments() {
	flagNumMachines = flag.Int("n", 10, "Number of Machines in the network in the range [2, 10]")
	filename := flag.String("f", "", "Filename of the log file")
	cacheSize = flag.Int("c", 10, "Size of the in-memory LRU cache")
	verbose = flag.Bool("v", false, "Indicates if you want messages to be printed out")
	flag.Parse()

	// convert filename to a full path for the localLogFile variable
	fullPath, err := filepath.Abs(*filename)
	if err != nil {
		log.Fatal("Invalid file - does not exist!")
	}
	localLogFile = &fullPath
}

func Init() {
	gob.Register(&grep.GrepQuery{})
	gob.Register(&grep.GrepOutput{})
}

func main() {
	ParseArguments()
	Init()

	peerServerAddresses := utils.GetPeerServerAddresses(MACHINE_NAME_FORMAT, PORT_FORMAT, *flagNumMachines)
	serverPort := utils.GetLocalhostPort(MACHINE_NAME_FORMAT, PORT_FORMAT, *flagNumMachines)
	engine := distributed_engine.CreateEngine(*localLogFile, serverPort, peerServerAddresses, *cacheSize, *verbose)

	_, _ = fmt.Fprintln(os.Stderr, "Setting up server. Listening to new connections...")
	engine.InitializeServer()
	utils.PromptWait("Wait for all machines to initialize and setup server") // prompting user to wait
	engine.ConnectToPeers()

	for {
		utils.PrintMessage("Enter Grep command:", *verbose)
		rawInput, input_err := utils.ReadUserInput()

		if rawInput == "exit" || input_err == io.EOF {
			engine.Shutdown()
			break
		}

		grepQuery, err := grep.CreateGrepQueryFromInput(rawInput)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}
		engine.Execute(grepQuery)
	}
}
