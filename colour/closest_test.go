package colour

import (
	"slices"
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestDistSquared(t *testing.T) {
	webOlive := webColours["olive"]
	webMaroon := webColours["maroon"]
	webGrey := webColours["grey"]
	testCases := []struct {
		testhelper.ID
		c, target rgba
		expDistSq int
	}{
		{
			ID:        testhelper.MkID("same colour, same Alpha"),
			c:         rgba{R: 1, G: 2, B: 3, A: 123},
			target:    rgba{R: 1, G: 2, B: 3, A: 123},
			expDistSq: 0,
		},
		{
			ID:        testhelper.MkID("same colour, different Alpha"),
			c:         rgba{R: 1, G: 2, B: 3, A: 123},
			target:    rgba{R: 1, G: 2, B: 3, A: 0},
			expDistSq: 0,
		},
		{
			ID:        testhelper.MkID("different colour, same Alpha"),
			c:         rgba{R: 1, G: 2, B: 3, A: 123},
			target:    rgba{R: 11, G: 12, B: 13, A: 123},
			expDistSq: 300,
		},
		{
			ID:        testhelper.MkID("different colour, different Alpha"),
			c:         rgba{R: 11, G: 12, B: 13, A: 123},
			target:    rgba{R: 1, G: 2, B: 3, A: 0},
			expDistSq: 300,
		},
		{
			ID:        testhelper.MkID("webOlive to webMaroon"),
			c:         webOlive,
			target:    webMaroon,
			expDistSq: 0x80 * 0x80,
		},
		{
			ID:        testhelper.MkID("webGrey to webMaroon"),
			c:         webGrey,
			target:    webMaroon,
			expDistSq: 0x80 * 0x80 * 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			d := distSquared(tc.target, tc.c)
			testhelper.DiffInt(t, tc.IDStr(), "distance", d, tc.expDistSq)
		})
	}
}

func TestClosestWithin(t *testing.T) {
	x11Cornsilk3 := x11Colours["cornsilk3"]
	webGrey := webColours["grey"]

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		fl     Families
		target rgba
		dist   float64
		expFCs []FamilyColour
	}{
		{
			ID:     testhelper.MkID("no family - use the StandardFamilies"),
			fl:     Families{},
			target: x11Cornsilk3,
			dist:   0,
			expFCs: []FamilyColour{
				{
					dist:   0,
					Family: X11Colours,
					CNames: []string{
						"cornsilk3",
					},
					Colour: x11Cornsilk3,
				},
			},
		},
		{
			ID: testhelper.MkID("bad - bad family"),
			ExpErr: testhelper.MkExpErr(
				`1 problem found: "nonesuch" is not` +
					` a valid Family (at position 0)`),
			fl:     Families{Family("nonesuch")},
			target: webGrey,
			dist:   0,
		},
		{
			ID: testhelper.MkID("bad - invalid proximity"),
			ExpErr: testhelper.MkExpErr(BadColourProximity +
				": the proximity must be greater than or equal to 0"),
			fl:     Families{WebColours},
			target: webGrey,
			dist:   -1,
		},
		{
			ID:     testhelper.MkID("one Family, 0 distance, 1 match"),
			fl:     Families{WebColours},
			target: webGrey,
			dist:   0,
			expFCs: []FamilyColour{
				{
					dist:   0,
					Family: WebColours,
					CNames: []string{
						"gray",
						"grey",
					},
					Colour: webGrey,
				},
			},
		},
		{
			ID:     testhelper.MkID("one Family, 42 distance, 1 match"),
			fl:     Families{WebColours},
			target: webGrey,
			dist:   42,
			expFCs: []FamilyColour{
				{
					dist:   0,
					Family: WebColours,
					CNames: []string{
						"gray",
						"grey",
					},
					Colour: webGrey,
				},
			},
		},
		{
			ID:     testhelper.MkID("standard families, 0 distance, 3 matches"),
			fl:     Families{StandardColours},
			target: webGrey,
			dist:   0,
			expFCs: []FamilyColour{
				{
					dist:   0,
					Family: CGAColours,
					CNames: []string{
						"dark gray",
						"dark grey",
						"dark-gray",
						"dark-grey",
					},
					Colour: webGrey,
				},
				{
					dist:   0,
					Family: HTMLColours,
					CNames: []string{
						"gray",
						"grey",
					},
					Colour: webGrey,
				},
				{
					dist:   0,
					Family: WebColours,
					CNames: []string{
						"gray",
						"grey",
					},
					Colour: webGrey,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actFCs, err := tc.fl.ClosestWithin(tc.target, tc.dist)
			testhelper.CheckExpErr(t, err, tc)

			if err == nil {
				slices.SortFunc(actFCs, FamilyColourCompare)

				dvErr := testhelper.DiffVals(actFCs, tc.expFCs)
				if dvErr != nil {
					t.Log(tc.IDStr())
					t.Log("\t: unexpected results:")

					for _, fc := range actFCs {
						t.Logf("\t:   %#v", fc)
					}

					t.Error("\t:", dvErr)
				}
			}
		})
	}
}

