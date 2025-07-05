package skipperflag

import (
	"SkipperTunnel/constants"
	"flag"
	"fmt"
	"strings"
	"time"
	"unicode"
)

func StartSkipper() (string, int, error) {
	var subdomain string
	var port int
	var flagError string

	printSkipperStart()

	flag.IntVar(&port, "port", 0, "The -p or -port flag is a valid & active port where your app is running. \n e.g -p 8080")
	flag.StringVar(&subdomain, "subdomain", "", "The -s or -subdomain flag is the subdomain that you want to use for your app. \n e.g -p miSubdomain, this will be showed as misubdomain.skipper.lat")
	flag.Parse()
	if port < 1024 || port > 10000 {
		flagError="The port flag is invalid, please use ports on the valid range like 1024 or bigger"
		constants.PrintWithColor(constants.Red, flagError)
		flag.CommandLine.Usage()
		return "", 0, fmt.Errorf(flagError)
	}

	err := ValidateSubdomain(subdomain)
	if err != nil {
		flag.Usage()
		return "", 0, err
	}

	return subdomain, port, nil
}

func ValidateSubdomain(subdomain string) error {
	for _, letterRune := range subdomain {
		if !unicode.IsLetter(letterRune) {
			constants.PrintWithColor(constants.Red, "the subdomain contains invalid characters")
			return fmt.Errorf("the subdomain contains invalid characters")
		}
	}
	if strings.ToLower(subdomain) == "www" {
		constants.PrintWithColor(constants.Red, "invalid subdomain")
		fmt.Println("invalid subdomain")
		return fmt.Errorf("invalid subdomain")
	}
	return nil
}






func printSkipperStart() {

	asciiArt := `
  ___________   .__                            
 /   _____/  | _|__|_____ ______   ___________ 
 \_____  \|  |/ /  \____ \\____ \_/ __ \_  __ \
 /        \    <|  |  |_> >  |_> >  ___/|  | \/
/_______  /__|_ \__|   __/|   __/ \___  >__|   
        \/     \/  |__|   |__|        \/  
		
		
`
	constants.PrintWithColor(constants.Blue, asciiArt)

	time.Sleep(time.Millisecond*1500)
}



