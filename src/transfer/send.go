package transfer

import (
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"net"
	"os"
)

func newListener(metadata Metadata) net.Listener {
	target := metadata.Host + ":" + metadata.Port
	fmt.Println("Creating RSA certificate for data encryption")
	cert, err := CreateInMemoryCert()
	if err != nil {
		fmt.Println("Error during certificate creation. Cannot ensure connection security")
	}
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	sock, err := tls.Listen("tcp4", target, &config)
	fmt.Println("Ready to transfer file, awaiting connection")
	if err != nil {
		panic("Could not set up TCP socket listening to " + target)
	}
	return sock
}

func acceptConn(metadata Metadata, sock net.Listener) net.Conn {
	conn, err := sock.Accept()
	target := metadata.Host + ":" + metadata.Port
	if err != nil {
		panic("Error stabilishing connection with " + target)
	}
	return conn
}

func writeMetadata(conn net.Conn, fMetadata FileMetadataHeader) {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, fMetadata.FnameSize)
	conn.Write(buffer)

	conn.Write([]byte(fMetadata.Fname))

	clear(buffer)
	binary.LittleEndian.PutUint64(buffer, fMetadata.Fsize)
	conn.Write(buffer)

}

func Send(metadata Metadata) {
	file, err := os.Open(metadata.Fname)
	if err != nil {
		panic("Could not open file " + metadata.Fname)
	}
	defer file.Close()

	stat, _ := os.Stat(metadata.Fname)

	sock := newListener(metadata)
	conn := acceptConn(metadata, sock)
	buffer := make([]byte, packetSize)

	defer sock.Close()

	fmt.Println("Starting file transfer")
	fMetadata := FileMetadataHeader{
		FnameSize: uint64(len(metadata.Fname)),
		Fname:     metadata.Fname,
		Fsize:     uint64(stat.Size()),
	}

	totalPackets := uint64(math.Floor(float64(fMetadata.Fsize / packetSize)))

	writeMetadata(conn, fMetadata)
	fmt.Println("------------------------------------")
	fmt.Printf("Transfering file: %s\n", fMetadata.Fname)
	fmt.Printf("File size: %d bytes\n", fMetadata.Fsize)

	bPacketNumber := make([]byte, 8)

	fmt.Printf(
		"Transfered %d of %d packets. Each packet have %d Kb.\n",
		0,
		totalPackets+1,
		packetSize/1024,
	)

	for {
		_, err := io.ReadFull(conn, bPacketNumber)
		if err != nil {
			fmt.Println("Connection lost. Trying to reconnect")
			conn.Close()
			conn = acceptConn(metadata, sock)
			fmt.Println("Reconnected")
			continue
		}

		packetNumber := binary.LittleEndian.Uint64(bPacketNumber)
		_, err = file.Seek(int64(packetSize*packetNumber), 0)
		if err != nil {
			conn.Close()
			panic("Error moving file offset")
		}

		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			conn.Close()
			panic("Error reading file from disk.")
		}

		n, err = conn.Write(buffer[:n])
		if n == 0 || err != nil {
			fmt.Println("Connection lost. Trying to reconnect")
			conn.Close()
			conn = acceptConn(metadata, sock)
			fmt.Println("Reconnected")
			_, err = file.Seek(-packetSize, 1)
			if err != nil {
				conn.Close()
				panic("Error moving file offset")
			}
			continue
		}

		fmt.Printf(
			"\033[1A\033[2KTransfered %d of %d packets. Each packet have %d Kb.\n",
			packetNumber+1,
			totalPackets+1,
			packetSize/1024,
		)

		if packetNumber == totalPackets {
			break
		}

		clear(buffer)
	}

	conn.Close()
	fmt.Println("File sent successfully!")
}
