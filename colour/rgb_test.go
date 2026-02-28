package colour

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestRGBAllowedValue(t *testing.T) {
	aVal := RGBAllowedValues()
	expectedVal := "a string" +
		" giving the Red/Green/Blue/Alpha values as follows:" +
		" RGB{R: #, G: #, B: #, A: #}" +
		" (defaults: B / G / R: 0x00, A: 0xff)." +
		" Upper and lowercase values are treated equally" +
		" and whitespace is allowed anywhere." +
		"\n\n" +
		`Or a literal hash ("#")` +
		" immediately followed by precisely 3 or 6 hexadecimal digits"

	if aVal != expectedVal {
		t.Log("bad Allowed Value string:")
		t.Log("\t: expected:", expectedVal)
		t.Log("\t:   actual:", aVal)
		t.Error("\t: unexpected RGBAllowedValue")
	}
}

func TestParseNDigitColour(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		s         string
		re        *regexp.Regexp
		mkDigits  func(string) string
		expColour rgba
	}{
		{
			ID: testhelper.MkID("bad - regexp returns too many parts"),
			ExpErr: testhelper.MkExpErr(
				`the test colour ("#ffffff") is badly formed:` +
					` 3 parts expected, 6 found`),
			s: "#ffffff",
			re: regexp.MustCompile(
				`#([[:xdigit:]])` +
					`([[:xdigit:]])` +
					`([[:xdigit:]])` +
					`([[:xdigit:]])` +
					`([[:xdigit:]])` +
					`([[:xdigit:]])`),
			mkDigits: func(xd string) string { return xd },
		},
		{
			ID: testhelper.MkID("bad - regexp returns too few parts"),
			ExpErr: testhelper.MkExpErr(
				`the test colour ("#ffffff") is badly formed:` +
					` 3 parts expected, 2 found`),
			s: "#ffffff",
			re: regexp.MustCompile(
				`#([[:xdigit:]][[:xdigit:]][[:xdigit:]])` +
					`([[:xdigit:]][[:xdigit:]][[:xdigit:]])`),
			mkDigits: func(xd string) string { return xd },
		},
		{
			ID: testhelper.MkID("bad - regexp returns bad parts"),
			ExpErr: testhelper.MkExpErr(
				`the test colour ("#ffffff") is badly formed:` +
					` digit 1(fff) cannot be converted to a number`),
			s: "#ffffff",
			re: regexp.MustCompile(
				`#([[:xdigit:]][[:xdigit:]][[:xdigit:]])` +
					`([[:xdigit:]])` +
					`([[:xdigit:]][[:xdigit:]])`),
			mkDigits: func(xd string) string { return xd },
		},
		{
			ID:        testhelper.MkID("3-digits, white"),
			s:         "#fff",
			re:        rgbAlt3RE,
			mkDigits:  func(xd string) string { return xd + xd },
			expColour: rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		},
		{
			ID:        testhelper.MkID("6-digits, white"),
			s:         "#ffffff",
			re:        rgbAlt6RE,
			mkDigits:  func(xd string) string { return xd },
			expColour: rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		},
		{
			ID:        testhelper.MkID("3-digits, black"),
			s:         "#000",
			re:        rgbAlt3RE,
			mkDigits:  func(xd string) string { return xd + xd },
			expColour: rgba{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		},
		{
			ID:        testhelper.MkID("6-digits, black"),
			s:         "#000000",
			re:        rgbAlt6RE,
			mkDigits:  func(xd string) string { return xd },
			expColour: rgba{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		},
		{
			ID:        testhelper.MkID("3-digits, red"),
			s:         "#f00",
			re:        rgbAlt3RE,
			mkDigits:  func(xd string) string { return xd + xd },
			expColour: rgba{R: 0xff, G: 0x00, B: 0x00, A: 0xff},
		},
		{
			ID:        testhelper.MkID("6-digits, red"),
			s:         "#ff0000",
			re:        rgbAlt6RE,
			mkDigits:  func(xd string) string { return xd },
			expColour: rgba{R: 0xff, G: 0x00, B: 0x00, A: 0xff},
		},
		{
			ID:        testhelper.MkID("3-digits, green"),
			s:         "#0f0",
			re:        rgbAlt3RE,
			mkDigits:  func(xd string) string { return xd + xd },
			expColour: rgba{R: 0x00, G: 0xff, B: 0x00, A: 0xff},
		},
		{
			ID:        testhelper.MkID("6-digits, green"),
			s:         "#00ff00",
			re:        rgbAlt6RE,
			mkDigits:  func(xd string) string { return xd },
			expColour: rgba{R: 0x00, G: 0xff, B: 0x00, A: 0xff},
		},
		{
			ID:        testhelper.MkID("3-digits, blue"),
			s:         "#00f",
			re:        rgbAlt3RE,
			mkDigits:  func(xd string) string { return xd + xd },
			expColour: rgba{R: 0x00, G: 0x00, B: 0xff, A: 0xff},
		},
		{
			ID:        testhelper.MkID("6-digits, blue"),
			s:         "#0000ff",
			re:        rgbAlt6RE,
			mkDigits:  func(xd string) string { return xd },
			expColour: rgba{R: 0x00, G: 0x00, B: 0xff, A: 0xff},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actColour, err := parseNDigitColour(tc.s, "test", tc.re, tc.mkDigits)
			testhelper.CheckExpErrWithID(t, tc.IDStr(), err, tc)

			if err == nil {
				dvErr := testhelper.DiffVals(actColour, tc.expColour)
				if dvErr != nil {
					t.Log(tc.IDStr())
					t.Logf("\t: expected colour: %#v", tc.expColour)
					t.Logf("\t:   actual colour: %#v", actColour)
					t.Errorf("\t: colours differ: %s", dvErr)
				}
			}
		})
	}
}

