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

	lines := grepOutput.NumLines()
	output := grepOutput.Output()

	if lines != expectedNumLines {
		t.Errorf("Expected %d number of lines, but got %d", expectedNumLines, lines)
	}

	if output != expectedOutput {
		t.Errorf("Expected output: %s, but got %s", expectedOutput, output)
	}
}

func TestExecuteGrepEndWithCmd(t *testing.T) {
	var expectedOutput = "2023-09-06 22:52:35,317 INFO: Cache size: 256 MB\n2023-09-06 22:52:35,319 WARNING: Low memory warning, available memory: 10 MB\n"
	var expectedNumLines = 2

	q, err := grep.CreateGrepQueryFromInput("grep MB$")

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// Call the Execute function from the grep package
	filename := "test_logs/test_log_file5.log"
	grepOutput := q.Execute(filename)

	lines := grepOutput.NumLines()
	output := grepOutput.Output()

	if lines != expectedNumLines {
		t.Errorf("Expected %d number of lines, but got %d", expectedNumLines, lines)
	}

	if output != expectedOutput {
		t.Errorf("Expected output: %s, but got %s", expectedOutput, output)
	}
}

// Tests or command
func TestExecuteGrepOrCmd(t *testing.T) {
	var expectedOutput = "WARNING: Configuration file outdated, please update\nINFO: Application started\nCRITICAL: Application halted: fatal error\n"
	var expectedNumLines = 3

	q, err := grep.CreateGrepQueryFromInput("grep 'Configuration\\|Application'")

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// Call the Execute function from the grep package
	filename := "test_logs/test_log_file4.log"
	grepOutput := q.Execute(filename)

	lines := grepOutput.NumLines()
	output := grepOutput.Output()

	if lines != expectedNumLines {
		t.Errorf("Expected %d number of lines, but got %d", expectedNumLines, lines)
	}

	if output != expectedOutput {
		t.Errorf("Expected output: %s, but got %s", expectedOutput, output)
	}
}

// tests "-c" arg that counts number of matches
// tests handle with quotes
func TestExecuteGrepNumMatchesCmd(t *testing.T) {
	var expectedOutput = "11\n"
	var expectedNumLines = 1

	q, err := grep.CreateGrepQueryFromInput("grep -c \"CRITICAL:\"")

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// Call the Execute function from the grep package
	filename := "test_logs/test_log_file6.log"
	grepOutput := q.Execute(filename)

	lines := grepOutput.NumLines()
	output := grepOutput.Output()

	if lines != expectedNumLines {
		t.Errorf("Expected %d number of lines, but got %d", expectedNumLines, lines)
	}

	if output != expectedOutput {
		t.Errorf("Expected output: %s, but got %s", expectedOutput, output)
	}
}

func TestExecuteGrepCaseInsensitiveCmd(t *testing.T) {
	var expectedOutput = "2023-09-06 22:52:35,317 DEBUG: Received API request: ...\n2023-09-06 22:52:35,319 DEBUG: Received API request: ...\n"
	var expectedNumLines = 2

	q, err := grep.CreateGrepQueryFromInput("grep -w \"api ReQUEST\"")

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	// Call the Execute function from the grep package
	filename := "test_logs/test_log_file5.log"
	grepOutput := q.Execute(filename)

	lines := grepOutput.NumLines()
	output := grepOutput.Output()

	if lines != expectedNumLines {
		t.Errorf("Expected %d number of lines, but got %d", expectedNumLines, lines)
	}

	if output != expectedOutput {
		t.Errorf("Expected output: %s, but got %s", expectedOutput, output)
	}
}
