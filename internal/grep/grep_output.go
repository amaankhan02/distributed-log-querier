package grep

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// TODO: should any of these fields be pointers?
type GrepOutput struct {
	// variable names starting with lowercase indicates it's private access
	output   string
	filename string
	numLines int
}

// Getter for numLines attribute
func (g GrepOutput) NumLines() int {
	return g.numLines
}

// Getter for filename attribute
func (g GrepOutput) Filename() string {
	return g.filename
}

// Formats the contents of the GrepOutput as a string
func (g GrepOutput) ToString() string {
	strFormat := "Filename: %s\nNumber of Lines: %d\nOutput:%s\n"
	return fmt.Sprintf(strFormat, g.filename, g.numLines, g.output)
}

// SerializeGrepOutput Serialize GrepOutput object into a byte array
// The returned format is good to send over a TCP socket
// Returns nil if it failed to serialize
func SerializeGrepOutput(grepOutput GrepOutput) []byte {
	binary_buff := new(bytes.Buffer)

	encoder := gob.NewEncoder(binary_buff)
	err := encoder.Encode(grepOutput)
	if err != nil {
		return nil
	}
	return binary_buff.Bytes()
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
func DeserializeGrepOutput(data []byte) *GrepOutput {
	grepOutput := new(GrepOutput)
	byteBuffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(byteBuffer)

	err := decoder.Decode(grepOutput)
	if err != nil {
		return nil
	}

	return grepOutput
}
