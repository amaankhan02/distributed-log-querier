package internal

import (
	"fmt"
	"os/exec"
	"strings"
)

func ExecuteGrep(cmdArgs []string, logFileName string) (string, int) {
	// make last arg the file to search -> which will be the log file for machine
	cmdArgs = append(cmdArgs, logFileName)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...) //Define grep command to run, store in cmd

	output, err := cmd.CombinedOutput() // run command and capture its output

	// make sure there were matches in doing this
	if err != nil { // TODO: why do we have to make sure there were matches tho?
		fmt.Println("No matches")
	}

	result := string(output)             // convert the output from bytes to strings. // TODO: why do we need it as a string?
	lines := strings.Count(result, "\n") // number of lines found

	return result, lines
}

// TODO: ask samaah what exactly this does and where it's used
// create key (grepCommand) by concatenating vals in cmdArgs
func createKey(cmdArgs []string) string {
	key := strings.Join(cmdArgs, "-")
	return key
}

// TODO: ask samaah how this function is supposed to be used and where and why
func GrepCmdToCmdArgs(grepCommand string) []string {
	cmdArgs := strings.Split(grepCommand, "-")
	return cmdArgs
}
