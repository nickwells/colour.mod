package colour

import (
	"regexp"
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestColoursMatchingByRegexp(t *testing.T) {
	webGrey := webColours["grey"]

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		fl           Families
		re           *regexp.Regexp
		namedColours []NamedColour
	}{
		{
			ID: testhelper.MkID("non-matching regexp - no results expected"),
			fl: Families{WebColours},
			re: regexp.MustCompilePOSIX("blah, blah, blah"),
		},
		{
			ID: testhelper.MkID("regexp (gr.y) - 2 results expected"),
			fl: Families{WebColours},
			re: regexp.MustCompilePOSIX("gr.y"),
			namedColours: []NamedColour{
				{
					name:   "web:gray",
					colour: webGrey,
				},
				{
					name:   "web:grey",
					colour: webGrey,
				},
			},
		},
		{
			ID:     testhelper.MkID("bad family - error expected"),
			ExpErr: testhelper.MkExpErr(`bad colour family: "nonesuch"`),
			fl:     Families{Family("nonesuch")},
			re:     regexp.MustCompilePOSIX("gr.y"),
		},
		{
			ID:     testhelper.MkID("no regexp - error expected"),
			ExpErr: testhelper.MkExpErr("no regular expression was provided"),
			fl:     Families{Family(WebColours)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ncs, err := ColoursMatchingByRegexp(tc.fl, tc.re)
			testhelper.CheckExpErr(t, err, tc)

			if err == nil {
				dvErr := testhelper.DiffVals(ncs, tc.namedColours)
				if dvErr != nil {
					t.Log(tc.IDStr())
					t.Log("\t: bad NamedColours list:")
					for _, nc := range ncs {
						t.Logf("\t:    %#v", nc)
					}
					t.Error("\t: differences:", dvErr)
				}
			}
		})
	}
}

func TestColoursMatchingByFunc(t *testing.T) {
	webGrey := webColours["grey"]

	alwaysFalseFunc := func(_ string, _ rgba) bool { return false }

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		fl           Families
		f            ColourMatchingFunc
		namedColours []NamedColour
	}{
		{
			ID: testhelper.MkID("no matches - no results expected"),
			fl: Families{WebColours},
			f:  alwaysFalseFunc,
		},
		{
			ID: testhelper.MkID("grey/grey - 2 results expected"),
			fl: Families{WebColours},
			f: func(name string, _ rgba) bool {
				if name == "grey" || name == "gray" {
					return true
				}
				return false
			},
			namedColours: []NamedColour{
				{
					name:   "web:gray",
					colour: webGrey,
				},
				{
					name:   "web:grey",
					colour: webGrey,
				},
			},
		},
		{
			ID:     testhelper.MkID("bad family - error expected"),
			ExpErr: testhelper.MkExpErr(`bad colour family: "nonesuch"`),
			fl:     Families{Family("nonesuch")},
			f:      alwaysFalseFunc,
		},
		{
			ID:     testhelper.MkID("no matchFunc - error expected"),
			ExpErr: testhelper.MkExpErr("no match function was provided"),
			fl:     Families{Family(WebColours)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ncs, err := ColoursMatchingByFunc(tc.fl, tc.f)
			testhelper.CheckExpErr(t, err, tc)

			if err == nil {
				dvErr := testhelper.DiffVals(ncs, tc.namedColours)
				if dvErr != nil {
					t.Log(tc.IDStr())
					t.Log("\t: bad NamedColours list:")
					for _, nc := range ncs {
						t.Logf("\t:    %#v", nc)
					}
					t.Error("\t: differences:", dvErr)
				}
			}
		})
	}
}

func reportNCDiffs(t *testing.T,
	nc NamedColour, expName string, expColour rgba,
) {
	t.Helper()

	if nc.Name() != expName {
		t.Log("MakeNamedColour failure")
		t.Logf("\t: expected: %q", expName)
		t.Logf("\t:   actual: %q", nc.Name())
		t.Error("\t: unexpected name")
	}

	if nc.Colour() != expColour {
		t.Log("MakeNamedColour failure")
		t.Logf("\t: expected: %#v", expColour)
		t.Logf("\t:   actual: %#v", nc.Colour())
		t.Error("\t: unexpected colour")
	}
}

