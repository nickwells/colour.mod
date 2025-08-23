package colour

import (
	"fmt"
	"image/color" //nolint:misspell
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
func (fl Families) familyColours() []familyColours {
	fcs := []familyColours{}

	for _, f := range fl {
		switch f {
		case WebColours:
			fcs = append(fcs, familyColours{f, webColours})
		case CGAColours:
			fcs = append(fcs, familyColours{f, cgaColours})
		case HTMLColours:
			fcs = append(fcs, familyColours{f, htmlColours})
		case X11Colours:
			fcs = append(fcs, familyColours{f, x11Colours})
		case PantoneColours:
			fcs = append(fcs, familyColours{f, pantoneColours})
		case FarrowAndBallColours:
			fcs = append(fcs, familyColours{f, farrowAndBallColours})
		case CrayolaColours:
			fcs = append(fcs, familyColours{f, crayolaColours})
		case XKCDColours:
			fcs = append(fcs, familyColours{f, xkcdColours})
		}
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

	for c := range colourMap {
		colours = append(colours, c)
	}

	return colours, nil
}
