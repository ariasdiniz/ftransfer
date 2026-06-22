package src

import "net"

func newConnector(metadata Metadata) *net.Conn {
	target := metadata.Host + ":" + metadata.Port
	conn, err := net.Dial("tcp4", target)
	if err != nil {
		panic("Error stabilishing connection with sender at " + target)
	}
	return &conn
}

func Receive(metadata Metadata) {
  
}
