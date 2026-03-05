package colour

import (
	"testing"

	"github.com/nickwells/colour.mod/v2/colourtesthelper"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestSaturation(t *testing.T) {
	var (
		red = rgba{R: 0xff, A: 0xff}
		//green = rgba{G: 0xff, A: 0xff}
		//blue  = rgba{B: 0xff, A: 0xff}
	)

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		baseC      rgba
		saturation float64
	}{
		{
			ID: testhelper.MkID("bad saturation: <0"),
			ExpErr: testhelper.MkExpErr(
				"the saturation (-1.00) must be >= 0"),
			baseC:      red,
			saturation: -1,
		},
		{
			ID: testhelper.MkID("bad saturation: >1"),
			ExpErr: testhelper.MkExpErr(
				"the saturation (1.10) must be <= 1"),
			baseC:      red,
			saturation: 1.1,
		},
		{
			ID:         testhelper.MkID("good saturation: 0.5"),
			baseC:      red,
			saturation: 0.5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			satC, err := Saturation(tc.baseC, tc.saturation)
			testhelper.CheckExpErrWithID(t, tc.IDStr(), err, tc)

			if err == nil {
				satCHSL, _ := RGBA2HSLAndHSV(satC)
				baseCHSL, _ := RGBA2HSLAndHSV(tc.baseC)

				const epsilon = 0.01
				testhelper.DiffFloat(t, tc.IDStr(), "hue",
					satCHSL.Hue, baseCHSL.Hue, epsilon)
				testhelper.DiffFloat(t, tc.IDStr(), "luminance",
					satCHSL.Luminance, baseCHSL.Luminance, epsilon)
				testhelper.DiffFloat(t, tc.IDStr(), "saturation",
					satCHSL.Saturation, tc.saturation, epsilon)
			}
		})
	}
}

func TestLuminance(t *testing.T) {
	var (
		red = rgba{R: 0xff, A: 0xff}
		//green = rgba{G: 0xff, A: 0xff}
		//blue  = rgba{B: 0xff, A: 0xff}
	)

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		baseC     rgba
		luminance float64
	}{
		{
			ID: testhelper.MkID("bad luminance: <0"),
			ExpErr: testhelper.MkExpErr(
				"the luminance (-1.00) must be >= 0"),
			baseC:     red,
			luminance: -1,
		},
		{
			ID: testhelper.MkID("bad luminance: >1"),
			ExpErr: testhelper.MkExpErr(
				"the luminance (1.10) must be <= 1"),
			baseC:     red,
			luminance: 1.1,
		},
		{
			ID:        testhelper.MkID("good luminance: 0.5"),
			baseC:     red,
			luminance: 0.5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			satC, err := Luminance(tc.baseC, tc.luminance)
			testhelper.CheckExpErrWithID(t, tc.IDStr(), err, tc)

			if err == nil {
				satCHSL, _ := RGBA2HSLAndHSV(satC)
				baseCHSL, _ := RGBA2HSLAndHSV(tc.baseC)

				const epsilon = 0.01
				testhelper.DiffFloat(t, tc.IDStr(), "hue",
					satCHSL.Hue, baseCHSL.Hue, epsilon)
				testhelper.DiffFloat(t, tc.IDStr(), "saturaiton",
					satCHSL.Saturation, baseCHSL.Saturation, epsilon)
				testhelper.DiffFloat(t, tc.IDStr(), "luminance",
					satCHSL.Luminance, tc.luminance, epsilon)
			}
		})
	}
}

func TestInvert(t *testing.T) {
	var (
		red      = rgba{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
		invRed   = rgba{R: 0x00, G: 0xff, B: 0xff, A: 0xff}
		green    = rgba{R: 0x00, G: 0xff, B: 0x00, A: 0xff}
		invGreen = rgba{R: 0xff, G: 0x00, B: 0xff, A: 0xff}
		blue     = rgba{R: 0x00, G: 0x00, B: 0xff, A: 0xff}
		invBlue  = rgba{R: 0xff, G: 0xff, B: 0x00, A: 0xff}
		black    = rgba{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
		white    = rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	)

	testCases := []struct {
		testhelper.ID
		c    rgba
		expC rgba
	}{
		{
			ID:   testhelper.MkID("red"),
			c:    red,
			expC: invRed,
		},
		{
			ID:   testhelper.MkID("green"),
			c:    green,
			expC: invGreen,
		},
		{
			ID:   testhelper.MkID("blue"),
			c:    blue,
			expC: invBlue,
		},
		{
			ID:   testhelper.MkID("black"),
			c:    black,
			expC: white,
		},
		{
			ID:   testhelper.MkID("white"),
			c:    white,
			expC: black,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actC := Invert(tc.c)
			colourtesthelper.DiffRGB(t, tc.IDStr(), "invert", actC, tc.expC)
		})
	}
}

func TestComplement(t *testing.T) {
	var (
		red   = rgba{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
		green = rgba{R: 0x00, G: 0xff, B: 0x00, A: 0xff}
		blue  = rgba{R: 0x00, G: 0x00, B: 0xff, A: 0xff}
		black = rgba{R: 0x00, G: 0x00, B: 0x00, A: 0xff}
		grey8 = rgba{R: 0x88, G: 0x88, B: 0x88, A: 0xff}
		white = rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff}
	)

	testCases := []struct {
		testhelper.ID
		c              rgba
		changeExpected bool
	}{
		{
			ID:             testhelper.MkID("red"),
			c:              red,
			changeExpected: true,
		},
		{
			ID:             testhelper.MkID("green"),
			c:              green,
			changeExpected: true,
		},
		{
			ID:             testhelper.MkID("blue"),
			c:              blue,
			changeExpected: true,
		},
		{
			ID: testhelper.MkID("black"),
			c:  black,
		},
		{
			ID: testhelper.MkID("grey"),
			c:  grey8,
		},
		{
			ID: testhelper.MkID("white"),
			c:  white,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			cc := Complement(tc.c)

			if !tc.changeExpected {
				colourtesthelper.DiffRGB(t,
					tc.IDStr(), "no change expected",
					cc, tc.c)
				return
			}

			ccHSL, _ := RGBA2HSLAndHSV(cc)
			cHSL, _ := RGBA2HSLAndHSV(tc.c)

			const (
				maxDegrees     = 360
				halfMaxDegrees = 180
			)

			expHue := cHSL.Hue + halfMaxDegrees
			if expHue > maxDegrees {
				expHue -= maxDegrees
			}

			const epsilon = 0.01
			testhelper.DiffFloat(t, tc.IDStr(), "hue",
				ccHSL.Hue, expHue, epsilon)
			testhelper.DiffFloat(t, tc.IDStr(), "saturaiton",
				ccHSL.Saturation, cHSL.Saturation, epsilon)
			testhelper.DiffFloat(t, tc.IDStr(), "luminance",
				ccHSL.Luminance, cHSL.Luminance, epsilon)
		})
	}
}
