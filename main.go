package main

import (
	"flag"

	"github.com/ariasdiniz/ftransfer/src/menu"
	"github.com/ariasdiniz/ftransfer/src/transfer"
)

func main() {
	conn := flag.String("conn", "receiver", "Type of connection, you can use \"receiver\" to receive a file or \"sender\" to send")
	fname := flag.String("fname", "", "The path of the file to send/store")
	host := flag.String("host", "0.0.0.0", "The host IPv4 address")
	port := flag.String("port", "8080", "The port to connect to host")
	pSize := flag.Int("psize", 100, "Packet size in Kb")
	flag.Parse()

	meta := transfer.Metadata{
		Conn:  *conn,
		Fname: *fname,
		Host:  *host,
		Port:  *port,
		Psize: *pSize * 1024,
	}

	if *fname == "" || *host == "" {
		meta = menu.ShowMenu(meta)
	}

	switch meta.Conn {
	case "receiver":
		transfer.Receive(meta)
	case "sender":
		transfer.Send(meta)
	default:
		panic("Invalid conn argument. Use receiver or sender")
	}
}
