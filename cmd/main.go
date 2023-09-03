package main

import (
	"cs425_mp1/internal" // adding "_" prevents having to write "internal." before each method call
	"fmt"
	"net"
	"sync"
)

const (
	MACHINE_NAME_FORMAT = "fa23-cs425-19%02d.cs.illinois.edu"
	NUM_MACHINES        = 10
	PORT                = "8080"
)

func main() {
	_, peerMachines := internal.GetMachineNames(MACHINE_NAME_FORMAT, NUM_MACHINES)
	ipAddresses := internal.GetIPAddresses(peerMachines)

	// used to wait for all incoming connections to finish
	var wg sync.WaitGroup // WaitGroup similar to using pthread_join (later though - this just initializes)

	for i, ipAddr := range ipAddresses {
		// Create TCP Listener for each IP address
		listener, err := net.Listen("tcp", ipAddr+":"+PORT) // TODO: check if this is right...
		if err != nil {
			fmt.Printf("Error creating listener for %s (Host: %s): %v\n", ipAddr, peerMachines[i], err)
			continue // skip this and just connect to the others
		}

		defer listener.Close()

		fmt.Printf("Server is listening on %s:%s\n", ipAddr, PORT)

		for {
			// Accept incoming connections
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Error accepting connection: ", err)
				continue
			}

			// Handle the connection in a separate goroutine
			wg.Add(1)
			go func() {
				defer wg.Done() // decrements by 1 after this func finishes
				handleConnection(conn)
			}()

		}
	}

	// Wait for all connections to finish - similar to pthread_join()
	wg.Wait()

}

func handleConnection(conn net.Conn) {
	// Handle the incoming connection here.
	// Read data from the peer and send responses as needed
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	// Example: Read data from the peer
	buffer := make([]byte, 1024)

	/*
		https://stackoverflow.com/a/27003111/7359915
			> NOTE: conn.Read() is non-blocking
			> To read a specific number of bytes, use io.ReadAtLeast or io.ReadFull
			> Or to read until some arbitrary condition is met, you should loop on the Read()
			  as long as there is no error (but error out on too-large inputs to prevent eating server resources)
			> If implementing a TEXT-Based protocol, consider net/textproto which puts bufio.Reader()
				in front of the connection, so you can read lines
	*/
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection: ", err)
		return
	}

	receivedData := buffer[:n]
	fmt.Printf("Received data from peer: %s\n", receivedData)

	// Example: Send a response back to the peer
	response := []byte("Hello from the server")
	_, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error sending response: ", err)
	}
}
