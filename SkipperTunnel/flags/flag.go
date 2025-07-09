package skipperflag

import (
	"SkipperTunnel/constants"
	"flag"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

/*This function returns 
string -> the valid subdomain
sstring-> the parsed localhost url with the port: e.g localhost:8080
an error */ 
func FlagValidation() (string, string, error) {

	var port int
	var subdomain string
	printSkipperStart()

	flag.IntVar(&port , "port", 0, "The -p or -port flag is a valid & active port where your app is running. \n e.g -p 8080")
	flag.StringVar(&subdomain, "subdomain", "", "The -s or -subdomain flag is the subdomain that you want to use for your app. \n e.g -p miSubdomain, this will be showed as misubdomain.skipper.lat")
	flag.Parse()
	if port < 1024 || port > 100000 {
		constants.PrintWithColor(constants.Red,"The port flag is invalid, please use ports on the valid range like 1024 or bigger")
		flag.CommandLine.Usage()
		return "", "", fmt.Errorf("invaild localhostPort")
	}

	localhostUrl:= "localhost:"+strconv.Itoa(port)

	subdomain= strings.ToLower(subdomain)
	subdomain = strings.TrimSpace(subdomain)

	err := ValidateSubdomain(subdomain)
	if err != nil {
		constants.PrintWithColor(constants.Red, err.Error())
		flag.Usage()
		return "", "",err
	}
	return subdomain, localhostUrl, nil
}



func ValidateSubdomain(subdomain string) error {
	if subdomain == "www" || subdomain==""{
		return fmt.Errorf("invalid subdomain")
	}

	for _, letterRune := range subdomain {
		if !unicode.IsLetter(letterRune) && !unicode.IsNumber(letterRune) {
			return fmt.Errorf("the subdomain contains invalid characters")
		}
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



