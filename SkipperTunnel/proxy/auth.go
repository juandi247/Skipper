package proxy

import (
	"SkipperTunnel/constants"
	"SkipperTunnel/frame"
	"fmt"
	"net"
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
	fmt.Println("Waiting to stablish the connection with the proxy...")
	frameType, _, _, payload, err := frame.ReadCompleteFrame(conn)
	if err != nil {
		constants.PrintWithColor(constants.Red, "error connecting with the skipperProxy, try again later")
		return err
	}

	switch frameType {
	case constants.Control_TunnelAck:
		fmt.Println("The connection with the proxy was successfull")
		return nil
	case constants.Control_TunnelError:
		constants.PrintWithColor(constants.Red, "error connecting with the skipperProxy. "+ string(payload))
		return fmt.Errorf("")
	}
	return err
}





