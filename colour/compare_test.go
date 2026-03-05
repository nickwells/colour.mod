package colour

import (
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestCompare(t *testing.T) {
	var (
		red     = rgba{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
		redish  = rgba{R: 0xfe, G: 0x01, B: 0x01, A: 0xfe}
		black   = rgba{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
		blackR1 = rgba{R: 0x01, G: 0x00, B: 0x00, A: 0xff}
	// 	blackish    = rgba{R: 0x01, G: 0x01, B: 0x01, A: 0xfe}
	// 	white    = rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	// 	whiteish    = rgba{R: 0xfe, G: 0xfe, B: 0xfe, A: 0xfe}
	)
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		c1, c2    rgba
		precision uint8
	}{
		{
			ID:        testhelper.MkID("same"),
			c1:        red,
			c2:        red,
			precision: 0,
		},
		{
			ID:        testhelper.MkID("red, redish, precision: 1"),
			c1:        red,
			c2:        redish,
			precision: 1,
		},
		{
			ID: testhelper.MkID("red, redish, precision: 0"),
			ExpErr: testhelper.MkExpErr(
				"colours ",
				"differ (precision: 0)",
				"R's differ by",
				"G's differ by",
				"B's differ by",
				"A's differ by",
			),
			c1:        red,
			c2:        redish,
			precision: 0,
		},
		{
			ID:        testhelper.MkID("red, redish, precision: 1"),
			c1:        red,
			c2:        redish,
			precision: 1,
		},
		{
			ID: testhelper.MkID("black, blackR1, precision: 0"),
			ExpErr: testhelper.MkExpErr(
				"colours ",
				"differ (precision: 0)",
				"R's differ by 1",
			),
			c1:        black,
			c2:        blackR1,
			precision: 0,
		},
		{
			ID:        testhelper.MkID("black, blackR1, precision: 1"),
			c1:        black,
			c2:        blackR1,
			precision: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := Compare(tc.c1, tc.c2, tc.precision)
			testhelper.CheckExpErrWithID(t, tc.IDStr(), err, tc)
		})
	}
}

func TestWithinDist(t *testing.T) {
	var (
		red   = rgba{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
		black = rgba{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
		white = rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	)
	testCases := []struct {
		testhelper.ID
		c1, c2 rgba
		dist   float64
		exp    bool
	}{
		{
			ID:   testhelper.MkID("same colour"),
			c1:   red,
			c2:   red,
			dist: 0,
			exp:  true,
		},
		{
			ID:   testhelper.MkID("black/white, dist: 0"),
			c1:   black,
			c2:   white,
			dist: 0,
			exp:  false,
		},
		{
			ID:   testhelper.MkID("black/white, dist: 441"),
			c1:   black,
			c2:   white,
			dist: 441,
			exp:  false,
		},
		{
			ID:   testhelper.MkID("black/white, dist: 442"),
			c1:   black,
			c2:   white,
			dist: 442,
			exp:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			act := WithinDist(tc.c1, tc.c2, tc.dist)
			testhelper.DiffBool(t,
				tc.IDStr(), "within distance",
				act, tc.exp)
		})
	}
}
