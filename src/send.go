package src

import (
	"fmt"
	"io"
	"net"
	"os"
)

func newListener(metadata Metadata) *net.Listener {
	target := metadata.Host + ":" + metadata.Port
	sock, err := net.Listen("tcp4", target)
	fmt.Println("Ready to transfer file, awaiting connection")
	if err != nil {
		panic("Could not set up TCP socket listening to " + target)
	}
	return &sock
}

func acceptConn(metadata Metadata, sock *net.Listener) *net.Conn {
	conn, err := (*sock).Accept()
	target := metadata.Host + ":" + metadata.Port
	if err != nil {
		panic("Error stabilishing connection with " + target)
	}
	return &conn
}

func Send(metadata Metadata) {
	file, err := os.Open(metadata.Fname)
	if err != nil {
		panic("Could not open file " + metadata.Fname)
	}
	defer file.Close()

	sock := newListener(metadata)
	conn := *acceptConn(metadata, sock)
	buffer := make([]byte, packetSize)

	defer (*sock).Close()

	fmt.Println("Starting file transfer")

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			panic("Error reading file from disk.")
		}

		if n == 0 || err == io.EOF {
			break
		}

		n, err = conn.Write(buffer)
		if n == 0 || err != nil {
			conn.Close()
			conn = *acceptConn(metadata, sock)
			_, err = file.Seek(-packetSize, 1)
			if err != nil {
				conn.Close()
				panic("Error moving file offset")
			}
			continue
		}

		clear(buffer)
	}

	conn.Close()
	fmt.Println("File sent successfully!")
}
