package internal

import (
	"fmt"
	"os"
)

// Utility function that prints a formatted error message to stderr
// and exits the program with exit code 1
// 
// Parameters:
// 		messageContext (string): context of where this error is occurring
func exitOnError(messageContext string, err error) {
	messageFormat := "--------Error--------\n\tContext: %s\n\tError Message: %s"
	message = fmt.Sprintf(messageFormat, messageContext, err.Error())
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

// Returns a tuple of the current machine's hostname and an array of all the machine names
func getMachineNames(machineNameFormat string, numMachines int) (string, []string) {
	thisMachineName, err := os.Hostname()
	if err != nil
		utils.exitOnError("Retrieving HostName in utils.getMachineNames()", err)

	machineNames = make([]string, 0)
	
	for i := 1; i <= numMachines; i++ {
		newName = fmt.Sprintf(machineNameFormat, i)
		if newName != thisMachineName
			machineNames = append(machineNames, newName)
	}

	return thisMachineName, machineNames
}