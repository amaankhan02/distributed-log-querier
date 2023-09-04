package grep

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// -------------------------- GrepOutput -------------------------- \\

// TODO: should any of these fields be pointers?
type GrepOutput struct {
	byteOutput []byte // output data as an array of bytes
	filename   string
}

// Return the output as a string
func (grepOutput GrepOutput) ToString() string {
	return string(grepOutput.byteOutput) // TODO: do i return referece/copy? should i allocate on heap?
}

// -------------------------- GrepQuery -------------------------- \\

// Represents details of the grep query a user my type in. Contains all information of that
// specific query, including any functions to execute the query or convert to a different form
// GrepQuery is independent of the filename, therefore the cmdArgs field does not contain the filename
type GrepQuery struct {
	cmdArgs        []string // slice of the command line arguments (w/o the filename)
	packagedString string   // command args as one string concatenated by "-" b/w each arg
}

func CreateGrepQueryFromInput(rawUserInput string) (GrepQuery, error) {
	var g = GrepQuery{}
	query, err := parseRawGrepQuery(rawUserInput)
	if err != nil {
		return g, err
	}
	g.cmdArgs = query
	g.packagedString = strings.Join(g.cmdArgs, "-")
	return g, nil
}

// Given a packagedString (grep query with cmd args split by "-") it returns a GrepQuery object
func CreateGrepQueryFromPackagedString(packagedString string) GrepQuery {
	var g = GrepQuery{}
	g.cmdArgs = strings.Split(packagedString, "-")
	g.packagedString = packagedString
	return g
}

// Executes the grep query on the file provided, and returns a GrepOutput object
func (q GrepQuery) Execute(filename string) GrepOutput {
	// make last arg the file to search -> which will be the log file for machine
	cmdLineArgs := append(q.cmdArgs, filename)
	cmd := exec.Command(cmdLineArgs[0], cmdLineArgs[1:]...) //Define grep command to run, store in cmd

	output, err := cmd.CombinedOutput() // run command and capture its output

	// make sure there were matches in doing this
	if err != nil { // TODO: why do we have to make sure there were matches tho?
		fmt.Println("No matches") // TODO: ask samaah what to do instead cuz we shouldn't print in this
	}

	var gOut = GrepOutput{output, filename}
	return gOut // TODO: does this return a copy or reference?? do i need to use pointers??

	// ! don't return the string, store the byte[] instead since strings are immutable
	//result := string(output)             // convert the output from bytes to strings. // TODO: why do we need it as a string?
	//lines := strings.Count(result, "\n") // number of lines found

	//return result, lines
}

// Parses the grep query user entered. Returns a slice containing the individual command arguments
func parseRawGrepQuery(userInput string) ([]string, error) {
	// split user input into command and arguments
	cmdArgs := strings.Fields(userInput) // split based off spaces
	cmdArgs = handleExtraQuotes(cmdArgs)

	// Make sure the user provided atleast two arguments
	if len(cmdArgs) < 2 {
		return nil, errors.New("Invalid input! Length of command arguments too small or invalid")
	}

	return cmdArgs, nil
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
