package grep

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"path/filepath"
)

// TODO: should any of these fields be pointers?
type GrepOutput struct {
	// variable names starting with lowercase indicates it's private access
	Output   string
	Filename string
	NumLines int
}

// Formats the contents of the GrepOutput as a string
func (g *GrepOutput) ToString() string {
	dashesWithFilename := "------------------------%s------------------------\n"
	// dashes := "------------------------------------------------\n"
	strFormat := dashesWithFilename + "Filename: %s\nNumber of Lines: %d\nOutput:\n%s\n"
	baseFileName := filepath.Base(g.Filename)
	return fmt.Sprintf(strFormat, baseFileName, baseFileName, g.NumLines, g.Output)
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
