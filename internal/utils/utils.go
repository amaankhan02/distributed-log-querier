package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Get the port that is assigned for this computer/server
// Returns a string with a colon preceding the port number. Ex: ":8081"
// Returns an empty string due to failure/error
func GetLocalhostPort(machineNameFormat string, portFormat string, numMachines int) string {
	thisMachineName, err := os.Hostname() // to make sure we don't add our hostname as the peer
	if err != nil {
		panic(err)
	}

	for i := 1; i <= numMachines; i++ {
		machineName := fmt.Sprintf(machineNameFormat, i)
		if machineName == thisMachineName {
			port := fmt.Sprintf(portFormat, i)
			return ":" + port
		}
	}
	return ""
}

// Return a tuple where first value is the this machine's ip address concatenated with its
// respective port that is assigned to it, and then second value of the tuple is a slice of
// the ip addresses concatenated with their respective ports of all the peer servers
func GetPeerServerAddresses(machineNameFormat string, portFormat string, numMachines int) []string {
	thisMachineName, err := os.Hostname() // to make sure we don't add our hostname as the peer
	if err != nil {
		panic(err)
	}

	// addresses will be the peer's ip addresses w/ their respective port numbers
	peerAddresses := make([]string, 0)

	for i := 1; i <= numMachines; i++ {
		peerHostName := fmt.Sprintf(machineNameFormat, i)
		peerPort := fmt.Sprintf(portFormat, i)
		peerIP, err := net.LookupIP(peerHostName)

		if err != nil {
			fmt.Printf("Failed to resolve IP addresses for %s: %v\n", thisMachineName, err)
			continue
		}

		if peerHostName != thisMachineName {
			peerAddresses = append(peerAddresses, peerIP[0].String()+":"+peerPort)
		}
	}

	return peerAddresses
}

// Reads from stdin and trims any additional whitespace on sides and returns as a string
func ReadUserInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	userInput, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	userInput = strings.TrimSpace(userInput)
	return userInput, nil
}

func DisplayGrepPrompt() {
	fmt.Println("Enter grep command: ")
	fmt.Print("$ ")
}

func GetLocalLogFile() string {
	panic("Not implemented!")
}

// Prompts user to wait by displaying the "message" to stdout, and waits for user
// to type in "r" for continue. Anything else will simply display the prompt again
func PromptWait(message string) {
	fmt.Println(message)

	for {
		fmt.Println("Type [r] when you are ready to continue (w/o brackets)")
		input, _ := ReadUserInput()
		if input == "r" {
			break
		}
	}
}
