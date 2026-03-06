package colour

import (
	"slices"
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestGenerateAltSpellings(t *testing.T) {
	//nolint:misspell
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
		{
			ID: testhelper.MkID("multiple changes US->UK"),
			s:  "gray aluminum harbor",
			expResults: []string{
				"gray aluminum harbor",
				"grey aluminium harbour",
			},
		},
		{
			ID: testhelper.MkID("multiple changes UK->US"),
			s:  "grey aluminium harbour",
			expResults: []string{
				"grey aluminium harbour",
				"gray aluminum harbor",
			},
		},
		{
			ID: testhelper.MkID("repeats UK->US"),
			s:  "grey grey grey",
			expResults: []string{
				"grey grey grey",
				"gray gray gray",
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

func TestIsAnAlias(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		name1   string
		name2   string
		expIs   bool
		expPref string
	}{
		{
			ID:      testhelper.MkID("same string"),
			name1:   "blah",
			name2:   "blah",
			expIs:   true,
			expPref: "blah",
		},
		{
			ID:      testhelper.MkID("different string"),
			name1:   "blah",
			name2:   "something else",
			expIs:   false,
			expPref: "",
		},
		{
			ID:      testhelper.MkID("grey/gray"),
			name1:   "grey",
			name2:   "gray",
			expIs:   true,
			expPref: "grey",
		},
		{
			ID:      testhelper.MkID("with apostrophes"),
			name1:   "xxx's",
			name2:   "xxxs",
			expIs:   true,
			expPref: "xxx's",
		},
		{
			ID:      testhelper.MkID("with spaces"),
			name1:   "xxx   s",
			name2:   "xxx-s",
			expIs:   true,
			expPref: "xxx   s",
		},
		{
			ID:      testhelper.MkID("with spaces and apostrophes"),
			name1:   "xxx   yyy's",
			name2:   "xxx-yyys",
			expIs:   true,
			expPref: "xxx   yyy's",
		},
		{
			ID:      testhelper.MkID("with grey/gray, spaces and apostrophes"),
			name1:   "grey xxx   yyy's",
			name2:   "gray-xxx-yyys",
			expIs:   true,
			expPref: "grey xxx   yyy's",
		},
		{
			ID:      testhelper.MkID("with just grey/gray differing"),
			name1:   "grey xxx   yyy's",
			name2:   "gray xxx   yyy's",
			expIs:   true,
			expPref: "grey xxx   yyy's",
		},
		{
			ID:      testhelper.MkID("with gray/grey, spaces and apostrophes"),
			name1:   "gray xxx   yyy's",
			name2:   "grey-xxx-yyys",
			expIs:   true,
			expPref: "gray xxx   yyy's",
		},
		{
			ID:      testhelper.MkID("nancys-grey-blushes 1"),
			name1:   "nancy's grey blushes",
			name2:   "nancys-grey-blushes",
			expIs:   true,
			expPref: "nancy's grey blushes",
		},
		{
			ID:      testhelper.MkID("nancy's gray blushes 2"),
			name1:   "nancy's grey blushes",
			name2:   "nancy's gray blushes",
			expIs:   true,
			expPref: "nancy's grey blushes",
		},
		{
			ID:      testhelper.MkID("nancys-gray-blushes 3"),
			name1:   "nancy's grey blushes",
			name2:   "nancys-gray-blushes",
			expIs:   true,
			expPref: "nancy's grey blushes",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actPref, actIs := IsAColourAlias(tc.name1, tc.name2)
			actPrefAlt, actIsAlt := IsAColourAlias(tc.name2, tc.name1)

			// first check that the order is irrelevant
			testhelper.DiffBool(t,
				tc.IDStr(), "is an alias-swapped",
				actIs, actIsAlt)
			testhelper.DiffString(t,
				tc.IDStr(), "pref-swapped",
				actPref, actPrefAlt)

			// now check for expected results
			testhelper.DiffBool(t,
				tc.IDStr(), "is an alias",
				actIs, tc.expIs)
			testhelper.DiffString(t,
				tc.IDStr(), "pref",
				actPref, tc.expPref)
		})
	}
}

func TestStripAliasesOf(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		s        string
		names    []string
		expNames []string
	}{
		{
			ID:       testhelper.MkID("empty slice"),
			s:        "hello",
			expNames: []string{"hello"},
		},
		{
			ID:       testhelper.MkID("target not in slice"),
			s:        "hello",
			names:    []string{"world", "flesh", "devil"},
			expNames: []string{"world", "flesh", "devil", "hello"},
		},
		{
			ID:       testhelper.MkID("target preferred"),
			s:        "hello world",
			names:    []string{"hello-world", "flesh", "devil"},
			expNames: []string{"flesh", "devil", "hello world"},
		},
		{
			ID:       testhelper.MkID("target not preferred"),
			s:        "hello-world",
			names:    []string{"hello world", "flesh", "devil"},
			expNames: []string{"flesh", "devil", "hello world"},
		},
		{
			ID: testhelper.MkID("multiple matches"),
			s:  "nancy's grey blushes",
			names: []string{
				"nancys-grey-blushes",
				"flesh",
				"nancys-gray-blushes",
				"devil",
				"nancy's gray blushes",
			},
			expNames: []string{"flesh", "devil", "nancy's grey blushes"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actNames := stripAliasesOf(tc.s, tc.names)
			if testhelper.DiffStringSlice(t,
				tc.IDStr(), "stripped names",
				actNames, tc.expNames) {
				t.Log("\t:   actual:", actNames)
				t.Log("\t: expected:", tc.expNames)
			}
		})
	}
}

func TestStripAliases(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		names    []string
		expNames []string
	}{
		{
			ID:       testhelper.MkID("no aliases"),
			names:    []string{"a", "b", "c"},
			expNames: []string{"a", "b", "c"},
		},
		{
			ID:       testhelper.MkID("one alias"),
			names:    []string{"a's", "b", "as", "c"},
			expNames: []string{"a's", "b", "c"},
		},
		{
			ID:       testhelper.MkID("several aliases"),
			names:    []string{"a's", "b  x", "as", "b-x", "c"},
			expNames: []string{"a's", "b  x", "c"},
		},
		{
			ID: testhelper.MkID("many, longer aliases"),
			names: []string{
				"encycolorpedia:dried and weathered bamboo",
				"encycolorpedia:dried-and-weathered-bamboo",
				"encycolorpedia:golden-gray bamboo",
				"encycolorpedia:golden-gray-bamboo",
				"encycolorpedia:golden-grey bamboo",
				"encycolorpedia:golden-grey-bamboo",
			},
			expNames: []string{
				"encycolorpedia:dried and weathered bamboo",
				"encycolorpedia:golden-grey bamboo",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			actNames := StripAliases(tc.names)

			if testhelper.DiffStringSlice(t,
				tc.IDStr(), "stripped list",
				actNames, tc.expNames) {
				t.Log("\t:   actual names:", actNames)
				t.Log("\t: expected names:", tc.expNames)
			}
		})
	}
}
