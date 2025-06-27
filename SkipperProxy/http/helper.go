package http

import (
	"SkipperProxy/constants"
	"strings"
)

func ParseSubdomain(host string) (string, bool) {

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

	// this is kust to evaluate that the subdomain is not like subdomain.anothersubdmain.another.skipper.lat
	// on our bussines logic for now we only allow a single subdomain
	containsADot := strings.Contains(subdomain, ".")
	if containsADot {
		return "", false
	}
	return subdomain, true
}
