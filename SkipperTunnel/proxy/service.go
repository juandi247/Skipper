package proxy

import (
	"SkipperTunnel/constants"
	"SkipperTunnel/frame"
	"context"
	"fmt"
	"net"
	"sync"
)

func StartReactor(ctx context.Context, conn net.Conn, errChan chan error, requestChan chan *frame.InternalFrame, sp *sync.Pool) {
	for {
		frameType, streamId, _, payload, err := frame.ReadCompleteFrame(conn)
		if err != nil {
			select {
			case <-errChan:
				return
			default:
				errChan<-err
				// fmt.Println("error reading the proxy packet", err)
				return
			}
		}

		switch frameType {
		case constants.ProxyRequestType:
			item := sp.Get()
			// type assert, in this case item is an empty interface{} so we need to use this type to give it a value typed
			// the go compiler doesnt know that there is an especiied type here.
			frame := item.(*frame.InternalFrame)
			frame.StreamId = streamId
			frame.Payload = payload
			// fmt.Println("Received a Request")
			// non blocking action becasue its a buffered channel
			requestChan <- frame
		case constants.TunnelPong:
			fmt.Println("tunnel is kept oppened")
		}
	}
}