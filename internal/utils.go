package internal

import (
	"fmt"
	"net"
	"os"
)

// DEPRECATED: realized we can just use panic() so might need this anymore
// Utility function that prints a formatted error message to stderr
// and exits the program with exit code 1
//
// Parameters:
//
//	messageContext (string): context of where this error is occurring
func ExitOnError(messageContext string, err error) {
	messageFormat := "--------Error--------\n\tContext: %s\n\tError Message: %s"
	message := fmt.Sprintf(messageFormat, messageContext, err.Error())
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

// Returns a tuple of the current machine's hostname and an array of all the machine names
func GetMachineNames(machineNameFormat string, numMachines int) (string, []string) {
	thisMachineName, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	machineNames := make([]string, 0)

	for i := 1; i <= numMachines; i++ {
		newName := fmt.Sprintf(machineNameFormat, i)
		if newName != thisMachineName {
			machineNames = append(machineNames, newName)
		}
	}

	return thisMachineName, machineNames
}

// Given a slice of hostnames, it returns a slice of their corresponding
// IP addresses
func GetIPAddresses(hostnames []string) []string {
	ret := make([]string, 0)
	for _, hostname := range hostnames {
		ipaddr, err := net.LookupIP(hostname)
		if err != nil {
			fmt.Printf("Failed to resolve IP addresses for %s: %v\n", hostname, err)
			continue
		}
		ret = append(ret, ipaddr[0].String()) // ipaddr[0] = ipv4, ipaddr[1] = ipv6
	}

	return ret
}
