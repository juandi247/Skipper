package tunnel

import (
	"SkipperProxy/constants"
	"SkipperProxy/frame"
	FramePayloadpb "SkipperProxy/gen"
	"fmt"
	"net"
	"sync"

	"google.golang.org/protobuf/proto"
)

type TunnelConnection struct {
	Subdomain   string
	Connection  net.Conn
	ip          net.Addr
	StreamId    uint64
	StreamMap map[uint64]chan *frame.InternalFrame
	Locker sync.Mutex
}

var InternalPayloadPool = &sync.Pool{
	New: func() interface{} {
		return &frame.InternalFrame{}
	},
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
		frameType, streamId, _, payload, err := frame.ReadCompleteFrame(tc.Connection, false)
		if err != nil {
			fmt.Println("error on reading loop", err)
			return
		}

		switch frameType {
		case constants.TunnelResponseType:

			item := InternalPayloadPool.Get()
			// type assert, in this case item is an empty interface{} so we need to use this type to give it a value typed
			// the go compiler doesnt know that there is an especiied type here.
			frame := item.(*frame.InternalFrame)
			frame.StreamId = streamId
			frame.Payload = payload

			tc.Locker.Lock()
			ResponseChannel, exists := tc.StreamMap[streamId]
			tc.Locker.Unlock()
			if !exists {
				continue
			}
			// non blocking action becasue its a buffered channel
			ResponseChannel <- frame
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