func TestParse3DigitColour(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		s         string
		expColour rgba
	}{
		{
			ID: testhelper.MkID("bad - empty string"),
			ExpErr: testhelper.MkExpErr(
				`the 3-digit colour ("") is badly formed`),
			s: "",
		},
		{
			ID: testhelper.MkID("bad - no hash"),
			ExpErr: testhelper.MkExpErr(
				`the 3-digit colour ("fff") is badly formed`),
			s: "fff",
		},
		{
			ID: testhelper.MkID("bad - too few digits"),
			ExpErr: testhelper.MkExpErr(
				`the 3-digit colour ("#ff") is badly formed`),
			s: "#ff",
		},
		{
			ID: testhelper.MkID("bad - too many digits"),
			ExpErr: testhelper.MkExpErr(
				`the 3-digit colour ("#ffff") is badly formed`),
			s: "#ffff",
		},
		{
			ID: testhelper.MkID("bad - non-hex digits"),
			ExpErr: testhelper.MkExpErr(
				`the 3-digit colour ("#ffg") is badly formed`),
			s: "#ffg",
		},
		{
			ID:        testhelper.MkID("white"),
			s:         "#fff",
			expColour: rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		},
		{
			ID:        testhelper.MkID("black"),
			s:         "#000",
			expColour: rgba{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		},
		{
			ID:        testhelper.MkID("red"),
			s:         "#f00",
			expColour: rgba{R: 0xff, G: 0x00, B: 0x00, A: 0xff},
		},
		{
			ID:        testhelper.MkID("green"),
			s:         "#0f0",
			expColour: rgba{R: 0x00, G: 0xff, B: 0x00, A: 0xff},
		},
		{
			ID:        testhelper.MkID("blue"),
			s:         "#00f",
			expColour: rgba{R: 0x00, G: 0x00, B: 0xff, A: 0xff},
		},
		{
			ID:        testhelper.MkID("123"),
			s:         "#123",
			expColour: rgba{R: 0x11, G: 0x22, B: 0x33, A: 0xff},
		},
		{
			ID:        testhelper.MkID("123 - with leading space"),
			s:         "   #123",
			expColour: rgba{R: 0x11, G: 0x22, B: 0x33, A: 0xff},
		},
		{
			ID:        testhelper.MkID("123 - with trailing space"),
			s:         "#123   ",
			expColour: rgba{R: 0x11, G: 0x22, B: 0x33, A: 0xff},
		},
		{
			ID:        testhelper.MkID("123 - with leading and trailing space"),
			s:         "   #123   ",
			expColour: rgba{R: 0x11, G: 0x22, B: 0x33, A: 0xff},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actColour, err := Parse3DigitColour(tc.s)
			testhelper.CheckExpErrWithID(t, tc.IDStr(), err, tc)

			if err == nil {
				dvErr := testhelper.DiffVals(actColour, tc.expColour)
				if dvErr != nil {
					t.Log(tc.IDStr())
					t.Logf("\t: expected colour: %#v", tc.expColour)
					t.Logf("\t:   actual colour: %#v", actColour)
					t.Errorf("\t: colours differ: %s", dvErr)
				}
			}
		})
	}
}

