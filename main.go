package main

import (
	"flag"
	"github.com/ariasdiniz/ftransfer/src"
)


func main() {
	conn := flag.String("conn", "receiver", "Type of connection, you can use \"receiver\" to receive a file or \"sender\" to send")
	fname := flag.String("fname", "", "The path of the file to send/store")
	host := flag.String("host", "", "The host IPv4 address")
	port := flag.Int("port", 8080, "The port to connect to host. Default 8080")

	flag.Parse()
	
	if *fname == "" || *host == "" {
		panic("You must set conn, fname and host!")
	}

	meta := src.Metadata{
		Conn: *conn,
		Fname: *fname,
		Host: *host,
		Port: *port,
	}

	switch *conn {
  case "receiver":
		src.Receive(meta)
	case "sender":
		src.Send(meta)
	default:
		panic("Invalid conn argument. Use receiver or sender")
	}
}
