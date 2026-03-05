package colour

import (
	"fmt"
	"image/color" //nolint:misspell
	"maps"
	"slices"

	"github.com/nickwells/english.mod/english"
)

// Families represents a collection of Family values
type Families []Family

// standardFamilies is the collection of families which are searched by default
var standardFamilies = Families{
	WebColours,
	CGAColours,
	HTMLColours,
	X11Colours,
	PantoneColours,
}

// familyColours returns a list of familyColours for each family in the list.
func (fl Families) familyColours() []familyToColourMap {
	fToCNM := map[Family]colourNameToRGBA{
		WebColours:            webColours,
		CGAColours:            cgaColours,
		HTMLColours:           htmlColours,
		X11Colours:            x11Colours,
		PantoneColours:        pantoneColours,
		FarrowAndBallColours:  farrowAndBallColours,
		CrayolaColours:        crayolaColours,
		XKCDColours:           xkcdColours,
		EncycolorpediaColours: encycolorpediaColours,
	}

	fcs := []familyToColourMap{}

	for _, f := range fl {
		cnm, ok := fToCNM[f]
		if !ok {
			panic(badFamilyErr(f))
		}

		fcs = append(fcs, familyToColourMap{f, cnm})
	}

	return fcs
}

// String returns the list of family IDs as a string
func (fl Families) String() string {
	fNames := []string{}
	for _, f := range fl {
		fNames = append(fNames, string(f))
	}

	return english.Join(fNames, ", ", " and ")
}

// ByName returns the list of family names (the Family value cast to
// lowercase) as a string, see [Family.Name].
func (fl Families) ByName() string {
	fNames := []string{}
	for _, f := range fl {
		fNames = append(fNames, f.Name())
	}

	return english.Join(fNames, ", ", " and ")
}

// Check checks the Family list for duplicate entries and invalid Family
// entries. It returns a non-nil error if any problems are found.
func (fl Families) Check() error {
	problems := []string{}
	dupMap := map[Family][]int{}

	for i, f := range fl {
		dupMap[f] = append(dupMap[f], i)
		if !f.IsValid() {
			problems = append(problems,
				fmt.Sprintf("%q is not a valid Family (at position %d)",
					f, i))
		}
	}

	for f, idxs := range dupMap {
		if len(idxs) > 1 {
			problems = append(problems,
				fmt.Sprintf(
					"%q appears %d times"+
						" in the Families list, at positions: %v",
					f, len(idxs), idxs))
		}
	}

	if len(problems) > 0 {
		slices.Sort(problems)

		return fmt.Errorf("%d %s found: %s",
			len(problems),
			english.Plural("problem", len(problems)),
			english.Join(problems, ", ", " and "))
	}

	return nil
}

// AllColours returns a slice containing all the colours in each of the given
// Family elements in the list. The returned value has distinct entries (no
// colour appears twice) but in a random order. A non-nil error is returned
// if any Familly in the list is not recognised.
func (fl Families) AllColours() ([]color.RGBA, error) { //nolint:misspell
	var colours []rgba

	colourMap := map[rgba]bool{}

	if len(fl) == 0 {
		fl = standardFamilies
	}

	for _, f := range fl {
		fi, ok := f.info()
		if !ok {
			return colours, badFamilyErr(f)
		}

		for _, m := range fi.colours {
			for _, c := range m.cMap {
				colourMap[c] = true
			}
		}
	}

	return slices.Collect(maps.Keys(colourMap)), nil
}

// AllColourNames returns a slice containing all the names of the colours in
// each of the given Family elements in the list. The returned value has
// distinct entries (no colour name appears twice) but in a random order. A
// non-nil error is returned if any Familly in the list is not recognised.
func (fl Families) AllColourNames() ([]string, error) { //nolint:misspell
	var colours []string

	colourMap := map[string]bool{}

	if len(fl) == 0 {
		fl = standardFamilies
	}

	for _, f := range fl {
		fi, ok := f.info()
		if !ok {
			return colours, badFamilyErr(f)
		}

		for _, m := range fi.colours {
			for name := range m.cMap {
				colourMap[name] = true
			}
		}
	}

	return slices.Collect(maps.Keys(colourMap)), nil
}
