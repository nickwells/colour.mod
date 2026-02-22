package colour

import (
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestMakeColoursBetween(t *testing.T) {
	var (
		webBlack = webColours["black"]
		midGrey  = rgba{R: 0x7f, G: 0x7f, B: 0x7f, A: 0xff}
		webWhite = webColours["white"]
		webRed   = webColours["red"]
		webGreen = webColours["green"]
	)

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		count        int
		lower, upper rgba
		expColours   []rgba
	}{
		{
			ID: testhelper.MkID("bad call, count: 0"),
			ExpErr: testhelper.MkExpErr("bad colour count: " +
				"the count (0) must be greater than zero"),
		},
		{
			ID: testhelper.MkID("bad call, count: -1"),
			ExpErr: testhelper.MkExpErr("bad colour count: " +
				"the count (-1) must be greater than zero"),
			count: -1,
		},
		{
			ID:         testhelper.MkID("good call, count: 1"),
			count:      1,
			lower:      webBlack,
			upper:      webWhite,
			expColours: []rgba{midGrey},
		},
		{
			ID:         testhelper.MkID("good call, lower==upper"),
			count:      99,
			lower:      webBlack,
			upper:      webBlack,
			expColours: []rgba{webBlack},
		},
		{
			ID:         testhelper.MkID("good call, lower!=upper, count: 2"),
			count:      2,
			lower:      webRed,
			upper:      webGreen,
			expColours: []rgba{webRed, webGreen},
		},
		{
			ID:    testhelper.MkID("good call, lower!=upper, count: 4"),
			count: 4,
			lower: webRed,
			upper: webGreen,
			expColours: []rgba{
				webRed,
				{R: 212, G: 141, B: 0, A: 255},
				{R: 113, G: 170, B: 0, A: 255},
				webGreen,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actColours, err := MakeColoursBetween(
				tc.count, tc.lower, tc.upper)

			testhelper.CheckExpErr(t, err, tc)

			if err == nil {
				testhelper.DiffSlice(t, tc.IDStr(), "colours",
					actColours, tc.expColours)
			}
		})
	}
}
