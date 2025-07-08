package forward

import (
	"context"
	"fmt"
	"net"
	"time"
)

func PingLocalhost(ctx context.Context, url string, errChan chan error) {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			fmt.Println("STOPPED THE GOROUTINE OF PINGING LOCALHOSTSTTT!!")
			return
		case <-ticker.C:
			_, err := net.DialTimeout("tcp", url, time.Second*1)
			if err != nil {
				fmt.Println("something failed on lcoalhost ping", err)
				errChan <- err
				return
			}
			fmt.Println("ping succesfull")
		}
	}
}