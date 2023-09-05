package network

import (
	"encoding/binary"
	"net"
)

/*
Send request in the format
Format: [size][data]

	[size] is the size of the data represented in a binary format - 4 Byte little-endian
	[data] is a []byte of the serialize GrepQuery/GrepOutput object gquery (use grep.SerializeGrepOutput())
*/
func SendRequest(data []byte, conn net.Conn) error {
	size := len(data)

	err := sendMessageSize(size, conn)
	if err != nil {
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func sendMessageSize(base10Number int, conn net.Conn) error {
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(base10Number))
	_, err := conn.Write(size) // _ is the number of bytes sent

	if err != nil {
		return err
	}

	return nil
}
