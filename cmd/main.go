package main

import (
	"cs425_mp1/internal/distributed_engine"
	"cs425_mp1/internal/grep"
	"cs425_mp1/internal/utils"
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	MACHINE_NAME_FORMAT = "fa23-cs425-19%02d.cs.illinois.edu"
	PORT_FORMAT         = "80%02d" // 8001, 8002, ... 8010 - based on the
	OUTPUT_JSON_FORMAT  = "test%d.json"
)

var flagNumMachines *int
var localLogFile *string // full path of the local log file of this machine
var cacheSize *int
var verbose *bool

var peerServerAddresses []string
var engine *distributed_engine.DistributedGrepEngine
var serverPort string

var testDir *string

func ParseArguments() {
	flagNumMachines = flag.Int("n", 10, "Number of Machines in the network in the range [2, 10]")
	localLogFile = flag.String("f", "", "Filename of the log file")
	cacheSize = flag.Int("c", 10, "Size of the in-memory LRU cache")
	verbose = flag.Bool("v", false, "Indicates if you want messages to be printed out")
	testDir = flag.String("t", "", "If you wish to run this program in TEST mode, put the directory you want your output JSON files to be stored")
	flag.Parse()
}

func Init() {
	gob.Register(&grep.GrepQuery{})
	gob.Register(&grep.GrepOutput{})

	peerServerAddresses = utils.GetPeerServerAddresses(MACHINE_NAME_FORMAT, PORT_FORMAT, *flagNumMachines)
	serverPort = utils.GetLocalhostPort(MACHINE_NAME_FORMAT, PORT_FORMAT, *flagNumMachines)
	if *testDir != "" {
		_, _ = fmt.Fprintf(os.Stderr, "Opening in [TEST] mode. Saving test output JSON files to %s\n", *testDir)
		dirPlusFile := filepath.Join(*testDir, OUTPUT_JSON_FORMAT)
		engine = distributed_engine.CreateEngine(*localLogFile, serverPort, peerServerAddresses, *cacheSize, *verbose, dirPlusFile)
	} else {
		engine = distributed_engine.CreateEngine(*localLogFile, serverPort, peerServerAddresses, *cacheSize, *verbose, "")
	}
}

func ProcessInput() (string, error) {
	var inputStr string
	//if readFromFile {
	utils.PrintMessage("Enter Grep command:", *verbose)
	rawInput, input_err := utils.ReadUserInput()
	if rawInput == "exit" || input_err == io.EOF {
		return "", errors.New("Break")
	}
	inputStr = rawInput
	return inputStr, nil
}

func SetupEngine() {
	_, _ = fmt.Fprintln(os.Stderr, "Setting up server. Listening to new connections...")
	engine.InitializeServer()
	engine.ConnectToPeers()
	_, _ = fmt.Fprintln(os.Stderr, "Connected to all machines")
}

func main() {
	ParseArguments()
	Init()
	SetupEngine()

	for {
		inputStr, err2 := ProcessInput()
		if err2 != nil {
			//engine.Shutdown()
			break
		}
		grepQuery, err := grep.CreateGrepQueryFromInput(inputStr)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}
		engine.Execute(grepQuery)
	}
}
