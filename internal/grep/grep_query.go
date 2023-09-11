package grep

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Represents details of the grep query a user my type in. Contains all information of that
// specific query, including any functions to execute the query or convert to a different form
// GrepQuery is independent of the filename, therefore the cmdArgs field does not contain the filename
type GrepQuery struct {
	CmdArgs        []string // slice of the command line arguments (w/o the filename)
	PackagedString string   // command args as one string concatenated by "-" b/w each arg
}

const DELIMITER = ";"

func CreateGrepQueryFromInput(rawUserInput string) (*GrepQuery, error) {
	g := &GrepQuery{}
	query, err := parseRawGrepQuery(rawUserInput)
	if err != nil {
		return g, err
	}

	g.CmdArgs = query
	g.PackagedString = strings.Join(g.CmdArgs, DELIMITER)

	return g, nil
}

// Given a packagedString (grep query with cmd args split by "-") it returns a GrepQuery object
func CreateGrepQueryFromPackagedString(packagedString string) *GrepQuery {
	g := &GrepQuery{}
	g.CmdArgs = strings.Split(packagedString, DELIMITER)
	g.PackagedString = packagedString
	return g
}

func SerializeGrepQuery(gquery *GrepQuery) ([]byte, error) {
	binary_buff := new(bytes.Buffer)

	encoder := gob.NewEncoder(binary_buff)
	err := encoder.Encode(gquery)
	// fmt.Printf("\n--------- SERIALIZE GREP QUERY ---------\nbinary_buff: %v\nbinary_buff.Bytes(): %v\n", binary_buff, binary_buff.Bytes())
	if err != nil {
		return nil, err
	}
	return binary_buff.Bytes(), nil
}

func DeserializeGrepQuery(data []byte) (*GrepQuery, error) {
	gquery := new(GrepQuery)
	byteBuffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(byteBuffer)

	err := decoder.Decode(gquery)
	if err != nil {
		return nil, err
	}

	return gquery, nil
}

// Executes the grep query on the file provided, and returns a GrepOutput object
func (q *GrepQuery) Execute(filename string) *GrepOutput {
	// make last arg the file to search -> which will be the log file for machine
	start := time.Now()
	cmdLineArgs := append(q.CmdArgs, filename)
	cmd := exec.Command(cmdLineArgs[0], cmdLineArgs[1:]...) //Define grep command to run, store in cmd

	binaryOutput, err := cmd.CombinedOutput() // run command and capture its output

	// make sure there were matches in doing this
	if err != nil {
		return &GrepOutput{Output: "", Filename: filepath.Base(filename), NumLines: 0}
	}

	outputStr := string(binaryOutput)
	numLines := strings.Count(outputStr, "\n")
	end := time.Now()
	elapsedTime := end.Sub(start)

	return &GrepOutput{outputStr, filepath.Base(filename), numLines, elapsedTime}
}

// Parses the grep query user entered. Returns a slice containing the individual command arguments
func parseRawGrepQuery(userInput string) ([]string, error) {
	// split user input into command and arguments
	cmdArgs := strings.Fields(userInput) // split based off spaces
	cmdArgs = handleExtraQuotes(cmdArgs)

	// Make sure the user provided atleast two arguments
	if len(cmdArgs) < 2 {
		return nil, errors.New("Invalid input! Length of command arguments too small or invalid")
	} else if cmdArgs[0] != "grep" {
		return nil, errors.New("Invalid command! Must be a grep command w/o putting the filename")
	}

	return cmdArgs, nil
}

// Helper function to be used in ParseRawGrepQuery
// loop through and see if any of the command arguments start with quotations " or ' & handle that
func handleExtraQuotes(cmdArgs []string) []string {
	resultGrepQuery := []string{}
	modifiedCmd := ""

	for _, cmd := range cmdArgs {

		if strings.HasPrefix(cmd, `"`) && strings.HasSuffix(cmd, `"`) {

			cmd = strings.Trim(cmd, `"`)
			resultGrepQuery = append(resultGrepQuery, cmd)
		} else {

			if strings.HasPrefix(cmd, `"`) {
				modifiedCmd = cmd
			} else if strings.HasSuffix(cmd, `"`) {
				modifiedCmd = modifiedCmd + " " + cmd
				trimmedCmd := strings.Trim(modifiedCmd, `"`)
				resultGrepQuery = append(resultGrepQuery, trimmedCmd)
				modifiedCmd = ""
			} else if modifiedCmd != "" {
				modifiedCmd = modifiedCmd + " " + cmd
			} else {
				resultGrepQuery = append(resultGrepQuery, cmd)
			}
		}
	}

	return resultGrepQuery
}
