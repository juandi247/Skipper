package skipperflag

import (
	"testing"
)

func TestSubdomainParsing(t *testing.T) {
	domains := []struct {
		input     string
		outputErr bool
	}{
		{"", true},
		{"    ", true},
		{"    ", true},
		{"  m  ", true},
		{".", true},
		{"..", true},
		{"mysubdomain.", true},
		{"skipper.lat", true},
		{"#$/%%$#!", true},
		{"s,kipper,", true},
		{"skipper!", true},
		{"skipper!23232.", true},
		{"skippe122323r3232..", true},
		{"$$$$$.", true},

		// valid subdomains (only letters and numbers)
		{"mimi", false},
		{"skipper", false},
		{"SKIPPER", false},
		{"skipperS", false},
		{"skipper", false},
		{"12", false},
		{"sk1223ipper", false},
		{"sNkippe122323r3232", false},
	}

	for _, value := range domains {
		t.Run("test with subdomain "+value.input, func(t *testing.T) {
			err := ValidateSubdomain(value.input)
			if (err != nil) != value.outputErr {
				t.Errorf("error with the test of %v, should return %v, but returned %v, %v\n \n", value.input, value.outputErr, (err != nil), err)
			}

		})
	}
}
