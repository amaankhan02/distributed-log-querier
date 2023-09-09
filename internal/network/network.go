package network

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
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

Returns an error of io.EOF or io.ErrUnexpectedEOF
*/
// TODO: change to use io.ReadFull() w/ bufio.Reader()
// TODO: change to return and handle the error = io.EOF to know when we are done w/ the port to exit out? not sure...
func ReadRequest(conn net.Conn) ([]byte, error) {
	reader := bufio.NewReader(conn)
	data_size, err := readMessageSize(reader, MESSAGE_SIZE_BYTES)

	if err != nil {
		return nil, err
	}
	buff := make([]byte, data_size)

	// ReadFull() reads exactly len(buff) bytes from reader into buff
	// returns the number of bytes copied, and an error if fewer than len(buff) bytes were read
	// the error is io.EOF only if no bytes were read. otherwise its an ErrUnexpectedEOF
	_, err = io.ReadFull(reader, buff)
	if err != nil { // fewer than len(buff) bytes were read, or none at all
		return nil, err
	}

	return buff, nil
}

/*
Helper function to read just the message size from the connection
*/
func readMessageSize(reader *bufio.Reader, messageSizeBytes int) (int, error) {
	if messageSizeBytes != 4 && messageSizeBytes != 8 {
		return 0, errors.New("Invalid argument for messageSizeBytes - must be either equal to 4 or 8")
	}

	buff := make([]byte, messageSizeBytes)
	_, err := io.ReadFull(reader, buff)
	if err != nil {
		return 0, err
	}

	if messageSizeBytes == 4 {
		return int(binary.LittleEndian.Uint32(buff)), nil
	} else {
		return int(binary.LittleEndian.Uint64(buff)), nil
	}
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
