package http

import (
	FramePayloadpb "SkipperProxy/gen"
	"fmt"
	"io"
	"net/http"
	"google.golang.org/protobuf/proto"
)

func SerializeHttpRequest(subdomain string, r *http.Request) ([]byte,uint32, error) {
	// headers parsing for seralization
	headersMap := make(map[string]*FramePayloadpb.HeaderValues)
	for key, value := range r.Header {
		headersMap[key] = &FramePayloadpb.HeaderValues{HeaderValues: value}
	}
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil,0, fmt.Errorf("could not read the body")
	}
	defer r.Body.Close()

	finalRequest := &FramePayloadpb.Request{
		Method:    r.Method,
		Proto:     r.Proto,
		TargetUri: subdomain,
		Path:      r.RequestURI,
		Headers:   headersMap,
		Body:      requestBody,
		RequestId: 1232323,
	}

	finalPayload, err := proto.Marshal(finalRequest)
	if err != nil {
		return nil, 0,fmt.Errorf("error marshaling the requst", err)
	}

	return finalPayload, uint32(len(requestBody)),err
}
