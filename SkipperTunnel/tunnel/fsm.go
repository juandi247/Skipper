package tunnel

import (
	"SkipperTunnel/constants"
	skipperflag "SkipperTunnel/flags"
	"SkipperTunnel/forward"
	"SkipperTunnel/proxy"
	"fmt"
	"net"
	"time"
)

type FsmFunc func(*Tunnel) FsmFunc

func (t *Tunnel) FsmStart() {
	var state FsmFunc = t.HandleInitialization()
	for state != nil {
		state = state(t)
	}
}

// States
func (t *Tunnel) HandleInitialization() FsmFunc {
	subdomain, localhostUrl, err := skipperflag.FlagValidation()

	if err != nil {
		fmt.Println("error: ", err)
		return nil
	}
	t.LocalhostUrl = localhostUrl
	t.Subdomain = subdomain
	return t.HandleLocalhostConnection()
}

func (t *Tunnel) HandleLocalhostConnection() FsmFunc {

	_, err := net.DialTimeout("tcp", t.LocalhostUrl, time.Second*1)

	if err != nil {
		constants.PrintWithColor(constants.Red, "failed to connect to your localhost app. Check that the port is correct")
		return nil
	}
	// todo: add to this goroutine the switch and power to detemirne the cancell of all things
	go forward.PingLocalhost(t.Ctx, t.LocalhostUrl)

	return t.HandleProxyConnection()
}

func (t *Tunnel) HandleProxyConnection() FsmFunc {

	proxyConn, err := net.Dial("tcp", t.ProxyUrl)
	if err != nil {
		fmt.Println("error ", err)
		return nil
	}
	t.ProxyConn = proxyConn
	err = proxy.SendRequestPacket(t.Subdomain, t.ProxyConn)
	if err != nil {
		// todo: send error to handler
		return nil
	}
	err = proxy.ReadProxyConnResponse(t.ProxyConn)
	if err != nil {
		// todo: send error to handler (because we need to stop the other goroutine of pinging localhost)
		return nil
	}

	return t.HandleActiveTunnel()
}

func (t *Tunnel) HandleActiveTunnel() FsmFunc {
	proxy.StartReactor(t.Ctx, t.ProxyConn)
	
	return nil
}
