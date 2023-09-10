package test

import (
	"cs425_mp1/internal/grep"
	"testing"
)

func TestSerializeDeserializeOutput(t *testing.T) {
	gOutput := grep.GrepOutput{
		Output:   "Sample output from text file",
		Filename: "example_test_file1.txt",
		NumLines: 1,
	}

	gOutputBytes, err := grep.SerializeGrepOutput(&gOutput)

	if err != nil {
		t.Errorf("Error thrown in serializing Grep Output: %s", err)
	}

	gOutputDeserialized, err := grep.DeserializeGrepOutput(gOutputBytes)

	if gOutput.Output != gOutputDeserialized.Output {
		t.Errorf("Expected output %s, but got %s", gOutput.Output, gOutputDeserialized.Output)
	}

	if gOutput.Filename != gOutputDeserialized.Filename {
		t.Errorf("Expected filename %s, but got %s", gOutput.Filename, gOutputDeserialized.Filename)
	}

	if gOutput.NumLines != gOutputDeserialized.NumLines {
		t.Errorf("Expected filename %d, but got %d", gOutput.NumLines, gOutputDeserialized.NumLines)
	}
}
