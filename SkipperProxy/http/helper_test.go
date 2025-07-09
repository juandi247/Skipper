package http

import "testing"

func TestSubdomainParsing(t *testing.T) {
	url:="skipper.lat"
	domains := []struct {
		input  string
		output bool
	}{
		{"", false},
		{".", false},
		{"..", false},
		{"mysubdomain.", false},
		{"skipper.lat", false},
		{"#$/%%$#!",false},
		{"skipper"+url, false},
		{"skipper!"+url, false},
		{"skipper!23232."+url, false},
		{"skipper"+url, false},
		{"skippe122323r3232.."+url, false},
		{"$$$$$."+url, false},

	// valid subdomains
		{"mimi."+url, true},
		{"skipper."+url, true},
		{"SKIPPER."+url, true},
		{"skipperS."+url, true},
		{"skipper."+url, true},
		{"12."+url, true},
		{"sk1223ipper."+url, true},
		{"skippe122323r3232."+url, true},
	}

	for _, value := range domains {
		t.Run("test with subdomain "+value.input, func(t *testing.T) {
			_, valid := ParseSubdomain(value.input)
			if valid != value.output {
				t.Errorf("error with the test of %v, should return %v, but returned %v \n \n", value.input, value.output, valid)
			}

		})
	}
}
