package http

import (
	"SkipperProxy/constants"
	"SkipperProxy/frame"
	"SkipperProxy/tunnel"
	"strings"
	"unicode"
)

func ParseSubdomain(host string) (string, bool) {
	host=strings.ToLower(host)

	if host == constants.SkipperUrl || host == constants.WWWSkipperUrl {
		return "", false
	}

	// we want to get out the last .skipper.lat (from the subdomain), we already know that there is a subdomain
	subdomain, exists := strings.CutSuffix(host, "."+constants.SkipperUrl)
	if !exists {
		return "", false
	}

	// now we evaluate if the subdomain contains previously a www. prefix
	subdomain, _ = strings.CutPrefix(subdomain, "www.")

	// this is to avoid double "." or "!"#$%&/()=?", etc on our subdomain call. The browser helps us to do that but still
	for _, letterRune:=range subdomain{
		if !unicode.IsLetter(letterRune) && !unicode.IsNumber(letterRune){
			return "", false
		}
	}
	return subdomain, true
}

func DefineStreamId(tc *tunnel.TunnelConnection,ch chan *frame.InternalFrame ) uint64 {
	tc.Locker.Lock()
	defer tc.Locker.Unlock()
	nextStreamId:=tc.StreamId+1
	tc.StreamMap[nextStreamId] = ch
	tc.StreamId= nextStreamId
	return tc.StreamId
}


func deleteChannelFromMap(tc *tunnel.TunnelConnection, streamId uint64) {
	tc.Locker.Lock()
	defer tc.Locker.Unlock()
	delete(tc.StreamMap, streamId)
}