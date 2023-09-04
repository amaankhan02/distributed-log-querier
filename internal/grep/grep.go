package grep

import (
	"fmt"
	"os/exec"
	"strings"
)

// Represents details of the grep query a user my type in. Contains all information of that
// specific query, including any functions to execute the query or convert to a different form
// GrepQuery is independent of the filename, therefore the cmdArgs field does not contain the filename
type GrepQuery struct {
	cmdArgs        []string // slice of the command line arguments (w/o the filename)
	packagedString string   // command args as one string concatenated by "-" b/w each arg
}

// TODO: should any of these fields be pointers?
type GrepOutput struct {
	byteOutput []byte // output data as an array of bytes
	filename   string
}

func (q GrepQuery) Execute(filename string) GrepOutput {
	// make last arg the file to search -> which will be the log file for machine
	cmdLineArgs := append(q.cmdArgs, filename)
	cmd := exec.Command(cmdLineArgs[0], cmdLineArgs[1:]...) //Define grep command to run, store in cmd

	output, err := cmd.CombinedOutput() // run command and capture its output

	// make sure there were matches in doing this
	if err != nil { // TODO: why do we have to make sure there were matches tho?
		fmt.Println("No matches") // TODO: ask samaah what to do instead cuz we shouldn't print in this
	}

	var gOut GrepOutput = 

	// ! don't return the string, store the byte[] instead since strings are immutable
	//result := string(output)             // convert the output from bytes to strings. // TODO: why do we need it as a string?
	//lines := strings.Count(result, "\n") // number of lines found

	//return result, lines
}

// TODO: ask samaah what exactly this does and where it's used
// Converts the grep commands in an array into one string split by dashes
func PackageGrepCommands(cmdArgs []string) string {
	key := strings.Join(cmdArgs, "-")
	return key
}

// Unpack the formatted grep command in a string split by dashes back into an array of cmd args
func UnpackGrepCommands(grepCommand string) []string {
	cmdArgs := strings.Split(grepCommand, "-")
	return cmdArgs
}

// Helper function to be used in ParseRawGrepQuery
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

// Parses the grep query user entered. Returns a slice containing the individual command arguments
func ParseRawGrepQuery(userInput string) []string {
	// split user input into command and arguments
	cmdArgs := strings.Fields(userInput) // split based off spaces
	cmdArgs = handleExtraQuotes(cmdArgs)

	// Make sure the user provided atleast two arguments
	if len(cmdArgs) < 2 {
		fmt.Println("Invalid input. Please provide a valid grep command.")
		return nil
	}

	return cmdArgs
}
