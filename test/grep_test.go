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
