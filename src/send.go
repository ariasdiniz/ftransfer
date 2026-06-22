package src

import (
	"fmt"
	"io"
	"net"
	"os"
)

func newListener(meta Metadata) *net.Conn {
	sock, err:= net.Listen("tcp4", meta.Host + ":" + meta.Port)
	if err != nil {
		panic("Could not set up TCP socket listening to port " + meta.Port)
	}
	conn, err := sock.Accept()
	if err != nil {
		panic("Error stabilishing connection")
	}
	return &conn
}

func Send(metadata Metadata) {
  file, err := os.Open(metadata.Fname)
	if err != nil {
		panic("Could not open file " + metadata.Fname)
	}

	sock := *newListener(metadata)
	defer sock.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			panic("Error reading file from disk.")
		}

		if n == 0 || err == io.EOF {
			break
		}
		n, err = sock.Write(buffer)
		if n == 0 || err != nil {
			panic("Error transfering bytes")
		}
	}

	fmt.Println("File sent successfully!")
}