func TestParse6DigitColour(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		s         string
		expColour rgba
	}{
		{
			ID: testhelper.MkID("bad - empty string"),
			ExpErr: testhelper.MkExpErr(
				`the 6-digit colour ("") is badly formed`),
			s: "",
		},
		{
			ID: testhelper.MkID("bad - no hash"),
			ExpErr: testhelper.MkExpErr(
				`the 6-digit colour ("ffffff") is badly formed`),
			s: "ffffff",
		},
		{
			ID: testhelper.MkID("bad - too few digits"),
			ExpErr: testhelper.MkExpErr(
				`the 6-digit colour ("#fffff") is badly formed`),
			s: "#fffff",
		},
		{
			ID: testhelper.MkID("bad - too many digits"),
			ExpErr: testhelper.MkExpErr(
				`the 6-digit colour ("#fffffff") is badly formed`),
			s: "#fffffff",
		},
		{
			ID: testhelper.MkID("bad - non-hex digits"),
			ExpErr: testhelper.MkExpErr(
				`the 6-digit colour ("#fffffg") is badly formed`),
			s: "#fffffg",
		},
		{
			ID:        testhelper.MkID("white"),
			s:         "#ffffff",
			expColour: rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		},
		{
			ID:        testhelper.MkID("black"),
			s:         "#000000",
			expColour: rgba{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
		},
		{
			ID:        testhelper.MkID("red"),
			s:         "#ff0000",
			expColour: rgba{R: 0xff, G: 0x00, B: 0x00, A: 0xff},
		},
		{
			ID:        testhelper.MkID("green"),
			s:         "#00ff00",
			expColour: rgba{R: 0x00, G: 0xff, B: 0x00, A: 0xff},
		},
		{
			ID:        testhelper.MkID("blue"),
			s:         "#0000ff",
			expColour: rgba{R: 0x00, G: 0x00, B: 0xff, A: 0xff},
		},
		{
			ID:        testhelper.MkID("123456"),
			s:         "#123456",
			expColour: rgba{R: 0x12, G: 0x34, B: 0x56, A: 0xff},
		},
		{
			ID:        testhelper.MkID("123456 - with leading space"),
			s:         "   #123456",
			expColour: rgba{R: 0x12, G: 0x34, B: 0x56, A: 0xff},
		},
		{
			ID:        testhelper.MkID("123456 - with trailing space"),
			s:         "#123456   ",
			expColour: rgba{R: 0x12, G: 0x34, B: 0x56, A: 0xff},
		},
		{
			ID: testhelper.MkID(
				"123456 - with leading and trailing space"),
			s:         "   #123456   ",
			expColour: rgba{R: 0x12, G: 0x34, B: 0x56, A: 0xff},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actColour, err := Parse6DigitColour(tc.s)
			testhelper.CheckExpErrWithID(t, tc.IDStr(), err, tc)

			if err == nil {
				dvErr := testhelper.DiffVals(actColour, tc.expColour)
				if dvErr != nil {
					t.Log(tc.IDStr())
					t.Logf("\t: expected colour: %#v", tc.expColour)
					t.Logf("\t:   actual colour: %#v", actColour)
					t.Errorf("\t: colours differ: %s", dvErr)
				}
			}
		})
	}
}

func TestParseColourDefinition(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		s         string
		expColour rgba
	}{
		{
			ID: testhelper.MkID("bad - empty string"),
			ExpErr: testhelper.MkExpErr(
				`the colour definition ("") is invalid`),
		},
		{
			ID:        testhelper.MkID("3-digit - 123"),
			s:         "#123",
			expColour: rgba{R: 0x11, G: 0x22, B: 0x33, A: 0xff},
		},
		{
			ID:        testhelper.MkID("6-digit - 123456"),
			s:         "#123456",
			expColour: rgba{R: 0x12, G: 0x34, B: 0x56, A: 0xff},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actColour, err := ParseColourDefinition(tc.s)
			testhelper.CheckExpErr(t, err, tc)

			if err == nil {
				dvErr := testhelper.DiffVals(actColour, tc.expColour)
				if dvErr != nil {
					t.Log(tc.IDStr())
					t.Logf("\t: expected colour: %#v", tc.expColour)
					t.Logf("\t:   actual colour: %#v", actColour)
					t.Errorf("\t: colours differ: %s", dvErr)
				}
			}
		})
	}
}

func TestIsAPotentialColourString(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		s    string
		expB bool
	}{
		{
			ID:   testhelper.MkID("empty string"),
			s:    "",
			expB: false,
		},
		{
			ID:   testhelper.MkID("random non-colour string"),
			s:    "blah",
			expB: false,
		},
		{
			ID:   testhelper.MkID("starts with a hash but has no digits"),
			s:    "#",
			expB: false,
		},
		{
			ID:   testhelper.MkID("starts with a hash but has only 2 digits"),
			s:    "#12",
			expB: false,
		},
		{
			ID:   testhelper.MkID("starts with a hash but has >3 and <6 digits"),
			s:    "#1234",
			expB: false,
		},
		{
			ID:   testhelper.MkID("starts with a hash but has >6 digits"),
			s:    "#1234567",
			expB: false,
		},
		{
			ID:   testhelper.MkID("starts with a hash but has bad digits"),
			s:    "#12x",
			expB: false,
		},
		{
			ID:   testhelper.MkID("starts with a hash but has bad digits"),
			s:    "#12x",
			expB: false,
		},
		{
			ID:   testhelper.MkID("starts with a hash, 3 chars but bad digits"),
			s:    "#12x",
			expB: false,
		},
		{
			ID:   testhelper.MkID("starts with a hash, 6 chars but bad digits"),
			s:    "#12x456",
			expB: false,
		},
		{
			ID:   testhelper.MkID("good - hash, 3 digits"),
			s:    "#126",
			expB: true,
		},
		{
			ID:   testhelper.MkID("good - hash, 6 digits"),
			s:    "#123456",
			expB: true,
		},
		{
			ID:   testhelper.MkID("good - rgb"),
			s:    "rgb{r:42}",
			expB: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actB := IsAPotentialColourString(tc.s)
			testhelper.DiffBool(t,
				tc.IDStr(), fmt.Sprintf("result for %q", tc.s),
				actB, tc.expB)
		})
	}
}
