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

// var inputFile *string // used for testing --> reads from input file instead of stdin
var peerServerAddresses []string
var engine *distributed_engine.DistributedGrepEngine
var serverPort string

// var isTest *bool
var testDir *string

//var readFromFile bool
//var scanner *bufio.Scanner

func ParseArguments() {
	flagNumMachines = flag.Int("n", 10, "Number of Machines in the network in the range [2, 10]")
	localLogFile = flag.String("f", "", "Filename of the log file")
	cacheSize = flag.Int("c", 10, "Size of the in-memory LRU cache")
	verbose = flag.Bool("v", false, "Indicates if you want messages to be printed out")
	testDir = flag.String("t", "", "If you wish to run this program in TEST mode, put the directory you want your output JSON files to be stored")
	//inputFile = flag.String("t", "", "Used for unit testing. Input file for reading grep queries instead of stdin")
	flag.Parse()

	// convert filename to a full path for the localLogFile variable
	//fullPath, err := filepath.Abs(*filename)
	//if err != nil {
	//	log.Fatal("Invalid file - does not exist!")
	//}
	//localLogFile = &fullPath

	// convert inputFile to a full path for the
	//if *_input_filename != "" {
	//	fullInputFile, err2 := filepath.Abs(*_input_filename)
	//	if err2 != nil {
	//		log.Fatal("Invalid file for input file flag")
	//	}
	//	inputFile = &fullInputFile
	//} else {
	//	inputFile = _input_filename
	//}

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

	//if *localLogFile == "" {
	//	readFromFile = false
	//} else {
	//	readFromFile = true
	//	file, errf := os.Open(*inputFile)
	//	if errf != nil {
	//		log.Fatal("Error optning input file")
	//	}
	//	defer func(file *os.File) {
	//		_ = file.Close()
	//	}(file)
	//	scanner = bufio.NewScanner(file)
	//}
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
	//} else {
	//	canScan := scanner.Scan()
	//	if !canScan {
	//		return "", errors.New("Break")
	//	} else {
	//		inputStr = strings.TrimSpace(scanner.Text())
	//	}
	//}
	return inputStr, nil
}

func SetupEngine() {
	_, _ = fmt.Fprintln(os.Stderr, "Setting up server. Listening to new connections...")
	engine.InitializeServer()
	//utils.PromptWait("Wait for all machines to initialize and setup server") // prompting user to wait
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
