package colour

import (
	"slices"
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestGenerateAltSpellings(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		s          string
		expResults []string
	}{
		{
			ID: testhelper.MkID("blah - nothing to change"),
			s:  "blah",
			expResults: []string{
				"blah",
			},
		},
		{
			ID: testhelper.MkID("gray"),
			s:  "gray",
			expResults: []string{
				"gray",
				"grey",
			},
		},
		{
			ID: testhelper.MkID("grey"),
			s:  "grey",
			expResults: []string{
				"gray",
				"grey",
			},
		},
		{
			ID: testhelper.MkID("gray at the start"),
			s:  "gray etc",
			expResults: []string{
				"gray etc",
				"grey etc",
			},
		},
		{
			ID: testhelper.MkID("grey at the start"),
			s:  "grey etc",
			expResults: []string{
				"gray etc",
				"grey etc",
			},
		},
		{
			ID: testhelper.MkID("gray at the end"),
			s:  "then gray",
			expResults: []string{
				"then gray",
				"then grey",
			},
		},
		{
			ID: testhelper.MkID("grey at the end"),
			s:  "then grey",
			expResults: []string{
				"then gray",
				"then grey",
			},
		},
		{
			ID: testhelper.MkID("gray in the middle"),
			s:  "then gray etc",
			expResults: []string{
				"then gray etc",
				"then grey etc",
			},
		},
		{
			ID: testhelper.MkID("grey in the middle"),
			s:  "then grey etc",
			expResults: []string{
				"then gray etc",
				"then grey etc",
			},
		},
		{
			ID: testhelper.MkID("greyson"),
			s:  "greyson",
			expResults: []string{
				"greyson",
			},
		},
		{
			ID: testhelper.MkID("grayson"),
			s:  "grayson",
			expResults: []string{
				"grayson",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			results := generateAltSpellings(tc.s)
			slices.Sort(results)
			slices.Sort(tc.expResults)

			testhelper.DiffStringSlice(t,
				tc.IDStr(), "AltSpellings",
				results, tc.expResults)
		})
	}
}

func TestTransformColourNames(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		names      []string
		expResults map[string]bool
	}{
		{
			ID: testhelper.MkID("no changes"),
			names: []string{
				"first",
				"second",
				"third",
			},
			expResults: map[string]bool{
				"first":  true,
				"second": true,
				"third":  true,
			},
		},
		{
			ID: testhelper.MkID("apostrophe changes"),
			names: []string{
				"noChange",
				"one's",
				"two'sAndThree's",
			},
			expResults: map[string]bool{
				"noChange":        true,
				"one's":           true,
				"ones":            true,
				"two'sAndThree's": true,
				"twosAndThrees":   true,
			},
		},
		{
			ID: testhelper.MkID("space and apostrophe changes"),
			names: []string{
				"noChange",
				"one's and two's",
				"three's and four's",
			},
			expResults: map[string]bool{
				"noChange":           true,
				"one's and two's":    true,
				"ones-and-twos":      true,
				"three's and four's": true,
				"threes-and-fours":   true,
			},
		},
		{
			ID: testhelper.MkID("space changes"),
			names: []string{
				"noChange",
				"one and two",
				"three and four",
			},
			expResults: map[string]bool{
				"noChange":       true,
				"one and two":    true,
				"one-and-two":    true,
				"three and four": true,
				"three-and-four": true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			results := transformColourNames(tc.names)
			for s := range results {
				if !tc.expResults[s] {
					t.Log(tc.IDStr())
					t.Errorf("\t: transformed name %q was not expected", s)
				}
			}

			for s := range tc.expResults {
				if !results[s] {
					t.Log(tc.IDStr())
					t.Errorf(
						"\t: transformed name %q was expected but not found", s)
				}
			}
		})
	}
}

