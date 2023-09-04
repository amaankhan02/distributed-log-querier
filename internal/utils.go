package internal

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
	//var thisMachineAddress string

	for i := 1; i <= numMachines; i++ {
		peerHostName := fmt.Sprintf(machineNameFormat, i)
		peerPort := fmt.Sprintf(portFormat, i)
		peerIP, err := net.LookupIP(peerHostName)

		if err != nil {
			fmt.Printf("Failed to resolve IP addresses for %s: %v\n", hostname, err)
			continue
		}

		if peerHostName != thisMachineName {
			peerAddresses = append(peerAddresses, peerIP[0].String()+":"+peerPort)
		}
		//else {
		//	thisMachineAddress = peerIP[0].String() + ":" + peerPort
		//}
	}

	return peerAddresses
}

// loop through and see if any of the command arguments start with quotoations " or ' & handle that
func handleExtraQuotes(cmdArgs []string) []string {
	var result []string
	var currentString string

	for _, cmd := range cmdArgs {

		if strings.HasPrefix(cmd, `"`) && strings.HasSuffix(cmd, `"`) {

			cmd = strings.Trim(cmd, `"`)
			result = append(result, cmd)
		} else {

			if strings.HasPrefix(cmd, `"`) {
				currentString = cmd
			} else if strings.HasSuffix(cmd, `"`) {
				currentString += " " + cmd
				result = append(result, strings.Trim(currentString, `"`))
				currentString = ""
			} else if currentString != "" {
				currentString += " " + cmd
			} else {
				result = append(result, cmd)
			}
		}
	}

	return result
}

func DisplayPrompt() {
	fmt.Println("Enter grep command: ")
	fmt.Print("$ ")
}

func GetUserInput() []string {
	// get input from user
	reader := bufio.NewReader(os.Stdin)
	userInput, _ := reader.ReadString('\n')
	userInput = strings.TrimSpace(userInput)

	// split user input into command and arguments
	cmdArgs := strings.Fields(userInput)

	cmdArgs = handleExtraQuotes(cmdArgs)

	// Make sure the user provided atleast two arguments
	if len(cmdArgs) < 2 {
		fmt.Println("Invalid input. Please provide a valid grep command.")
		return nil
	}

	return cmdArgs
}
