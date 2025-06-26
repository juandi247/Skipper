package tunnel

import (
	"SkipperProxy/constants"
	"SkipperProxy/frame"
	FramePayloadpb "SkipperProxy/gen"
	"fmt"
	"net"

	"google.golang.org/protobuf/proto"
)

type TunnelConnection struct {
	Subdomain  string
	Connection net.Conn
	ip         net.Addr

	// todo:!!!!!
	// multipleixing fields
	// nextstream
	// pdenign things, (idk yet)
}

func CreateTunnelConnection(subdomain string, conn net.Conn, ip net.Addr) *TunnelConnection {
	return &TunnelConnection{
		Subdomain:  subdomain,
		Connection: conn,
		ip:         ip,
	}
}

func (tc *TunnelConnection) StartReadLoop() {
	for {
		frameType, streamId, payloadLen, payload, err := frame.ReadCompleteFrame(tc.Connection, false)
		if err != nil {
			fmt.Println("error on reading loop", err)
			return
		}

		switch frameType {
		case constants.TunnelResponseType:
			fmt.Println("llego la respuesta", payload)
			fmt.Println(streamId, payloadLen)
		case constants.TunnelPong:
			fmt.Println("tunnel is kept oppened")
		}

	}
}

func (tc *TunnelConnection) SendAcknowledgeConnection() error {
	frame := frame.CreateFrame(1, constants.Control_TunnelAck, 0, 0)
	// there is No payload we just send the type Acknowledge to verify the succesfull connection
	buffer := frame.Encode(nil)
	_, err := tc.Connection.Write(buffer)
	if err != nil {
		return fmt.Errorf("the message could not be writeen to the connection %d", err)
	}
	return nil
}

func SendControlConnectionError(conn net.Conn, controlError string) error {
	payload := &FramePayloadpb.Control_Tunnel_Error{
		Error: controlError,
	}
	payloadBytes, err := proto.Marshal(payload)

	if err != nil {
		fmt.Println("error!! ", err)
		panic("ERROR marshaling the payload")
	}
	frame := frame.CreateFrame(1, constants.Control_TunnelError, 0, uint32(len(payloadBytes)))
	buffer := frame.Encode(payloadBytes)
	_, err = conn.Write(buffer)
	if err != nil {
		return fmt.Errorf("the message could not be writeen to the connection %v", err)
	}
	return nil

}
