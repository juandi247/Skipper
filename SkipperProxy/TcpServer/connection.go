package tcpserver

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func handleTcpConnection(c net.Conn){

	defer c.Close()

	for{
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil{
			fmt.Print(err)
			return
		}

		temp:= strings.TrimSpace(string(netData))
		if temp== "STOP"{
			break
		}

		result:= strconv.Itoa(random())+"\n"
		c.Write([]byte(string(result)))

	}


}
func random() int {
	return 12345 
}

