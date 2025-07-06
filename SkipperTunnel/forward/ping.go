package forward

import (
	"context"
	"fmt"
	"net"
	"time"
)

func PingLocalhost(ctx context.Context, url string) {
	for {
		_, err := net.DialTimeout("tcp", url, time.Second*1)
		if err!=nil{
			fmt.Println("something failed on lcoalhost ping", err)
			// todo: add the send to the error channel handler, meaning that we need to cancell all the opened goroutines
			return
		}
		fmt.Println("ping succesfull")
		time.Sleep(time.Second*2)
	}
}