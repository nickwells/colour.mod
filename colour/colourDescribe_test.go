package colour_test

import (
	"image/color" //nolint:misspell
	"testing"

	"github.com/nickwells/colour.mod/colour"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestDescribeColour(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		c      color.RGBA //nolint:misspell
		expStr string
	}{
		{
			ID:     testhelper.MkID("unnamed colour"),
			c:      color.RGBA{0x01, 0x02, 0x03, 0xff},        //nolint:misspell
			expStr: "color.RGBA{R:0x1, G:0x2, B:0x3, A:0xff}", //nolint:misspell
		},
		{
			ID:     testhelper.MkID("just one name X11:blue2"),
			c:      color.RGBA{0x00, 0x00, 0xee, 0xff}, //nolint:misspell
			expStr: "blue2",
		},
		{
			ID:     testhelper.MkID("many names, same family X11:grey1/gray1"),
			c:      color.RGBA{0x03, 0x03, 0x03, 0xff}, //nolint:misspell
			expStr: "gray1",
		},
		{
			ID:     testhelper.MkID("many names and family aliases: black"),
			c:      color.RGBA{0x00, 0x00, 0x00, 0xff}, //nolint:misspell
			expStr: "black",
		},
		{
			ID:     testhelper.MkID("many matches but with the same name"),
			c:      color.RGBA{0x00, 0x64, 0x00, 0xff}, //nolint:misspell
			expStr: "darkgreen",
		},
		{
			ID:     testhelper.MkID("many matches, different names"),
			c:      color.RGBA{0xc0, 0xc0, 0xc0, 0xff}, //nolint:misspell
			expStr: "silver (Web and HTML) or light gray (CGA)",
		},
		{
			ID:     testhelper.MkID("many values, different names"),
			c:      color.RGBA{0x00, 0xff, 0x00, 0xff}, //nolint:misspell
			expStr: "lime (Web and HTML) or green (CGA and X11)",
		},
	}

	for _, tc := range testCases {
		str := colour.Describe(tc.c)
		testhelper.DiffString(t, tc.IDStr(), "description",
			str, tc.expStr)
	}
}
