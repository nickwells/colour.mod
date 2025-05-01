package colour_test

import (
	"image/color" //nolint:misspell
	"testing"

	"github.com/nickwells/colour.mod/colour"
	"github.com/nickwells/colour.mod/colourtesthelper"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestCGACololur(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		i         int
		expColour color.RGBA //nolint:misspell
	}{
		{
			ID:        testhelper.MkID("0 (black)"),
			i:         0,
			expColour: color.RGBA{0, 0, 0, 0xff}, //nolint:misspell
		},
		{
			ID:        testhelper.MkID("15 (white)"),
			i:         15,
			expColour: color.RGBA{0xff, 0xff, 0xff, 0xff}, //nolint:misspell
		},
		{
			ID:        testhelper.MkID("bad index:99"),
			ExpErr:    testhelper.MkExpErr("bad CGA colour index: 99"),
			i:         99,
			expColour: color.RGBA{0, 0, 0, 0}, //nolint:misspell
		},
	}

	for _, tc := range testCases {
		c, err := colour.CGAColour(tc.i)
		testhelper.CheckExpErr(t, err, tc)
		colourtesthelper.DiffRGBA(t, tc.IDStr(), "CGA by number",
			c, tc.expColour)
	}
}
