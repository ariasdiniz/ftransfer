package transfer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
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

func readNameSize(conn net.Conn) uint64 {
	var fnameSize uint64
	buffer := make([]byte, 8)
	io.ReadFull(conn, buffer)
	fnameSize = binary.LittleEndian.Uint64(buffer)
	if fnameSize > packetSize {
		return packetSize
	}
	return fnameSize
}

func receivePacket(buffer *[]byte, conn net.Conn, packetNumber int) (int, error) {
	bPacketNumber := make([]byte, 8)
	binary.LittleEndian.PutUint64(bPacketNumber, uint64(packetNumber))

	n, err := conn.Write(bPacketNumber)
	if n != 8 || err != nil {
		return 0, errors.New("Error communication packet number")
	}

	n, err = io.ReadFull(conn, *buffer)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func readFileMetadata(conn net.Conn, metadata Metadata) FileMetadataHeader {
	fMetadata := FileMetadataHeader{
		FnameSize: 0,
		Fname:     "",
		Fsize:     0,
	}
	if metadata.Fname == "" {
		fnameSize := readNameSize(conn)
		fname := make([]byte, fnameSize)
		io.ReadFull(conn, fname)
		fMetadata.Fname = filepath.Base(string(fname))
	} else {
		fMetadata.Fname = metadata.Fname
	}
	fMetadata.FnameSize = uint64(len(fMetadata.Fname))
	fSize := make([]byte, 8)
	io.ReadFull(conn, fSize)
	fMetadata.Fsize = binary.LittleEndian.Uint64(fSize)
	return fMetadata
}

func Receive(metadata Metadata) {
	conn := newConnector(metadata)
	buffer := make([]byte, packetSize)

	fMetadata := readFileMetadata(conn, metadata)

	fmt.Println("------------------------------------")
	fmt.Printf("Receiving file: %s\n", fMetadata.Fname)
	fmt.Printf("File size: %d bytes.\n", fMetadata.Fsize)

	file, err := os.Create(fMetadata.Fname)
	if err != nil {
		conn.Close()
		panic("Error creating file " + fMetadata.Fname)
	}
	defer file.Close()

	offset, err := file.Seek(0, 0)
	if err != nil || offset != 0 {
		conn.Close()
		panic("Error opening the file " + fMetadata.Fname)
	}

	fmt.Printf("Starting file transfer, receiving %s\n", fMetadata.Fname)
	totalPackets := int(math.Ceil(float64(fMetadata.Fsize) / packetSize))

	fmt.Printf(
		"Transfered %d of %d packets. Each packet have %d bytes.\n",
		0,
		totalPackets,
		packetSize,
	)

	for packetNumber := range totalPackets {
		n, err := receivePacket(&buffer, conn, packetNumber)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			fmt.Println("Connection lost. Trying to reconnect")
			conn.Close()
			conn = newConnector(metadata)
			fmt.Println("Reconnected")
			continue
		}

		_, err = file.Write(buffer)
		if err != nil {
			conn.Close()
			panic("Error writing to file")
		}

		clear(buffer)

		fmt.Printf(
			"\033[1A\033[2KTransfered %d of %d packets. Each packet have %d bytes.\n",
			packetNumber+1,
			totalPackets,
			packetSize,
		)

		if n != 1024 {
			break
		}
	}
	conn.Close()
	fmt.Println("File received successfully and stored at " + fMetadata.Fname)
}
