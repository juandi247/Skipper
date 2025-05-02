package tcpserver

import (
	"fmt"
	"net"
	"time"
	"math/rand"
)

func StartTcp() {
	l, err := net.Listen("tcp4", ":9000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	rand.Seed(time.Now().Unix())

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleTcpConnection(c)
	}
}