func TestAddColourAliases(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		orig, exp colourNameToRGBA
	}{
		{
			ID: testhelper.MkID("example"),
			orig: colourNameToRGBA{
				"wevet":              {R: 0xee, G: 0xe9, B: 0xe7, A: 0xff},
				"ammonite":           {R: 0xdd, G: 0xd8, B: 0xcf, A: 0xff},
				"purbeck stone":      {R: 0xc4, G: 0xbe, B: 0xb4, A: 0xff},
				"cornforth white":    {R: 0xd1, G: 0xcb, B: 0xc3, A: 0xff},
				"all white":          {R: 0xfb, G: 0xf8, B: 0xf4, A: 0xff},
				"blackened":          {R: 0xdd, G: 0xdb, B: 0xd9, A: 0xff},
				"dimpse":             {R: 0xd9, G: 0xd8, B: 0xd3, A: 0xff},
				"pavilion gray":      {R: 0xc8, G: 0xc3, B: 0xbc, A: 0xff},
				"strong white":       {R: 0xe5, G: 0xe0, B: 0xdb, A: 0xff},
				"skimming stone":     {R: 0xdf, G: 0xd6, B: 0xcb, A: 0xff},
				"elephant's breath":  {R: 0xcc, G: 0xbf, B: 0xb3, A: 0xff},
				"dove tale":          {R: 0xbb, G: 0xb1, B: 0xab, A: 0xff},
				"school house white": {R: 0xe6, G: 0xdf, B: 0xd1, A: 0xff},
				"gray":               {R: 0x80, G: 0x80, B: 0x80, A: 0xFF},
				"grey":               {R: 0x80, G: 0x80, B: 0x80, A: 0xFF},
			},
			exp: colourNameToRGBA{
				"wevet":              {R: 0xee, G: 0xe9, B: 0xe7, A: 0xff},
				"ammonite":           {R: 0xdd, G: 0xd8, B: 0xcf, A: 0xff},
				"purbeck stone":      {R: 0xc4, G: 0xbe, B: 0xb4, A: 0xff},
				"purbeck-stone":      {R: 0xc4, G: 0xbe, B: 0xb4, A: 0xff},
				"cornforth white":    {R: 0xd1, G: 0xcb, B: 0xc3, A: 0xff},
				"cornforth-white":    {R: 0xd1, G: 0xcb, B: 0xc3, A: 0xff},
				"all white":          {R: 0xfb, G: 0xf8, B: 0xf4, A: 0xff},
				"all-white":          {R: 0xfb, G: 0xf8, B: 0xf4, A: 0xff},
				"blackened":          {R: 0xdd, G: 0xdb, B: 0xd9, A: 0xff},
				"dimpse":             {R: 0xd9, G: 0xd8, B: 0xd3, A: 0xff},
				"pavilion gray":      {R: 0xc8, G: 0xc3, B: 0xbc, A: 0xff},
				"pavilion-gray":      {R: 0xc8, G: 0xc3, B: 0xbc, A: 0xff},
				"pavilion grey":      {R: 0xc8, G: 0xc3, B: 0xbc, A: 0xff},
				"pavilion-grey":      {R: 0xc8, G: 0xc3, B: 0xbc, A: 0xff},
				"strong white":       {R: 0xe5, G: 0xe0, B: 0xdb, A: 0xff},
				"strong-white":       {R: 0xe5, G: 0xe0, B: 0xdb, A: 0xff},
				"skimming stone":     {R: 0xdf, G: 0xd6, B: 0xcb, A: 0xff},
				"skimming-stone":     {R: 0xdf, G: 0xd6, B: 0xcb, A: 0xff},
				"elephant's breath":  {R: 0xcc, G: 0xbf, B: 0xb3, A: 0xff},
				"elephants-breath":   {R: 0xcc, G: 0xbf, B: 0xb3, A: 0xff},
				"dove tale":          {R: 0xbb, G: 0xb1, B: 0xab, A: 0xff},
				"dove-tale":          {R: 0xbb, G: 0xb1, B: 0xab, A: 0xff},
				"school house white": {R: 0xe6, G: 0xdf, B: 0xd1, A: 0xff},
				"school-house-white": {R: 0xe6, G: 0xdf, B: 0xd1, A: 0xff},
				"gray":               {R: 0x80, G: 0x80, B: 0x80, A: 0xFF},
				"grey":               {R: 0x80, G: 0x80, B: 0x80, A: 0xFF},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			result := withAliases(tc.orig)

			for name, actColour := range result {
				expColour, ok := tc.exp[name]
				if !ok {
					t.Log(tc.IDStr())
					t.Logf("\t: colour[%q] = %v", name, actColour)
					t.Errorf("\t: unexpected colour %q", name)

					continue
				}

				if actColour != expColour {
					t.Log(tc.IDStr())
					t.Logf("\t: actual   colour[%q] = %v", name, actColour)
					t.Logf("\t: expected colour[%q] = %v", name, expColour)
					t.Errorf("\t: bad result for %q", name)
				}
			}

			for name := range tc.exp {
				if _, ok := result[name]; !ok {
					t.Log(tc.IDStr())
					t.Errorf("\t: expected colour missing for %q", name)
				}
			}
		})
	}
}
