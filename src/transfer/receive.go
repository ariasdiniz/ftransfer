package transfer

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

func newConnector(metadata Metadata) net.Conn {
	target := metadata.Host + ":" + metadata.Port
	conn, err := net.Dial("tcp4", target)
	if err != nil {
		panic("Error stabilishing connection with sender at " + target)
	}
	return conn
}

func readInt(conn net.Conn) uint16 {
	var fnameSize uint16
	buffer := make([]byte, 2)
	io.ReadFull(conn, buffer)
	fnameSize = binary.LittleEndian.Uint16(buffer)
	return fnameSize
}

func Receive(metadata Metadata) {
	conn := newConnector(metadata)

	reader := io.Reader(conn)
	buffer := make([]byte, packetSize)
	
	if metadata.Fname == "" {
	  fnameSize := readInt(conn)
		fname := make([]byte, fnameSize)
		io.ReadFull(conn, fname)
		metadata.Fname = filepath.Base(string(fname))
	}

	file, err := os.Create(metadata.Fname)
	if err != nil {
		panic("Error creating file " + metadata.Fname)
	}
	defer file.Close()

	offset, err := file.Seek(0, 0)
	if err != nil || offset != 0 {
		panic("Error opening the file " + metadata.Fname)
	}

	fmt.Printf("Starting file transfer, receiving %s\n", metadata.Fname)

	for {
		n, err := io.ReadFull(reader, buffer)
		if err != nil && err != io.EOF {
			conn.Close()
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

	conn.Close()
	fmt.Println("File received successfully and stored at " + metadata.Fname)
}
