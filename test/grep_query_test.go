package test

import (
	"cs425_mp1/internal/grep"
	"testing"
)

func TestExecuteGrepSimple(t *testing.T) {
	var expectedOutput = "This is a sample text file for testing the grep command in Go.\n"
	var expectedNumLines = 1

	q, err := grep.CreateGrepQueryFromInput("grep sample")

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// Call the Execute function from the grep package
	filename := "C:\\Users\\samaa\\Documents\\2023-2024\\DistributedSystems\\MP1\\cs425_mp1\\test\\test_logs\\test_text_file1.txt" // Replace with the actual file path
	grepOutput := q.Execute(filename)

	lines := grepOutput.NumLines
	output := grepOutput.Output

	if lines != expectedNumLines {
		t.Errorf("Expected %d number of lines, but got %d", expectedNumLines, lines)
	}

	if output != expectedOutput {
		t.Errorf("Expected output: %s, but got %s", expectedOutput, output)
	}
}

func TestExecuteEndWithCmd(t *testing.T) {
	var expectedOutput = "2023-09-06 22:52:35,317 INFO: Cache size: 256 MB\n2023-09-06 22:52:35,319 WARNING: Low memory warning, available memory: 10 MB\n"
	var expectedNumLines = 2

	q, err := grep.CreateGrepQueryFromInput("grep MB$")

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// Call the Execute function from the grep package
	filename := "test_logs/test_log_file5.log"
	grepOutput := q.Execute(filename)

	lines := grepOutput.NumLines
	output := grepOutput.Output

	if lines != expectedNumLines {
		t.Errorf("Expected %d number of lines, but got %d", expectedNumLines, lines)
	}

	if output != expectedOutput {
		t.Errorf("Expected output: %s, but got %s", expectedOutput, output)
	}
}

// Tests or command
func TestExecuteOrCmd(t *testing.T) {
	var expectedOutput = "WARNING: Configuration file outdated, please update\nINFO: Application started\nCRITICAL: Application halted: fatal error\n"
	var expectedNumLines = 3

	q, err := grep.CreateGrepQueryFromInput("grep 'Configuration\\|Application'")

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// Call the Execute function from the grep package
	filename := "test_logs/test_log_file4.log"
	grepOutput := q.Execute(filename)

	lines := grepOutput.NumLines
	output := grepOutput.Output

	if lines != expectedNumLines {
		t.Errorf("Expected %d number of lines, but got %d", expectedNumLines, lines)
	}

	if output != expectedOutput {
		t.Errorf("Expected output: %s, but got %s", expectedOutput, output)
	}
}

// tests "-c" arg that counts number of matches
// tests handle with quotes
func TestExecuteNumMatchesCmd(t *testing.T) {
	var expectedOutput = "11\n"
	var expectedNumLines = 1

	q, err := grep.CreateGrepQueryFromInput("grep -c \"CRITICAL:\"")

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// Call the Execute function from the grep package
	filename := "test_logs/test_log_file6.log"
	grepOutput := q.Execute(filename)

	lines := grepOutput.NumLines
	output := grepOutput.Output

	if lines != expectedNumLines {
		t.Errorf("Expected %d number of lines, but got %d", expectedNumLines, lines)
	}

	if output != expectedOutput {
		t.Errorf("Expected output: %s, but got %s", expectedOutput, output)
	}
}

func TestExecuteCaseInsensitiveCmd(t *testing.T) {
	var expectedOutput = "2023-09-06 22:52:35,317 DEBUG: Received API request: ...\n2023-09-06 22:52:35,319 DEBUG: Received API request: ...\n"
	var expectedNumLines = 2

	q, err := grep.CreateGrepQueryFromInput("grep -i \"api ReQUEST\"")

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// Call the Execute function from the grep package
	filename := "test_logs/test_log_file5.log"
	grepOutput := q.Execute(filename)

	lines := grepOutput.NumLines
	output := grepOutput.Output

	if lines != expectedNumLines {
		t.Errorf("Expected %d number of lines, but got %d", expectedNumLines, lines)
	}

	if output != expectedOutput {
		t.Errorf("Expected output: %s, but got %s", expectedOutput, output)
	}
}

func TestExecuteLinesBeforeCmd(t *testing.T) {
	var expectedOutput = "2023-09-06 22:54:16,273 INFO: File saved successfully\n2023-09-06 22:54:16,278 INFO: Request received from client\n"
	var expectedNumLines = 2

	q, err := grep.CreateGrepQueryFromInput("grep -B1 \"INFO: Request received from client\"")

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// Call the Execute function from the grep package
	filename := "test_logs/test_log_file6.log"
	grepOutput := q.Execute(filename)

	lines := grepOutput.NumLines
	output := grepOutput.Output

	if lines != expectedNumLines {
		t.Errorf("Expected %d number of lines, but got %d", expectedNumLines, lines)
	}

	if output != expectedOutput {
		t.Errorf("Expected output: %s, but got %s", expectedOutput, output)
	}
}

