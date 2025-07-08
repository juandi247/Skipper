package http

import (
	"SkipperProxy/constants"
	"SkipperProxy/frame"
	"SkipperProxy/tunnel"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func ClosureFunc(tm tunnel.TunnelManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subdomain, exists := ParseSubdomain(r.Host)

		if !exists {
			http.ServeFile(w, r, "templates/index.html")
			return
		}

		connectedTunnel, err := tm.GetTunnel(subdomain)

		if err != nil {
			http.ServeFile(w, r, "templates/error.html")
			fmt.Println("the subdomain doenst exist")
			return
		}

		responseChannel := make(chan *frame.InternalFrame, 1)
		nextStreamId := DefineStreamId(connectedTunnel, responseChannel)
		defer deleteChannelFromMap(connectedTunnel, nextStreamId)
		finalPayload, payloadLength, err := SerializeHttpRequest(subdomain, r)

		if err != nil {
			http.ServeFile(w, r, "templates/timeout.html")
			return
		}

		requestFrame := frame.CreateFrame(
			1,                          //version
			constants.ProxyRequestType, //requestype
			nextStreamId,               //streamID
			payloadLength)

		buffer := requestFrame.Encode(finalPayload)
		_, err = connectedTunnel.Connection.Write(buffer)

		if err != nil {
			http.ServeFile(w, r, "templates/timeout.html")
			return
		}

		ResponseFrame := <-responseChannel
		Response, err := DeserializeResponse(ResponseFrame.Payload)
		if err != nil {
			http.ServeFile(w, r, "templates/timeout.html")
			fmt.Println("ERRORRRR despues de desserliar", err)
			return
		}
		tunnel.InternalPayloadPool.Put(ResponseFrame)

		for key, mivalue := range Response.GetHeaders() {
			for _, value := range mivalue.GetHeaderValues() {
				w.Header().Set(key, value)
			}
		}

		w.WriteHeader(int(Response.GetStatusCode()))
		bodyReader := bytes.NewReader(Response.GetBody())
		io.Copy(w, bodyReader)

	}
}