func TestClosestN(t *testing.T) {
	webOlive := webColours["olive"]
	webGreen := webColours["green"]
	webGrey := webColours["grey"]
	webMaroon := webColours["maroon"]
	x11Cornsilk3 := x11Colours["cornsilk3"]
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		fl     Families
		target rgba
		n      int
		expFCs []FamilyColour
	}{
		{
			ID:     testhelper.MkID("bad - negative n"),
			ExpErr: testhelper.MkExpErr("bad colour count: -1 - must be >= 0"),
			fl:     Families{WebColours},
			n:      -1,
		},
		{
			ID:     testhelper.MkID("no families - using StandardColours"),
			n:      1,
			target: x11Cornsilk3,
			expFCs: []FamilyColour{
				{
					Family: X11Colours,
					Colour: x11Cornsilk3,
					CNames: []string{"cornsilk3"},
				},
			},
		},
		{
			ID:     testhelper.MkID("1 result, 1 family"),
			fl:     Families{WebColours},
			n:      1,
			target: webOlive,
			expFCs: []FamilyColour{
				{
					Family: WebColours,
					Colour: webOlive,
					CNames: []string{"olive"},
				},
			},
		},
		{
			ID:     testhelper.MkID("0 requested, 0 returned"),
			fl:     Families{WebColours},
			n:      0,
			target: webOlive,
		},
		{
			ID:     testhelper.MkID("4 results, 1 family"),
			fl:     Families{WebColours},
			n:      4,
			target: webOlive,
			expFCs: []FamilyColour{
				{
					dist:   0,
					Family: WebColours,
					Colour: webOlive,
					CNames: []string{"olive"},
				},
				{
					dist:   0x80 * 0x80,
					Family: WebColours,
					Colour: webGreen,
					CNames: []string{"green"},
				},
				{
					dist:   0x80 * 0x80,
					Family: WebColours,
					Colour: webMaroon,
					CNames: []string{"maroon"},
				},
				{
					dist:   0x80 * 0x80,
					Family: WebColours,
					Colour: webGrey,
					CNames: []string{"gray", "grey"},
				},
			},
		},
		{
			ID:     testhelper.MkID("2 results, 1 family"),
			fl:     Families{WebColours},
			n:      2,
			target: webOlive,
			expFCs: []FamilyColour{
				{
					dist:   0,
					Family: WebColours,
					Colour: webOlive,
					CNames: []string{"olive"},
				},
				{
					dist:   0x80 * 0x80,
					Family: WebColours,
					Colour: webGreen,
					CNames: []string{"green"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actFCs, err := tc.fl.ClosestN(tc.target, tc.n)
			testhelper.CheckExpErr(t, err, tc)

			if err == nil {
				slices.SortFunc(actFCs, FamilyColourCompare)

				dvErr := testhelper.DiffVals(actFCs, tc.expFCs)
				if dvErr != nil {
					t.Log(tc.IDStr())
					t.Log("\t: unexpected results:")

					for _, fc := range actFCs {
						t.Logf("\t:   %#v", fc)
					}

					t.Error("\t:", dvErr)
				}
			}
		})
	}
}

func TestFamilyColour_FullNames(t *testing.T) {
	webGrey := webColours["grey"]
	testCases := []struct {
		testhelper.ID
		fam     Family
		cNames  []string
		expName string
	}{
		{
			ID:      testhelper.MkID("one cName"),
			fam:     "test",
			cNames:  []string{"cName1"},
			expName: `"test:cName1"`,
		},
		{
			ID:      testhelper.MkID("two cNames"),
			fam:     "test",
			cNames:  []string{"cName1", "cName2"},
			expName: `"test:cName1" or "test:cName2"`,
		},
		{
			ID:      testhelper.MkID("three cNames"),
			fam:     "test",
			cNames:  []string{"cName1", "cName2", "cName3"},
			expName: `"test:cName1", "test:cName2" or "test:cName3"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			fc := FamilyColour{
				dist:   0,
				Family: tc.fam,
				CNames: tc.cNames,
				Colour: webGrey,
			}
			name := fc.FullNames()
			testhelper.DiffString(t, tc.IDStr(), "FullName", name, tc.expName)
		})
	}
}
