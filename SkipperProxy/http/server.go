package http

import (
	"SkipperProxy/constants"
	"SkipperProxy/frame"
	FramePayloadpb "SkipperProxy/gen"
	"SkipperProxy/tunnel"
	"fmt"
	"io"
	"net/http"
	"google.golang.org/protobuf/proto"
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

		fmt.Println(subdomain)

		//todo: context timeout of 10/15 seconds, if there is a timtou show 404

// todo: this should be on the getTunnelFunc frim the tunnel manager
// the mutex should not be touched by this handler func
		tm.Mutex.Lock()
		value, exists := tm.TunnelConnectionsMap[subdomain]
		tm.Mutex.Unlock()
		if !exists {
			w.Write([]byte("404 not found subdomain on the tunnel"))
		}

	
		fmt.Println(value)

// headers parsing for seralization
		headersMap:=make(map[string]*FramePayloadpb.HeaderValues)

		for key, value:=range  r.Header{
			headersMap[key]= &FramePayloadpb.HeaderValues{HeaderValues: value}
		}

		requestBody,err:=io.ReadAll(r.Body)

		if err!=nil{
			w.Write([]byte("error reading the body"))
			return
		}
		defer r.Body.Close()

		finalRequest:= &FramePayloadpb.Request{
			Method:r.Method ,
			Proto: r.Proto,
			TargetUri: subdomain,
			Path: r.RequestURI, 
			Headers: headersMap,
			Body: requestBody,
			RequestId: 1232323,
		}


		requestFrame:= frame.CreateFrame(
			1,  //version
			constants.ProxyRequestType,  //requestype
			122112, //streamID
			uint32(len(requestBody))) //PayloadLength)
 


		finalPayload, err:= proto.Marshal(finalRequest)
		if err!=nil{
			return
		}
		requestFrame.Encode(finalPayload)




		// when we send the request, we should create a channel, that will receive it and will have the streamID
// !aplly the reactor pattern
		// thats a blocking action

		// select that waits the channel or the contextdeadline

		// make a reader of headers, check if its a cached thing
		// check if its json, xml, the format
		// depending on that make a switch that evaluates that and makes a w.Write to the user

		// the handler is over

		fmt.Printf("the path of the request was %v \n", r.RequestURI)
		hello := "hello world"
		w.Write([]byte(hello))
	}
}
