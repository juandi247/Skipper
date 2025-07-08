package tunnel

import (
	"SkipperTunnel/constants"
	skipperflag "SkipperTunnel/flags"
	"SkipperTunnel/forward"
	"SkipperTunnel/proxy"
	"SkipperTunnel/worker"
	"fmt"
	"net"
	"time"
)

type FsmFunc func(*Tunnel) FsmFunc

/*
This is my FSM implementation for the skipper tunnel.
Its inspired by Rob Pike´s code for go lexer, that just opened my mind
on how to write idiomatic fsm without a big switch :)
*/
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
		return nil
	}
	err = proxy.ReadProxyConnResponse(t.ProxyConn)
	if err != nil {
		return nil
	}

	return t.HandleActiveTunnel()
}

func (t *Tunnel) HandleActiveTunnel() FsmFunc {
	go forward.PingLocalhost(t.Ctx, t.LocalhostUrl, t.ErrChan)

	// todo: check number of goroutines
	for i:=0; i<= 35000;i++ {
		go worker.Worker(t.Ctx, t.ProxyConn, t.RequestChan, t.LocalhostUrl)
	}

	go proxy.StartReactor(t.Ctx, t.ProxyConn, t.ErrChan, t.RequestChan, t.syncPool)

	select {
	case <-t.Ctx.Done():
		fmt.Println("the user cancelled everything with control c probalby")
	case err := <-t.ErrChan:
		constants.PrintWithColor(constants.Red, "we received an error from the channel"+err.Error())
	}
	// cleanup the connection
	t.ProxyConn.Close()
	return nil
}
