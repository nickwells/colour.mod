package colour

import (
	"image/color" //nolint:misspell
	"testing"

	"github.com/nickwells/colour.mod/colourtesthelper"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestColourFamily(t *testing.T) {
	if len(searchOrder) != len(cFamMap)-2 {
		t.Errorf("There are %d entries in the  search order list"+
			" & %d in the families map (with two families excluded)",
			len(searchOrder), len(cFamMap))
	}
}

func TestFamilyMethods(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic

		f          Family
		expStr     string
		expLiteral string
		isValid    bool
	}{
		{
			ID:         testhelper.MkID("good - Any"),
			f:          AnyColours,
			expStr:     "Any",
			expLiteral: "AnyColours",
			isValid:    true,
		},
		{
			ID:         testhelper.MkID("good - Web"),
			f:          WebColours,
			expStr:     "Web",
			expLiteral: "WebColours",
			isValid:    true,
		},
		{
			ID:         testhelper.MkID("good - CGA"),
			f:          CGAColours,
			expStr:     "CGA",
			expLiteral: "CGAColours",
			isValid:    true,
		},
		{
			ID:         testhelper.MkID("good - X11"),
			f:          X11Colours,
			expStr:     "X11",
			expLiteral: "X11Colours",
			isValid:    true,
		},
		{
			ID:         testhelper.MkID("good - HTML"),
			f:          HTMLColours,
			expStr:     "HTML",
			expLiteral: "HTMLColours",
			isValid:    true,
		},
		{
			ID:         testhelper.MkID("good - Pantone"),
			f:          PantoneColours,
			expStr:     "Pantone",
			expLiteral: "PantoneColours",
			isValid:    true,
		},
		{
			ID:       testhelper.MkID("bad - 99"),
			ExpPanic: testhelper.MkExpPanic("BadFamily:99"),
			f:        Family(99),
			expStr:   "BadFamily:99",
			isValid:  false,
		},
	}

	for _, tc := range testCases {
		str := tc.f.String()
		testhelper.DiffString(t, tc.IDStr(), "String", str, tc.expStr)

		var literal string

		panicked, panicVal := testhelper.PanicSafe(func() {
			literal = tc.f.Literal()
		})
		if !panicked {
			testhelper.DiffString(t, tc.IDStr(), "Literal",
				literal, tc.expLiteral)
		}

		testhelper.CheckExpPanicError(t, panicked, panicVal, tc)
		isValid := tc.f.IsValid()
		testhelper.DiffBool(t, tc.IDStr(), "Valid check", isValid, tc.isValid)
	}
}

func TestFamilyList(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		fl     []Family
		expStr string
	}{
		{
			ID: testhelper.MkID("none"),
		},
		{
			ID:     testhelper.MkID("just one"),
			fl:     []Family{WebColours},
			expStr: "Web",
		},
		{
			ID:     testhelper.MkID("two"),
			fl:     []Family{WebColours, X11Colours},
			expStr: "Web and X11",
		},
		{
			ID:     testhelper.MkID("three"),
			fl:     []Family{WebColours, X11Colours, PantoneColours},
			expStr: "Web, X11 and Pantone",
		},
	}

	for _, tc := range testCases {
		str := familyList(tc.fl)
		testhelper.DiffString(t, tc.IDStr(), "family list",
			str, tc.expStr)
	}
}

func TestFamilyNames(t *testing.T) {
	var (
		webCount     = len(webColours)
		cgaCount     = len(cgaColours)
		x11Count     = len(x11Colours)
		htmlCount    = len(htmlColours)
		pantoneCount = len(pantoneColours)
		anyCount     = webCount +
			cgaCount +
			x11Count +
			htmlCount +
			pantoneCount
	)

	testCases := []struct {
		testhelper.ID
		f            Family
		expNameCount int
	}{
		{
			ID:           testhelper.MkID("Web"),
			f:            WebColours,
			expNameCount: webCount,
		},
		{
			ID:           testhelper.MkID("CGA"),
			f:            CGAColours,
			expNameCount: cgaCount,
		},
		{
			ID:           testhelper.MkID("X11"),
			f:            X11Colours,
			expNameCount: x11Count,
		},
		{
			ID:           testhelper.MkID("HTML"),
			f:            HTMLColours,
			expNameCount: htmlCount,
		},
		{
			ID:           testhelper.MkID("Pantone"),
			f:            PantoneColours,
			expNameCount: pantoneCount,
		},
	}

	for _, tc := range testCases {
		names := tc.f.ColourNames()
		testhelper.DiffInt(t, tc.IDStr(), "number of colour names",
			len(names), tc.expNameCount)
	}

	// Handle the 'Any' family separately - we remove duplicates so we can
	// only say that the number of names must be less than the sum of the
	// constituent name sets
	names := AnyColours.ColourNames()
	if len(names) > anyCount {
		t.Log("Any count\n")
		t.Logf("\t:             the sum of parts is: %d\n", anyCount)
		t.Logf("\t: the number of distinct names is: %d\n", len(names))
		t.Logf("\t:                            diff: %d\n",
			len(names)-anyCount)
		t.Error("full set of colour names should be <= it's constituents")
	}
}

func TestColourByName(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		f         Family
		cName     string
		expColour color.RGBA //nolint:misspell
	}{
		{
			ID:        testhelper.MkID("green"),
			f:         AnyColours,
			cName:     "green",
			expColour: color.RGBA{0x00, 0x80, 0x00, 0xff}, //nolint:misspell
		},
		{
			ID:        testhelper.MkID("lime"),
			f:         AnyColours,
			cName:     "lime",
			expColour: color.RGBA{0x00, 0xff, 0x00, 0xff}, //nolint:misspell
		},
		{
			ID:        testhelper.MkID("lime"),
			f:         HTMLColours,
			cName:     "lime",
			expColour: color.RGBA{0x00, 0xff, 0x00, 0xff}, //nolint:misspell
		},
		{
			ID:        testhelper.MkID("bad family"),
			ExpErr:    testhelper.MkExpErr(ErrBadFamily.Error()),
			f:         Family(99),
			cName:     "lime",
			expColour: color.RGBA{0x00, 0x00, 0x00, 0x00}, //nolint:misspell
		},
		{
			ID:        testhelper.MkID("bad colour (WebColours)"),
			ExpErr:    testhelper.MkExpErr(ErrBadColour.Error()),
			f:         WebColours,
			cName:     "NO SUCH COLOUR",
			expColour: color.RGBA{0x00, 0x00, 0x00, 0x00}, //nolint:misspell
		},
		{
			ID:        testhelper.MkID("bad colour (AnyColours)"),
			ExpErr:    testhelper.MkExpErr(ErrBadColour.Error()),
			f:         AnyColours,
			cName:     "NO SUCH COLOUR",
			expColour: color.RGBA{0x00, 0x00, 0x00, 0x00}, //nolint:misspell
		},
	}

	for _, tc := range testCases {
		c, err := tc.f.Colour(tc.cName)
		testhelper.CheckExpErr(t, err, tc)
		colourtesthelper.DiffRGBA(t, tc.IDStr(), tc.f.String()+".Colour",
			c, tc.expColour)
	}
}
