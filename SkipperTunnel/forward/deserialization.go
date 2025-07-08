package forward

import (
	FramePayloadpb "SkipperTunnel/gen"
	"fmt"
	"google.golang.org/protobuf/proto"
)

func DeserializeRequest(payload []byte) (*FramePayloadpb.Request, error) {
	frame := &FramePayloadpb.Request{}
	err := proto.Unmarshal(payload, frame)
	if err != nil {
		return nil, fmt.Errorf("eRror unmersshling the request from the proxy", err)
	}

	fmt.Println("vamos a uimripmri ciertas cosas")
	return frame, nil
}
