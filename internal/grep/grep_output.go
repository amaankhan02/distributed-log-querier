package grep

import (
	"bytes"
	"encoding/gob"
)

// TODO: should any of these fields be pointers?
type GrepOutput struct {
	output   string
	filename string
	numLines int
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
