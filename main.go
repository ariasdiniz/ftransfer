package main

import (
	"flag"
)

func main() {
	conn := flag.String("conn", "receiver", "Type of connection, you can use \"receiver\" to receive a file or \"sender\" to send")
	fname := flag.String("fname", "", "The path of the file to send/store")
	target := flag.String("target", "", "The target IPv4 address")

	flag.Parse()
	
	if *fname == "" || *target == "" {
		panic("You must set conn, fname and target!")
	}

	switch *conn {
  case "receiver":
		// TODO
	case "sender":
		// TODO
	default:
		panic("Invalid conn argument. Use receiver or sender")
	}
}
