package colour

import (
	"fmt"
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestDescribeColour(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		c      rgba
		expStr string
	}{
		{
			ID: testhelper.MkID("unnamed colour"),
			c:  rgba{R: 0x01, G: 0x02, B: 0x03, A: 0xff},
			//nolint:misspell
			expStr: "color.RGBA{R:0x01, G:0x02, B:0x03, A:0xff}",
		},
		{
			ID:     testhelper.MkID("just one name X11:blue2"),
			c:      rgba{R: 0x00, G: 0x00, B: 0xee, A: 0xff},
			expStr: "blue2",
		},
		{
			ID:     testhelper.MkID("many names, same family X11:grey1/gray1"),
			c:      rgba{R: 0x03, G: 0x03, B: 0x03, A: 0xff},
			expStr: "gray1",
		},
		{
			ID:     testhelper.MkID("many names and family aliases: black"),
			c:      rgba{R: 0x00, G: 0x00, B: 0x00, A: 0xff},
			expStr: "black",
		},
		{
			ID:     testhelper.MkID("many matches but with the same name"),
			c:      rgba{R: 0x00, G: 0x64, B: 0x00, A: 0xff},
			expStr: "darkgreen",
		},
		{
			ID:     testhelper.MkID("many matches, different names"),
			c:      rgba{R: 0xc0, G: 0xc0, B: 0xc0, A: 0xff},
			expStr: `"Web:silver", "HTML:silver" or "CGA:light gray"`,
		},
		{
			ID:     testhelper.MkID("many values, different names"),
			c:      rgba{R: 0x00, G: 0xff, B: 0x00, A: 0xff},
			expStr: `"Web:lime", "HTML:lime", "CGA:green" or "X11:green"`,
		},
	}

	for _, tc := range testCases {
		str := Describe(tc.c)
		testhelper.DiffString(t,
			tc.IDStr(), fmt.Sprintf("description of %v", tc.c),
			str, tc.expStr)
	}
}
