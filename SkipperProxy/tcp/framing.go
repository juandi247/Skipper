package tcp

const MAX_PAYLOAD_LENGTH = 10000000000

type TcpFrame struct {
	Magic      string
	Version    uint8
	FrameType       uint8
	Reserved   [2]byte
	StreamId   uint64
	PayloadLen uint32
}

type Payload []byte

func CreateFrame(version uint8, frameType uint8, streamId uint64, payloadLen uint32) (*TcpFrame){
	return &TcpFrame{
		Magic: "SKPR", //should always have 4 bytes
		Version: version,
		FrameType: frameType,
		PayloadLen: payloadLen,
	}
}