package constants

import "fmt"

const SkipperMagic = "SKPR"

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
