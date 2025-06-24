package frame

import (
	"SkipperProxy/constants"
	FramePayloadpb "SkipperProxy/gen"
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net"
)

const MAX_PAYLOAD_LENGTH = 10000000000

type TcpFrame struct {
	Magic      string
	Version    uint8
	FrameType  uint8
	Reserved   [2]byte
	StreamId   uint64
	PayloadLen uint32
}

type Payload []byte

// todo: serialize data payload and save data on a buffer to make a conn.Write
func CreateFrame(version uint8, frameType uint8, streamId uint64, payloadLen uint32) *TcpFrame {
	return &TcpFrame{
		Magic:      constants.SkipperMagic,
		Version:    version,
		FrameType:  frameType,
		StreamId:   streamId,
		PayloadLen: payloadLen,
	}
}

func (f *TcpFrame) Encode(payloadBuffer Payload) []byte {
	FinalBuffer := make([]byte, 20+f.PayloadLen)

	copy(FinalBuffer, []byte(f.Magic))  
	FinalBuffer[4] = f.Version   
	FinalBuffer[5] = f.FrameType 
	//!IMPORTANT here we sart from 8, because the 6,7 position are for padding so they are 0zeroed
	binary.BigEndian.PutUint64(FinalBuffer[8:16], f.StreamId) 
	binary.BigEndian.PutUint32(FinalBuffer[16:20], f.PayloadLen) 

	if f.PayloadLen > 0 {
		copy(FinalBuffer[20:], payloadBuffer)
		return FinalBuffer
	}
	return FinalBuffer
}

func DecodeHeader(buffer []byte) (*TcpFrame, error) {
	if len(buffer) != 20 {
		return nil, fmt.Errorf("the header didnt conatin 20 bytes, this is not suposed to happen?")
	}

	if string(buffer[0:3]) != constants.SkipperMagic {
		return nil, fmt.Errorf("the frame message was incomplete, and deosnt have the MAGIC ")
	}

	f := TcpFrame{}
	f.Magic = string(buffer[0:4])
	f.Version = uint8(buffer[4])
	f.FrameType = uint8(buffer[5])
	f.Reserved = [2]byte(buffer[6:8])
	f.StreamId = binary.BigEndian.Uint64(buffer[8:16])
	f.PayloadLen = binary.BigEndian.Uint32(buffer[16:20])

	return &f, nil
}

func SendAcknowledgeConnection(conn net.Conn) error {
	frame := CreateFrame(1, constants.Control_TunnelAck, 0, 0)
	// there is No payload we just send the type Acknowledge to verify the succesfull connection
	buffer := frame.Encode(nil)
	_, err := conn.Write(buffer)
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

	frame := CreateFrame(1, constants.Control_TunnelError, 0, uint32(len(payloadBytes)))
	buffer := frame.Encode(payloadBytes)
	_, err = conn.Write(buffer)
	if err != nil {
		return fmt.Errorf("the message could not be writeen to the connection %v", err)
	}
	return nil

}
