package proxy

import (
	"SkipperTunnel/constants"
	"SkipperTunnel/frame"
	"fmt"
	"net"
	"time"
)

func SendRequestPacket(subdomain string, conn net.Conn) error {
	RequestFrame := frame.CreateFrame(1, constants.Control_TunnelRequest, 0, uint32(len(subdomain)))
	RequestBuffer := RequestFrame.Encode([]byte(subdomain))
	_, err := conn.Write(RequestBuffer)
	if err != nil {
		return err
	}
	return nil
}

func ReadProxyConnResponse(conn net.Conn) error {
	timeError := conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	if timeError != nil {
		fmt.Println("ERROR on dedlinee")
		return timeError
	}

	fmt.Println("we are on the reading loop for acknowledge")
	frameType, _, _, payload, err := frame.ReadCompleteFrame(conn)
	if err != nil {
		fmt.Println("error reading the proxy packet", err)
		return fmt.Errorf("error reading loop or connection closed", err)
	}

	switch frameType {
	case constants.Control_TunnelAck:
		fmt.Println("GOOD WE RECEIVED ALL GOOD, to start the reading of the requestss")
		return nil
	case constants.Control_TunnelError:
		fmt.Println(string(payload))
		return err
	}
	fmt.Println("porque llegamos aca?")
	return err
}





