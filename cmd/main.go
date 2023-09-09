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
	"path/filepath"
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
	//NUM_MACHINES        = 2        // num of total machines in the network, although you should be able to use less
	PORT_FORMAT = "80%02d" // 8001, 8002, ... 8010 - based on the
)

var flagNumMachines *int
var localLogFile *string // full path of the local log file of this machine

func ParseArguments() {
	flagNumMachines = flag.Int("m", 10, "Number of Machines in the network in the range [2, 10]")
	filename := flag.String("f", "", "Filename of the log file")
	flag.Parse()

	fullPath, err := filepath.Abs(*filename)
	if err != nil {
		log.Fatal("Invalid file - does not exist!")
	}
	localLogFile = &fullPath
}

func main() {
	ParseArguments()

	gob.Register(&grep.GrepQuery{})
	gob.Register(&grep.GrepOutput{})

	peerServerAddresses := utils.GetPeerServerAddresses(MACHINE_NAME_FORMAT, PORT_FORMAT, *flagNumMachines)
	serverPort := utils.GetLocalhostPort(MACHINE_NAME_FORMAT, PORT_FORMAT, *flagNumMachines)
	//localLogFile := utils.GetLocalLogFile()
	engine := distributed_engine.CreateEngine(*localLogFile, serverPort, peerServerAddresses)

	fmt.Println("Setting up server. Listening to new connections...")
	engine.InitializeServer()
	utils.PromptWait("Wait for all machines to initialize and setup server") // prompting user to wait
	engine.ConnectToPeers()

	for {
		utils.DisplayGrepPrompt()
		rawInput, input_err := utils.ReadUserInput()

		if rawInput == "exit" || input_err == io.EOF {
			engine.Shutdown()
		}

		grepQuery, err := grep.CreateGrepQueryFromInput(rawInput)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		engine.Execute(grepQuery)
	}
}
