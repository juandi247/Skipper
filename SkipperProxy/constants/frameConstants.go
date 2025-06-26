package constants

import "fmt"

const SkipperMagic = "SKPR"

var SkipperMagicBuffer = [4]byte{SkipperMagic[0], SkipperMagic[1], SkipperMagic[2], SkipperMagic[3]}

// Framing Types
const (
	Control_TunnelRequest = iota + 1
	Control_TunnelAck
	Control_TunnelError
	ProxyRequestType
	TunnelResponseType
	ProxyPing
	TunnelPong
)

func ValidateFramingType(number uint8) error {
	if number <= 0 || number >= 7 {
		return fmt.Errorf("frame type unexistent")
	}
	return nil
}
