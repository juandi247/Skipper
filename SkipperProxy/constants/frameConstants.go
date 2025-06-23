package constants

const SkipperMagic = "SKPR"

// Framing Types
const (
	Control_TunnelRequest = iota + 1
	Control_TunnelAck
	Control_TunnelError
	RequestType
	ResponseType
	ErrorType
)