func TestExecuteComplexRegEx1(t *testing.T) {
	var expectedOutput = "CRITICAL: Server outage: emergency shutdown\nERROR: Internal server error: Unable to process request\nERROR: File not found: 'file.txt'\nCRITICAL: Application halted: fatal error\n"
	var expectedNumLines = 4

	q, err := grep.CreateGrepQueryFromInput("grep ^ERROR\\|^CRITICAL")

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// Call the Execute function from the grep package
	filename := "test_logs/test_log_file4.log"
	grepOutput := q.Execute(filename)

	lines := grepOutput.NumLines
	output := grepOutput.Output

	if lines != expectedNumLines {
		t.Errorf("Expected %d number of lines, but got %d", expectedNumLines, lines)
	}

	if output != expectedOutput {
		t.Errorf("Expected output: %s, but got %s", expectedOutput, output)
	}
}

func TestExecuteComplexRegEx2(t *testing.T) {
	var expectedOutput = "WARNING: Invalid input detected\nWARNING: Network connection unstable\nERROR: Internal server error: Unable to process request\nWARNING: File deletion warning, 'temp.txt' will be removed\nWARNING: Invalid input detected\nWARNING: Configuration file outdated, please update\nWARNING: Invalid input detected\nWARNING: SSL certificate expiration warning, renew needed\nWARNING: SSL certificate expiration warning, renew needed\nERROR: File not found: 'file.txt'\n"
	var expectedNumLines = 10

	q, err := grep.CreateGrepQueryFromInput("grep -E ^[eE]|^W")

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// Call the Execute function from the grep package
	filename := "test_logs/test_log_file4.log"
	grepOutput := q.Execute(filename)

	lines := grepOutput.NumLines
	output := grepOutput.Output

	if lines != expectedNumLines {
		t.Errorf("Expected %d number of lines, but got %d", expectedNumLines, lines)
	}

	if output != expectedOutput {
		t.Errorf("Expected output: %s, but got %s", expectedOutput, output)
	}
}

func TestExecuteComplexRegEx3(t *testing.T) {
	var expectedOutput = "WARNING: Invalid input detected\nWARNING: Network connection unstable\nINFO: API version 1.2.3 is now deprecated\nERROR: Internal server error: Unable to process request\nWARNING: File deletion warning, 'temp.txt' will be removed\nWARNING: Invalid input detected\nWARNING: Configuration file outdated, please update\nINFO: Database initialized successfully\nINFO: Application started\nWARNING: Invalid input detected\nWARNING: SSL certificate expiration warning, renew needed\nWARNING: SSL certificate expiration warning, renew needed\nINFO: Scheduled maintenance task started\nINFO: Received 200 OK response from API endpoint"
	var expectedNumLines = 14

	q, err := grep.CreateGrepQueryFromInput("grep -e ^[F-Z] -e request$")

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// Call the Execute function from the grep package
	filename := "test_logs/test_log_file4.log"
	grepOutput := q.Execute(filename)

	lines := grepOutput.NumLines
	output := grepOutput.Output

	if lines != expectedNumLines {
		t.Errorf("Expected %d number of lines, but got %d", expectedNumLines, lines)
	}

	if output != expectedOutput {
		t.Errorf("Expected output: %s, but got %s", expectedOutput, output)
	}
}

func TestExecuteComplexRegEx4(t *testing.T) {
	var expectedOutput = ""
	var expectedNumLines = 0

	q, err := grep.CreateGrepQueryFromInput("grep ^$")

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// Call the Execute function from the grep package
	filename := "test_logs/test_log_file6.log"
	grepOutput := q.Execute(filename)

	lines := grepOutput.NumLines
	output := grepOutput.Output

	if lines != expectedNumLines {
		t.Errorf("Expected %d number of lines, but got %d", expectedNumLines, lines)
	}

	if output != expectedOutput {
		t.Errorf("Expected output: %s, but got %s", expectedOutput, output)
	}
}

func TestSerializeDeserializeQuery(t *testing.T) {
	gQuery := grep.GrepQuery{
		CmdArgs:        []string{"grep", "sample", "example_file_name.txt"},
		PackagedString: "grep;sample;example_file_name.txt",
	}

	gQueryBytes, err := grep.SerializeGrepQuery(&gQuery)

	if err != nil {
		t.Errorf("Error thrown in serializing Grep Query: %s", err)
	}

	gQueryDeserialized, err := grep.DeserializeGrepQuery(gQueryBytes)

	if gQuery.PackagedString != gQueryDeserialized.PackagedString {
		t.Errorf("Expected package string %s, but got %s", gQuery.PackagedString, gQueryDeserialized.PackagedString)
	} else {

		if len(gQuery.CmdArgs) != len(gQueryDeserialized.CmdArgs) {
			t.Errorf("Expected length of cmdArgs %d, but got %d", len(gQuery.CmdArgs), len(gQueryDeserialized.CmdArgs))
		}
		for i, cmd := range gQuery.CmdArgs {
			if cmd != gQueryDeserialized.CmdArgs[i] {
				t.Errorf("Expected cmd %s, but got %s", cmd, gQueryDeserialized.CmdArgs)
			}
		}
	}
}