func TestMakeNamedColour(t *testing.T) {
	webGrey := webColours["grey"]

	{
		name := "hello"

		nc := MakeNamedColour(name, webGrey)

		reportNCDiffs(t, nc, name, webGrey)
	}

	{
		fc := FamilyColour{
			dist:   0,
			Family: "f",
			CNames: []string{"grey", "also grey"},
			Colour: webGrey,
		}
		expName := fc.FullNames()
		nc := MakeNamedColourFromFamilyColour(fc)

		reportNCDiffs(t, nc, expName, webGrey)
	}

	{
		fc := FamilyColour{
			dist:   0,
			Family: "f",
			CNames: []string{"grey", "also grey"},
			Colour: webGrey,
		}
		ncs := MakeDistinctNamedColoursFromFamilyColour(fc)

		for i, nc := range ncs {
			expName := string(fc.Family) + ":" + fc.CNames[i]
			reportNCDiffs(t, nc, expName, webGrey)
		}
	}
}

func TestParseNamedColour(t *testing.T) {
	webGrey := webColours["grey"]
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		fl    Families
		s     string
		expNC NamedColour
	}{
		{
			ID: testhelper.MkID("bad RGB string - missing }"),
			ExpErr: testhelper.MkExpErr(
				"the colour definition starts with " +
					`"RGB {" but has no trailing "}"`),
			fl: Families{WebColours},
			s:  "RGB {",
		},
		{
			ID: testhelper.MkID("bad RGB string - R missing :"),
			ExpErr: testhelper.MkExpErr(`bad colour component: "R",` +
				" the name and value should be separated by a colon(:)"),
			fl: Families{WebColours},
			s:  "RGB {R}",
		},
		{
			ID: testhelper.MkID("bad RGB string - G missing :"),
			ExpErr: testhelper.MkExpErr(`bad colour component: "G",` +
				" the name and value should be separated by a colon(:)"),
			fl: Families{WebColours},
			s:  "RGB {R:0,G}",
		},
		{
			ID: testhelper.MkID("bad RGB string - R bad value"),
			ExpErr: testhelper.MkExpErr(`cannot convert the "R"` +
				` value ("hello") to a valid number: invalid syntax`),
			fl: Families{WebColours},
			s:  "RGB {R:hello}",
		},
		{
			ID: testhelper.MkID("bad RGB string - R negative value"),
			ExpErr: testhelper.MkExpErr(`cannot convert the "R"` +
				` value ("-42") to a valid number: invalid syntax`),
			fl: Families{WebColours},
			s:  "RGB {R:-42}",
		},
		{
			ID: testhelper.MkID("bad RGB string - R just too big"),
			ExpErr: testhelper.MkExpErr(`cannot convert the "R"` +
				` value ("0x100") to a valid number: value out of range`),
			fl: Families{WebColours},
			s:  "RGB {R:0x100}",
		},
		{
			ID: testhelper.MkID("bad RGB string - R too big"),
			ExpErr: testhelper.MkExpErr(`cannot convert the "R"` +
				` value ("0xfff") to a valid number: value out of range`),
			fl: Families{WebColours},
			s:  "RGB {R:0xfff}",
		},
		{
			ID: testhelper.MkID("bad RGB string - unknown component"),
			ExpErr: testhelper.MkExpErr(`unknown colour component: "X",` +
				` allowed values: "A", "B", "G" or "R"`),
			fl: Families{WebColours},
			s:  "RGB {X:0xff}",
		},
		{
			ID: testhelper.MkID("good RGB string - R = 0xff"),
			fl: Families{WebColours},
			s:  "RGB {R:0xff}",
			expNC: NamedColour{
				name:   "RGB {R:0xff}",
				colour: rgba{R: 0xff, A: 0xff},
			},
		},
		{
			ID: testhelper.MkID("good RGB string - r = 0xff"),
			fl: Families{WebColours},
			s:  "RGB {r:0xff}",
			expNC: NamedColour{
				name:   "RGB {r:0xff}",
				colour: rgba{R: 0xff, A: 0xff},
			},
		},
		{
			ID: testhelper.MkID("good RGB string - all LC and spaced"),
			fl: Families{WebColours},
			s:  "rgb { r : 0xff , g : 0 , b : 0 , a : 0 } ",
			expNC: NamedColour{
				name:   "rgb { r : 0xff , g : 0 , b : 0 , a : 0 } ",
				colour: rgba{R: 0xff},
			},
		},
		{
			ID:     testhelper.MkID("bad Family:Colour - family: nonesuch"),
			ExpErr: testhelper.MkExpErr(`bad colour family name: "nonesuch"`),
			fl:     Families{WebColours},
			s:      "nonesuch:grey",
		},
		{
			ID: testhelper.MkID("bad Family:Colour - family: webb"),
			ExpErr: testhelper.MkExpErr(`bad colour family name: "webb",` +
				` did you mean "web"`),
			fl: Families{WebColours},
			s:  "webb:grey",
		},
		{
			ID:     testhelper.MkID("bad Family:Colour - colour: nonesuch"),
			ExpErr: testhelper.MkExpErr(`bad colour name: "nonesuch"`),
			fl:     Families{WebColours},
			s:      "web:nonesuch",
		},
		{
			ID: testhelper.MkID("bad Family:Colour - colour: maroonn"),
			ExpErr: testhelper.MkExpErr(`bad colour name: "maroonn"` +
				`, did you mean "maroon"?`),
			fl: Families{WebColours},
			s:  "web:maroonn",
		},
		{
			ID: testhelper.MkID("bad Family:Colour - colour: MAROONN"),
			ExpErr: testhelper.MkExpErr(`bad colour name: "maroonn"` +
				`, did you mean "maroon"?`),
			fl: Families{WebColours},
			s:  "web:MAROONN",
		},
		{
			ID: testhelper.MkID("good Family:Colour - web:grey"),
			fl: Families{WebColours},
			s:  "web:grey",
			expNC: NamedColour{
				name:   "web:grey",
				colour: webGrey,
			},
		},
		{
			ID: testhelper.MkID("good Family:Colour - web:grey - with space"),
			fl: Families{WebColours},
			s:  " web : grey ",
			expNC: NamedColour{
				name:   " web : grey ",
				colour: webGrey,
			},
		},
		{
			ID: testhelper.MkID("good Family:Colour - web:grey - space and UC"),
			fl: Families{WebColours},
			s:  " WEB : GREY ",
			expNC: NamedColour{
				name:   " WEB : GREY ",
				colour: webGrey,
			},
		},
		{
			ID: testhelper.MkID("bad Colour name - colour: maroonn"),
			ExpErr: testhelper.MkExpErr(`bad colour name: "maroonn"` +
				`, did you mean "maroon"?`),
			fl: Families{WebColours},
			s:  "maroonn",
		},
		{
			ID: testhelper.MkID("bad Colour name - colour: MAROONN"),
			ExpErr: testhelper.MkExpErr(`bad colour name: "maroonn"` +
				`, did you mean "maroon"?`),
			fl: Families{WebColours},
			s:  "MAROONN",
		},
		{
			ID: testhelper.MkID("good Colour name - grey"),
			fl: Families{WebColours},
			s:  "grey",
			expNC: NamedColour{
				name:   "grey",
				colour: webGrey,
			},
		},
		{
			ID: testhelper.MkID("good Colour name - grey - space and UC"),
			fl: Families{WebColours},
			s:  " GREY ",
			expNC: NamedColour{
				name:   " GREY ",
				colour: webGrey,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			nc, err := ParseNamedColour(tc.fl, tc.s)

			testhelper.CheckExpErr(t, err, tc)

			if err == nil {
				dvErr := testhelper.DiffVals(nc, tc.expNC)
				if dvErr != nil {
					t.Log("ParseNamedColour failure")
					t.Logf("\t: expected: %#v", tc.expNC)
					t.Logf("\t:   actual: %#v", nc)
					t.Error("\t: unexpected NamedColour")
				}
			}
		})
	}
}
