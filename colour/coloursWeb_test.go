package colour

import (
	"testing"

	"github.com/nickwells/colour.mod/v2/colourtesthelper"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestCGACololur(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		i         int
		expColour rgba
	}{
		{
			ID:        testhelper.MkID("0 (black)"),
			i:         0,
			expColour: rgba{R: 0, G: 0, B: 0, A: 0xff},
		},
		{
			ID:        testhelper.MkID("15 (white)"),
			i:         15,
			expColour: rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		},
		{
			ID:        testhelper.MkID("bad index:99"),
			ExpErr:    testhelper.MkExpErr("bad CGA colour index: 99"),
			i:         99,
			expColour: rgba{R: 0, G: 0, B: 0, A: 0},
		},
	}

	for _, tc := range testCases {
		c, err := CGAColour(tc.i)
		testhelper.CheckExpErr(t, err, tc)
		colourtesthelper.DiffRGBA(t, tc.IDStr(), "CGA by number",
			c, tc.expColour)
	}
}
