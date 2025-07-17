package colour

import (
	"testing"

	"github.com/nickwells/colour.mod/v2/colourtesthelper"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestFamilyNames(t *testing.T) {
	var (
		webCount     = len(webColours)
		cgaCount     = len(cgaColours)
		x11Count     = len(x11Colours)
		htmlCount    = len(htmlColours)
		pantoneCount = len(pantoneColours)
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
		names, err := tc.f.ColourNames()
		if err != nil {
			t.Error(tc.IDStr(), "unexpected error:", err)
		}

		testhelper.DiffInt(t, tc.IDStr(), "number of colour names",
			len(names), tc.expNameCount)
	}
}

func TestColourByName(t *testing.T) {
	const (
		badFamily = "No-Such-Family"
		badColour = "No-Such-Colour"
	)

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		f         Family
		cName     string
		expColour rgba
	}{
		{
			ID:        testhelper.MkID("green"),
			f:         StandardColours,
			cName:     "green",
			expColour: rgba{R: 0x00, G: 0x80, B: 0x00, A: 0xff},
		},
		{
			ID:        testhelper.MkID("lime"),
			f:         StandardColours,
			cName:     "lime",
			expColour: rgba{R: 0x00, G: 0xff, B: 0x00, A: 0xff},
		},
		{
			ID:        testhelper.MkID("lime"),
			f:         HTMLColours,
			cName:     "lime",
			expColour: rgba{R: 0x00, G: 0xff, B: 0x00, A: 0xff},
		},
		{
			ID:        testhelper.MkID("bad family"),
			ExpErr:    testhelper.MkExpErr(badFamilyErr(badFamily).Error()),
			f:         Family(badFamily),
			cName:     "lime",
			expColour: rgba{R: 0x00, G: 0x00, B: 0x00, A: 0x00},
		},
		{
			ID:        testhelper.MkID("bad colour (WebColours)"),
			ExpErr:    testhelper.MkExpErr(badColourErr(badColour).Error()),
			f:         WebColours,
			cName:     badColour,
			expColour: rgba{R: 0x00, G: 0x00, B: 0x00, A: 0x00},
		},
		{
			ID:        testhelper.MkID("bad colour (StandardColours)"),
			ExpErr:    testhelper.MkExpErr(badColourErr(badColour).Error()),
			f:         StandardColours,
			cName:     badColour,
			expColour: rgba{R: 0x00, G: 0x00, B: 0x00, A: 0x00},
		},
	}

	for _, tc := range testCases {
		c, err := tc.f.Colour(tc.cName)
		testhelper.CheckExpErr(t, err, tc)
		colourtesthelper.DiffRGBA(t, tc.IDStr(), string(tc.f)+".Colour",
			c, tc.expColour)
	}
}
