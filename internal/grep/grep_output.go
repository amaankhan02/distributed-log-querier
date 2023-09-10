package grep

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"path/filepath"
	"time"
)

type GrepOutput struct {
	Output        string
	Filename      string
	NumLines      int
	ExecutionTime time.Duration
}

// Formats the contents of the GrepOutput as a string
func (g *GrepOutput) ToString() string {
	//dashesWithFilename := "------------------------%s------------------------\n"
	strFormat := "Filename: %s\nNum Lines: %d\nExecution Time: %dns\nOutput:\n%s\n"
	baseFileName := filepath.Base(g.Filename)
	return fmt.Sprintf(strFormat, baseFileName, g.NumLines, g.ExecutionTime.Nanoseconds(), g.Output)
}

// SerializeGrepOutput Serialize GrepOutput object into a byte array
// The returned format is good to send over a TCP socket
// Returns nil if it failed to serialize
func SerializeGrepOutput(grepOutput *GrepOutput) ([]byte, error) {
	binary_buff := new(bytes.Buffer)

	encoder := gob.NewEncoder(binary_buff)
	err := encoder.Encode(grepOutput)
	if err != nil {
		return nil, err
	}
	return binary_buff.Bytes(), nil
}

// DeserializeGrepOutput Deserializes the byte array into a GrepOutput object
// Parameters:
//
//	data: deserialized version of a GrepOutput struct.
//
// Returns
//
//	Pointer to a GrepOutput struct allocated on the heap.
//	nil if it failed to deserialized
func DeserializeGrepOutput(data []byte) (*GrepOutput, error) {
	grepOutput := new(GrepOutput)
	byteBuffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(byteBuffer)

	err := decoder.Decode(grepOutput)
	if err != nil {
		return nil, err
	}

	return grepOutput, nil
}

// Compares GrepOutput fields but does not compare execution time as that is not necessary for comparison in our cases
func GrepOutputsAreEqual(grepOutput1 *GrepOutput, grepOutput2 *GrepOutput) bool {
	return grepOutput1.Output == grepOutput2.Output && grepOutput1.NumLines == grepOutput2.NumLines && grepOutput1.Filename == grepOutput2.Filename
}
