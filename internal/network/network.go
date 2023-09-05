package network

import (
	"encoding/binary"
	"net"
)

const MESSAGE_SIZE_BYTES = 4 // number of bytes used in the protocol to define the size of the message

/*
Send request in the format
Format: [size][data]

	[size] is the size of the data represented in a binary format - 4 Byte little-endian
	[data] is a []byte of the serialize GrepQuery/GrepOutput object gquery (use grep.SerializeGrepOutput())
*/
func SendRequest(data []byte, conn net.Conn) error {
	size := len(data)

	err := sendMessageSize(size, conn, MESSAGE_SIZE_BYTES)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

/*
Read request from connection. Reads the message size and correctly gets all the []byte of data
and returns it. Caller is expected to deserialize this []byte of data as this function does not
do that.
*/
func ReadRequest(conn net.Conn) ([]byte, error) {

}

/*
Helper function to read just the message size from the connection
*/
func readMessageSize(conn net.Conn) error {

}

/*
Helper function to send just the message size as a MESSAGE_SIZE_BYTES number (like a 4byte number)
in little-endian format
*/
func sendMessageSize(base10Number int, conn net.Conn, messageSizeBytes int) error {
	size := make([]byte, MESSAGE_SIZE_BYTES)
	if messageSizeBytes == 4 {
		binary.LittleEndian.PutUint32(size, uint32(base10Number))
	} else if messageSizeBytes == 8 {
		binary.LittleEndian.PutUint64(size, uint64(base10Number))
	}

	_, err := conn.Write(size) // _ is the number of bytes sent

	if err != nil {
		return err
	}

	return nil
}
