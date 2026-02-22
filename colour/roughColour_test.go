package colour

import (
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestRoughly(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		c     rgba
		expRC RoughColour
	}{
		{
			ID:    testhelper.MkID("red shoud be roughly red"),
			c:     rgba{R: 0xff, A: 0xff},
			expRC: RoughlyRed,
		},
		{
			ID:    testhelper.MkID("green shoud be roughly green"),
			c:     rgba{G: 0xff, A: 0xff},
			expRC: RoughlyGreen,
		},
		{
			ID:    testhelper.MkID("blue shoud be roughly blue"),
			c:     rgba{B: 0xff, A: 0xff},
			expRC: RoughlyBlue,
		},
		{
			ID:    testhelper.MkID("cyan shoud be roughly cyan"),
			c:     rgba{G: 0xff, B: 0xff, A: 0xff},
			expRC: RoughlyCyan,
		},
		{
			ID:    testhelper.MkID("magenta shoud be roughly magenta"),
			c:     rgba{R: 0xff, B: 0xff, A: 0xff},
			expRC: RoughlyMagenta,
		},
		{
			ID:    testhelper.MkID("yellow shoud be roughly yellow"),
			c:     rgba{R: 0xff, G: 0xff, A: 0xff},
			expRC: RoughlyYellow,
		},
		{
			ID:    testhelper.MkID("black shoud be roughly black"),
			c:     rgba{R: 0, G: 0, B: 0, A: 0xff},
			expRC: RoughlyBlack,
		},
		{
			ID:    testhelper.MkID("white shoud be roughly white"),
			c:     rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
			expRC: RoughlyWhite,
		},
		{
			ID:    testhelper.MkID("grey shoud be roughly grey"),
			c:     rgba{R: 0x80, G: 0x80, B: 0x80, A: 0xff},
			expRC: RoughlyGrey,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actRC := Roughly(tc.c)
			testhelper.DiffInt(t, tc.IDStr(), "rough colour", actRC, tc.expRC)
		})
	}
}
