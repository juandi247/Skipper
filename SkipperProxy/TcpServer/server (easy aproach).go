package tcpserver

import (
	"fmt"
	"net"
)

func StartTcp() {
	l, err := net.Listen("tcp4", ":9000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleTcpConnection(c)
	}
}
