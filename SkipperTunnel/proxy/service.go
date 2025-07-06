package proxy

import (
	"SkipperTunnel/constants"
	"SkipperTunnel/frame"
	"context"
	"fmt"
	"net"
)

func StartReactor(ctx context.Context, conn net.Conn) error {
	for {
		fmt.Println("we are on the reading loop general")
		frameType, _, _, payload, err := frame.ReadCompleteFrame(conn)
		if err != nil {
			fmt.Println("error reading the proxy packet", err)
			return fmt.Errorf("error reading loop or connection closed", err)
		}

		switch frameType {
		case constants.ProxyRequestType:
			// deserilaize and decode all and redierct to a workerpool

		case constants.ProxyPing:
			// this is just pinging
			fmt.Println(string(payload))
			return err
		}
		fmt.Println("porque llegamos aca?")
		return err
	}
}

