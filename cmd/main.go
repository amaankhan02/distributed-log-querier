package main

import (
	"cs425_mp1/internal/grep"
	"cs425_mp1/internal/grep_engine"
	"cs425_mp1/internal/utils"
	"fmt"
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

func main() {
	peerServerAddresses := utils.GetPeerServerAddresses(MACHINE_NAME_FORMAT, PORT_FORMAT, NUM_MACHINES)
	serverPort := utils.GetLocalhostPort(MACHINE_NAME_FORMAT, PORT_FORMAT, NUM_MACHINES)
	localLogFile := utils.GetLocalLogFile()
	engine := grep_engine.CreateEngine(localLogFile, serverPort, peerServerAddresses)

	fmt.Println("Setting up server. Listening to new connections...")
	engine.InitializeServer()
	utils.PromptWait("Wait for all machines to initialize and setup server") // prompting user to wait
	engine.ConnectToPeers()

	for {
		utils.DisplayGrepPrompt()
		rawInput := utils.ReadUserInput()

		grepQuery, err := grep.CreateGrepQueryFromInput(rawInput)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		engine.DistributedExecute(grepQuery)
	}
}
