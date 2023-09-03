package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"internal"
)

const MACHINE_NAME_FORMAT = "fa23-cs425-19%02d.cs.illinois.edu"
const NUM_MACHINES = 10
const PORT = "9000"

func main() {
	/*
		Create loop to
			a) Show prompt to user
			b) Read user's grep query (user input) OR wait for query to be recieved from another machine
			...
	*/

	thisMachine, peerMachines := getMachineNames(MACHINE_NAME_FORMAT, NUM_MACHINES)

	// used to wait for all incoming connections to finish
	var wg sync.WaitGroup // WaitGroup similar to using pthread_join (later though - this just initializes)

	// Start server on port
	listener, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		utils.exitOnError("Listen()", err)
	}

	defer listener.Close() // closes listener right before exiting function

	fmt.Println("Listening for incoming connections on port")

	// Accept incoming connections from other machines
	for {
		conn, err := listener.Accept()

	}
}
