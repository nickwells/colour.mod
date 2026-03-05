package colour

import (
	"math"
	"testing"

	"github.com/nickwells/colour.mod/v2/colourtesthelper"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestToGrey(t *testing.T) {
	ffByPALRed := uint8(math.Round(0xff * wtRedPAL))
	redGrey := rgba{R: ffByPALRed, G: ffByPALRed, B: ffByPALRed, A: 0xff}
	ffByPALGreen := uint8(math.Round(0xff * wtGreenPAL))
	greenGrey := rgba{R: ffByPALGreen, G: ffByPALGreen, B: ffByPALGreen, A: 0xff}
	ffByPALBlue := uint8(math.Round(0xff * wtBluePAL))
	blueGrey := rgba{R: ffByPALBlue, G: ffByPALBlue, B: ffByPALBlue, A: 0xff}
	testCases := []struct {
		testhelper.ID
		c       rgba
		expGrey rgba
	}{
		{
			ID:      testhelper.MkID("black goes to black"),
			c:       rgba{A: 0xff},
			expGrey: rgba{A: 0xff},
		},
		{
			ID:      testhelper.MkID("white goes to white"),
			c:       rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
			expGrey: rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		},
		{
			ID:      testhelper.MkID("red goes to PAL redGrey"),
			c:       rgba{R: 0xff, A: 0xff},
			expGrey: redGrey,
		},
		{
			ID:      testhelper.MkID("green goes to PAL greenGrey"),
			c:       rgba{G: 0xff, A: 0xff},
			expGrey: greenGrey,
		},
		{
			ID:      testhelper.MkID("blue goes to PAL blueGrey"),
			c:       rgba{B: 0xff, A: 0xff},
			expGrey: blueGrey,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actGrey := ToGrey(tc.c)

			colourtesthelper.DiffRGB(t, tc.IDStr(), "grey",
				actGrey, tc.expGrey)
		})
	}
}

func TestToGreyBT709(t *testing.T) {
	ffByBT709Red := uint8(math.Round(0xff * wtRedBT709))
	ffByBT709Green := uint8(math.Round(0xff * wtGreenBT709))
	ffByBT709Blue := uint8(math.Round(0xff * wtBlueBT709))
	redGrey := MakeGrey(ffByBT709Red)
	greenGrey := MakeGrey(ffByBT709Green)
	blueGrey := MakeGrey(ffByBT709Blue)
	testCases := []struct {
		testhelper.ID
		c       rgba
		expGrey rgba
	}{
		{
			ID:      testhelper.MkID("black goes to black"),
			c:       rgba{A: 0xff},
			expGrey: rgba{A: 0xff},
		},
		{
			ID:      testhelper.MkID("white goes to white"),
			c:       rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
			expGrey: rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		},
		{
			ID:      testhelper.MkID("red goes to BT709 redGrey"),
			c:       rgba{R: 0xff, A: 0xff},
			expGrey: redGrey,
		},
		{
			ID:      testhelper.MkID("green goes to BT709 greenGrey"),
			c:       rgba{G: 0xff, A: 0xff},
			expGrey: greenGrey,
		},
		{
			ID:      testhelper.MkID("blue goes to BT709 blueGrey"),
			c:       rgba{B: 0xff, A: 0xff},
			expGrey: blueGrey,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actGrey := ToGreyBT709(tc.c)

			colourtesthelper.DiffRGB(t, tc.IDStr(), "grey",
				actGrey, tc.expGrey)
		})
	}
}

func TestToGreyBT2100(t *testing.T) {
	ffByBT2100Red := uint8(math.Round(0xff * wtRedBT2100))
	ffByBT2100Green := uint8(math.Round(0xff * wtGreenBT2100))
	ffByBT2100Blue := uint8(math.Round(0xff * wtBlueBT2100))
	redGrey := MakeGrey(ffByBT2100Red)
	greenGrey := MakeGrey(ffByBT2100Green)
	blueGrey := MakeGrey(ffByBT2100Blue)
	testCases := []struct {
		testhelper.ID
		c       rgba
		expGrey rgba
	}{
		{
			ID:      testhelper.MkID("black goes to black"),
			c:       rgba{A: 0xff},
			expGrey: rgba{A: 0xff},
		},
		{
			ID:      testhelper.MkID("white goes to white"),
			c:       rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
			expGrey: rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		},
		{
			ID:      testhelper.MkID("red goes to BT2100 redGrey"),
			c:       rgba{R: 0xff, A: 0xff},
			expGrey: redGrey,
		},
		{
			ID:      testhelper.MkID("green goes to BT2100 greenGrey"),
			c:       rgba{G: 0xff, A: 0xff},
			expGrey: greenGrey,
		},
		{
			ID:      testhelper.MkID("blue goes to BT2100 blueGrey"),
			c:       rgba{B: 0xff, A: 0xff},
			expGrey: blueGrey,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actGrey := ToGreyBT2100(tc.c)

			colourtesthelper.DiffRGB(t, tc.IDStr(), "grey",
				actGrey, tc.expGrey)
		})
	}
}

func TestToGreyEqual(t *testing.T) {
	weight := 1.0 / 3.0
	ffByEqual := uint8(math.Round(0xff * weight))
	grey := MakeGrey(ffByEqual)
	testCases := []struct {
		testhelper.ID
		c       rgba
		expGrey rgba
	}{
		{
			ID:      testhelper.MkID("black goes to black"),
			c:       rgba{A: 0xff},
			expGrey: rgba{A: 0xff},
		},
		{
			ID:      testhelper.MkID("white goes to white"),
			c:       rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
			expGrey: rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		},
		{
			ID:      testhelper.MkID("red goes to Equal redGrey"),
			c:       rgba{R: 0xff, A: 0xff},
			expGrey: grey,
		},
		{
			ID:      testhelper.MkID("green goes to Equal greenGrey"),
			c:       rgba{G: 0xff, A: 0xff},
			expGrey: grey,
		},
		{
			ID:      testhelper.MkID("blue goes to Equal blueGrey"),
			c:       rgba{B: 0xff, A: 0xff},
			expGrey: grey,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actGrey := ToGreyEqual(tc.c)

			colourtesthelper.DiffRGB(t, tc.IDStr(), "grey",
				actGrey, tc.expGrey)
		})
	}
}
