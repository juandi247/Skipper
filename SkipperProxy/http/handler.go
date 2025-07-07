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
			w.Write([]byte("we are gonna show skipper.lat page"))
		}

		fmt.Println("we start byt getting the tunnel")
		connectedTunnel, err := tm.GetTunnel(subdomain)

		if err != nil {
			fmt.Println("the subdomain doenst exist")
			w.Write([]byte("doesnt exists subdomain!"))
			return
		}

		fmt.Println("now we serilaize")
		responseChannel := make(chan *frame.InternalFrame, 1)
		nextStreamId := DefineStreamId(connectedTunnel, responseChannel)
		finalPayload, payloadLength, err := SerializeHttpRequest(subdomain, r)
		fmt.Println("VAMOS A VER EL LENGTH DEL PAUYLOAD")
	
		fmt.Println(payloadLength)

		if err != nil {
			fmt.Println("error when serialezing the request", err)
			w.Write([]byte("error seralizatin the request"))
			return
		}
		fmt.Println("we serialized the request")

		fmt.Println("stream id defined")
		// defer deleteChannelFromMap(connectedTunnel, nextStreamId)

		requestFrame := frame.CreateFrame(
			1,                          //version
			constants.ProxyRequestType, //requestype
			nextStreamId,               //streamID
			payloadLength)

		buffer := requestFrame.Encode(finalPayload)

	
		_, err = connectedTunnel.Connection.Write(buffer)

		if err != nil {
			fmt.Println("ERRRORRRR sedning the request ", err)
			return
		}
		select {
		case ResponseFrame := <-responseChannel:
			Response, err := DeserializeResponse(ResponseFrame.Payload)
			if err != nil {
				fmt.Println("ERRORRRR despues de desserliar", err)
				return
			}
			tunnel.InternalPayloadPool.Put(ResponseFrame)

		
			fmt.Println("cosoo")

			for key, mivalue:=range Response.GetHeaders(){
				fmt.Println("Para esta key", key)
				for _, value:= range mivalue.GetHeaderValues(){
				w.Header().Set(key,value)
				fmt.Print("valor",value, " ")
				}
				fmt.Println("-----------")
			}
		
			w.WriteHeader(int(Response.GetStatusCode()))
			bodyReader := bytes.NewReader(Response.GetBody())
			io.Copy(w, bodyReader)
			
		}
	}
}
