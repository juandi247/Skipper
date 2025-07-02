package http

import (
	"SkipperProxy/constants"
	"SkipperProxy/frame"
	"SkipperProxy/tunnel"
	"fmt"
	"net/http"
)

type httpServer struct {
	muxer *http.ServeMux
	port  string
	// RequestTimeout time

}

func CreateHttpServer(port string, tm *tunnel.TunnelManager) *httpServer {
	httpMultiplexer := http.NewServeMux()
	// we register it as wildcard
	httpMultiplexer.HandleFunc("/", ClosureFunc(tm))
	fmt.Println("creamos server")
	return &httpServer{
		muxer: httpMultiplexer,
		port:  port,
	}
}

func (s *httpServer) StartServer() error {
	fmt.Println("intentmos empezarlo")
	err := http.ListenAndServe(s.port, s.muxer)
	if err != nil {
		fmt.Println("ERORRRR")
		return fmt.Errorf(err.Error())
	}
	fmt.Println("empezo el de http")
	return nil
}

func ClosureFunc(tm *tunnel.TunnelManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		subdomain, exists := ParseSubdomain(r.Host)

		if !exists {
			w.Write([]byte("we are gonna show skipper.lat page"))
		}

		tunnel, err := tm.GetTunnel(subdomain)

		if err != nil {
			fmt.Println("the subdomain doenst exist")
			w.Write([]byte("doesnt exists subdomain!"))
			return
		}

		finalPayload, payloadLength, err := SerializeHttpRequest(subdomain, r)

		if err != nil {
			fmt.Println("error when serialezing the request", err)
			w.Write([]byte("error seralizatin the request"))
			return
		}

		nextStreamId := ObtainStreamId(tunnel)

		requestFrame := frame.CreateFrame(
			1,                          //version
			constants.ProxyRequestType, //requestype
			nextStreamId,               //streamID
			payloadLength)

		requestFrame.Encode(finalPayload)

		// buffered channel that waits for the response
		responseChannel := make(chan *frame.InternalFrame, 1)
		SaveResponseChannel(tunnel,nextStreamId,responseChannel)
		
		select {
		case <-responseChannel:
			fmt.Println("cosoo")
			w.Write([]byte("AAAAA RECBIIMOS LE CHHANELL EL EMSNAJE"))

		default:
			fmt.Println("ERRO ACA")
		}
		// make a reader of headers, check if its a cached thing
		// check if its json, xml, the format
		// depending on that make a switch that evaluates that and makes a w.Write to the user

		// the handler is over

		//! Delete the value from the map where the channel and the stream id was
		// ! if not would be a memory leak

	}
}
