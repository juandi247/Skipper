package worker

import (
	"SkipperTunnel/constants"
	"SkipperTunnel/forward"
	"SkipperTunnel/frame"
	FramePayloadpb "SkipperTunnel/gen"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
	"google.golang.org/protobuf/proto"
)

// this worker would contain all the logic of the requestss
func Worker(ctx context.Context, proxyConn net.Conn, reactorChan chan *frame.InternalFrame, localhostUrl string) {
	for {
		select {
		case <-ctx.Done():
			return
		case chanRequest := <-reactorChan:
			rawRequest, err := forward.DeserializeRequest(chanRequest.Payload)
			
			if err != nil {
				fmt.Println("error deserializig request", err)
				continue
			}

			bodyReader := bytes.NewReader(rawRequest.GetBody())
			request, err := http.NewRequest(rawRequest.GetMethod(), "http://"+localhostUrl+rawRequest.GetPath(), bodyReader)
			if err != nil {
				fmt.Println("erro creating request", err)
				continue
			}
			headerMap := make(map[string][]string)

		
			for key, value := range rawRequest.GetHeaders() {
				headerMap[key] = value.HeaderValues
			}
			request.Header = headerMap


			httpClient:= &http.Client{
				 Timeout: time.Second*5,
			}
			httpResponse, err:= httpClient.Do(request)
			if err != nil {
				fmt.Println("error making request", err)
				continue
			}

			payload, payloadLength, err:= SerializeHttpResponse(httpResponse)

			if err!=nil{
				fmt.Println("error serlializging request", err)
				continue
			}

			tcpFrame := frame.CreateFrame(1, constants.TunnelResponseType, chanRequest.StreamId, payloadLength)

			buffer:= tcpFrame.Encode(payload)
			_, err = proxyConn.Write(buffer)

			if err!=nil{
				fmt.Println("error", err)
				continue
			}
		}
	}
}





func SerializeHttpResponse(r *http.Response) ([]byte, uint32, error) {
	// headers parsing for seralization
	headersMap := make(map[string]*FramePayloadpb.HeaderValues)
	for key, value := range r.Header {
		headersMap[key] = &FramePayloadpb.HeaderValues{HeaderValues: value}
	}
	responseBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("could not read the body")
	}
	r.Body.Close()

	finalRequest := &FramePayloadpb.Response{
		Status:     r.Status,
		StatusCode: int32(r.StatusCode),
		ProtoMajor: int32(r.ProtoMajor),
		ProtoMinor: int32(r.ProtoMinor),
		Proto:      r.Proto,
		Headers: headersMap,
		Body:  responseBody,
	}

	finalPayload, err := proto.Marshal(finalRequest)
	if err != nil {
		return nil, 0, fmt.Errorf("error marshaling the requst", err)
	}
	return finalPayload, uint32(len(finalPayload)), err
}
