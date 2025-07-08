package forward

import (
	FramePayloadpb "SkipperTunnel/gen"
	"google.golang.org/protobuf/proto"
)

func DeserializeRequest(payload []byte) (*FramePayloadpb.Request, error) {
	frame := &FramePayloadpb.Request{}
	err := proto.Unmarshal(payload, frame)
	if err != nil {
		return nil, err
	}
	return frame, nil
}
