package frame

import (
	"SkipperProxy/constants"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"slices"
	"time"
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

type InternalFrame struct{
	StreamId uint64
	Payload
}

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

	copy(FinalBuffer, constants.SkipperMagicBuffer[:])
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

func DecodeHeader(buffer []byte) (uint8, uint64, uint32, error) {
	if len(buffer) != 20 {
		return 0, 0, 0, fmt.Errorf("the header didnt conatin 20 bytes, this is not suposed to happen?")
	}

	// validate magicNumber
	if !slices.Equal(buffer[0:4], constants.SkipperMagicBuffer[:]) {
		return 0, 0, 0, fmt.Errorf("the frame message was incomplete, and deosnt have the MAGIC ")
	}

	// validate the typ
	err := constants.ValidateFramingType(buffer[5])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("didnt receive right after the connection the REquest Type good")
	}

	FrameType := uint8(buffer[5])
	StreamId := binary.BigEndian.Uint64(buffer[8:16])
	PayloadLen := binary.BigEndian.Uint32(buffer[16:20])

	return FrameType, StreamId, PayloadLen, nil
}

/*
1.Read the 20 bytes from the buffer (a valid frame)
2.Decode the frame and return the important values (in this decode we validate the magic number, validate framing type etc)
3.create a payloadbuffer if the payload has anything and return it, if not return the payload null
*/
func ReadCompleteFrame(conn net.Conn, isControlType bool) (uint8, uint64, uint32, Payload, error) {

	// this returns a time
	deadlineTime := time.Now()
	/* here the setREadDEadline only recevies a strcut deadline (no pointer)
	   but we need to check the dureation of ddline, so we need to add that
	    .add so it doesnt get exectuted right on
	*/
	if err := conn.SetReadDeadline(deadlineTime.Add(2 * time.Second)); err != nil {
		fmt.Println("Error en readline deadlineeee mimiim")
		return 0, 0, 0, nil, fmt.Errorf(err.Error(), "error reading because of timouttt")
	}

	fmt.Println("we are reding the frame")
	buffer := make([]byte, 20)
	_, err := io.ReadFull(conn, buffer)
	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("ERROR Reading the buffer")
	}
	frameType, streamId, Payloadlength, err := DecodeHeader(buffer)

	if err != nil {
		return 0, 0, 0, nil, fmt.Errorf("could not parse the header", err)

	}

	// we create the complete buffer if the payload has content (length is more than cero and we return the payload here too)
	if Payloadlength > 0 {
		payloadBuffer := make([]byte, Payloadlength)
		_, err = io.ReadFull(conn, payloadBuffer)
		if err != nil {
			return 0, 0, 0, nil, fmt.Errorf("ERROR, the payload length was sent incomplete")
		}
		return frameType, streamId, Payloadlength, payloadBuffer, nil
	}

	return frameType, streamId, Payloadlength, nil, nil
}
