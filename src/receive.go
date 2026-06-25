package src

import (
	"fmt"
	"io"
	"net"
	"os"
)

func newConnector(metadata Metadata) *net.Conn {
	target := metadata.Host + ":" + metadata.Port
	conn, err := net.Dial("tcp4", target)
	if err != nil {
		panic("Error stabilishing connection with sender at " + target)
	}
	return &conn
}

func Receive(metadata Metadata) {
	file, err := os.Create(metadata.Fname)
	if err != nil {
		panic("Error creating file " + metadata.Fname)
	}
	defer file.Close()

	offset, err := file.Seek(0, 0)
	if err != nil || offset != 0 {
		panic("Error opening the file " + metadata.Fname)
	}

	conn := newConnector(metadata)
	defer (*conn).Close()

	reader := io.Reader(*conn)
	buffer := make([]byte, packageSize)

	fmt.Println("Starting file transfer")

	for {
		n, err := io.ReadFull(reader, buffer)
		if err != nil && err != io.EOF {
			conn = newConnector(metadata)
			continue
		}

		_, err = file.Write(buffer)
		if err != nil {
			panic("Error writing to file")
		}

		clear(buffer)

		if n != 1024 {
			break
		}
	}

	fmt.Println("File received successfully and stored at " + metadata.Fname)
}
