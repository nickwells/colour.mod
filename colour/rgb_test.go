package colour

import "testing"

func TestRGBAllowedValue(t *testing.T) {
	aVal := RGBAllowedValues()
	expectedVal := "a string" +
		" giving the Red/Green/Blue/Alpha values as follows:" +
		" RGB{R: #, G: #, B: #, A: #}" +
		" (defaults: B / G / R: 0x00, A: 0xff)." +
		" Upper and lowercase values are treated equally" +
		" and whitespace is allowed anywhere."

	if aVal != expectedVal {
		t.Log("bad Allowed Value string:")
		t.Log("\t: expected:", expectedVal)
		t.Log("\t:   actual:", aVal)
		t.Error("\t: unexpected RGBAllowedValue")
	}
}
