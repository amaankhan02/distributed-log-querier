package internal

import (
	"fmt"
	"os"
)

// Utility function that prints a formatted error message to stderr
// and exits the program with exit code 1
func exitOnError(err error) {
	message := "Error: " + err.Error() + "\n"
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

// Returns a tuple of the current machine's hostname and an array of all the machine names
func getMachineNames(machineNameFormat string, numMachines int) (string, []string) {
	thisMachineName, err := os.Hostname()
	if err != nil
		utils.exitOnError(err)

	machineNames = make([]string, 0)
	
	for i := 1; i <= numMachines; i++ {
		newName = fmt.Sprintf(machineNameFormat, i)
		if newName != thisMachineName
			machineNames = append(machineNames, newName)
	}

	return thisMachineName, machineNames
}