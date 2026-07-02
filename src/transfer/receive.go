package transfer

import (
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
)

func newConnector(metadata Metadata) *tls.Conn {
	target := metadata.Host + ":" + metadata.Port
	config := tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp4", target, &config)
	if err != nil {
		panic("Error stabilishing connection with sender at " + target)
	}
	return conn
}

func packetSizeDefReceiver(conn *tls.Conn, metadata Metadata) uint64 {
	pBuffer := make([]byte, 8)
	pSize := uint64(metadata.Psize)
	binary.LittleEndian.PutUint64(pBuffer, pSize)
	conn.Write(pBuffer)
	io.ReadFull(conn, pBuffer)
	pSize = binary.LittleEndian.Uint64(pBuffer)
	pSize = min(pSize, maxPacketSize)
	pSize = max(pSize, minPacketSize)
	return pSize
}

func readNameSize(conn *tls.Conn, pSize uint64) uint64 {
	var fnameSize uint64
	buffer := make([]byte, 8)
	io.ReadFull(conn, buffer)
	fnameSize = binary.LittleEndian.Uint64(buffer)
	if fnameSize > pSize {
		return pSize
	}
	return fnameSize
}

func receivePacket(buffer *[]byte, conn *tls.Conn, packetNumber int) (int, error) {
	bPacketNumber := make([]byte, 8)
	binary.LittleEndian.PutUint64(bPacketNumber, uint64(packetNumber))

	n, err := conn.Write(bPacketNumber)
	if n != 8 || err != nil {
		return 0, errors.New("Error communication packet number")
	}

	n, err = io.ReadFull(conn, *buffer)
	if err != nil && err != io.ErrUnexpectedEOF {
		return 0, err
	}

	return n, nil
}

func readFileMetadata(conn *tls.Conn, metadata Metadata, pSize uint64) FileMetadataHeader {
	fMetadata := FileMetadataHeader{
		FnameSize: 0,
		Fname:     "",
		Fsize:     0,
	}
	if metadata.Fname == "" {
		fnameSize := readNameSize(conn, pSize)
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
	pSize := packetSizeDefReceiver(conn, metadata)
	buffer := make([]byte, pSize)

	fMetadata := readFileMetadata(conn, metadata, pSize)

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
	totalPackets := int(math.Ceil(float64(fMetadata.Fsize) / float64(pSize)))

	fmt.Printf(
		"Transfered %d of %d packets. Each packet have %d Kb.\n",
		0,
		totalPackets,
		pSize/1024,
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

		_, err = file.Write(buffer[:n])
		if err != nil {
			conn.Close()
			panic("Error writing to file")
		}

		clear(buffer)

		fmt.Printf(
			"\033[1A\033[2KTransfered %d of %d packets. Each packet have %d Kb.\n",
			packetNumber+1,
			totalPackets,
			pSize/1024,
		)

		if n != int(pSize) {
			break
		}
	}
	conn.Close()
	fmt.Println("File received successfully and stored at " + fMetadata.Fname)
}
