package colour

import (
	"errors"
	"fmt"
	"image/color" //nolint:misspell
	"maps"
	"regexp"
	"slices"
	"strings"

	"github.com/nickwells/strdist.mod/v2/strdist"
)

const (
	familyColourAllowedValue = "a family name, a colon (:) and a colour name"
	colourNameAllowedValue   = "a colour name"
)

// NamedColourAllowedValues returns a string describing the values which can
// be parsed into a named colour.
func NamedColourAllowedValues(fl Families) string {
	var fDesc string

	const stdFamilies = " in the standard colour-name families"

	switch len(fl) {
	case 0:
		fDesc = stdFamilies
	case 1:
		switch fl[0] {
		case StandardColours:
			fDesc = stdFamilies
		default:
			fDesc = " in the " + string(fl[0]) + " colour-name family"
		}
	default:
		fDesc = " in one of the colour-name families: " + fl.String()
	}

	return "Either" +
		"\n" +
		colourNameAllowedValue + fDesc +
		"\n" +
		"or " + familyColourAllowedValue +
		"\n" +
		"or " + RGBAllowedValues()
}

// NamedColour records a colour and the associated name
type NamedColour struct {
	name   string
	colour color.RGBA //nolint:misspell
}

// Name returns the colour name
func (nc NamedColour) Name() string {
	return nc.name
}

// Colour returns the colour
func (nc NamedColour) Colour() color.RGBA { //nolint:misspell
	return nc.colour
}

// MakeNamedColour returns a NamedColour
func MakeNamedColour(name string, c color.RGBA) NamedColour { //nolint:misspell
	return NamedColour{
		name:   name,
		colour: c,
	}
}

// MakeNamedColourFromFamilyColour generates a NamedColour with its colour
// set to the FamilyColour's colour and it's name generated from the Family
// and CNames fields
func MakeNamedColourFromFamilyColour(fc FamilyColour) NamedColour {
	return MakeNamedColour(fc.FullNames(), fc.Colour)
}

// MakeDistinctNamedColoursFromFamilyColour generates a NamedColour with its
// colour set to the FamilyColour's colour and it's name generated from the
// Family and CNames fields
func MakeDistinctNamedColoursFromFamilyColour(fc FamilyColour) []NamedColour {
	namedColours := []NamedColour{}

	for _, cn := range fc.CNames {
		namedColours = append(namedColours,
			MakeNamedColour(fc.Family.Name()+":"+cn, fc.Colour))
	}

	return namedColours
}

// ParseNamedColour creates a NamedColour from the given string. It will
// first check that the string matches an RGB specification and if so it will
// use that. Otherwise if the string matches a family name and colour name
// (separated by a ':') then it will return the named colour from the given
// family. Lastly it will find the named colour in the given families list.
func ParseNamedColour(fl Families, s string) (NamedColour, error) {
	nc := NamedColour{name: s}

	var err error

	if IsAPotentialColourString(s) {
		nc.colour, err = ParseColourDefinition(s)

		return nc, err
	}

	if familyName, colourName, found := strings.Cut(s, ":"); found {
		nc.colour, err = getColourByFamilyAndColourName(familyName, colourName)

		return nc, err
	}

	nc.colour, err = getColourByColourName(fl, s)

	return nc, err
}

// getColourByFamilyAndColourName gets the colour value from the family and
// colour names. It returns a non-nil error if the value can not be set.
func getColourByFamilyAndColourName(fName, cName string,
) (
	c color.RGBA, err error, //nolint:misspell
) {
	origFName := fName
	fName = strings.TrimSpace(fName)
	fName = strings.ToLower(fName)

	cName = strings.TrimSpace(cName)
	cName = strings.ToLower(cName)

	f := Family(fName)
	if !f.IsValid() {
		return c, fmt.Errorf("bad colour family name: %q%s",
			origFName,
			strdist.SuggestionString(
				strdist.SuggestedVals(fName,
					slices.Collect(maps.Keys(AllowedFamilies())),
				)))
	}

	c, err = f.Colour(cName)
	if err != nil {
		altNames := ""
		if cNames, err := f.ColourNames(); err == nil {
			altNames = strdist.SuggestionString(
				strdist.SuggestedVals(cName, cNames))
		}

		return c, fmt.Errorf("%s%s",
			err, altNames)
	}

	return c, nil
}

// getColourByColourName sets the RGB from the colour name. It
// returns a non-nil error if the value can not be set.
func getColourByColourName(fl Families, cName string) (
	c color.RGBA, err error, //nolint:misspell
) {
	cName = strings.TrimSpace(cName)
	cName = strings.ToLower(cName)

	if len(fl) == 0 {
		c, err = StandardColours.Colour(cName)
	} else {
		for _, f := range fl {
			c, err = f.Colour(cName)
			if err == nil {
				break
			}
		}
	}

	if err != nil {
		if errors.Is(err, errors.New(BadColourName)) {
			allNames, colourErr := fl.AllColourNames()
			if colourErr == nil {
				return c, fmt.Errorf("bad colour name: %q%s",
					cName, strdist.SuggestionString(
						strdist.SuggestedVals(cName, allNames)))
			}
		}

		return c, err
	}

	return c, nil
}

// ColoursMatchingByRegexp returns a set of NamedColours where the colour
// names in the families match the regular expression. The resulting colour
// names will include the family from which they were matched, separated by a
// colon (:).
func ColoursMatchingByRegexp(fl Families, re *regexp.Regexp) (
	[]NamedColour, error,
) {
	if re == nil {
		return nil, errors.New("no regular expression was provided")
	}

	return ColoursMatchingByFunc(fl,
		func(name string, _ rgba) bool {
			return re.MatchString(name)
		})
}

// MatchingFunc is the type of a function to be supplied to the
// ColoursMatchingByFunc function. It takes the colour name and the RGBA
// colour value and it should return true if the colour should be in the
// results set.
type MatchingFunc func(string, color.RGBA) bool //nolint:misspell

// ColoursMatchingByFunc returns a set of NamedColours where the matchFunc
// when passed the colour name and colour returns true. The resulting colour
// names will include the family from which they were matched, separated by a
// colon (:). The results will be sorted according to the NamedColourCompare
// function.
func ColoursMatchingByFunc(fl Families, matchFunc MatchingFunc) (
	[]NamedColour, error,
) {
	if matchFunc == nil {
		return nil, errors.New("no match function was provided")
	}

	ncs := []NamedColour{}

	if len(fl) == 0 {
		fl = standardFamilies
	}

	for _, f := range fl {
		fi, ok := f.info()
		if !ok {
			return ncs, badFamilyErr(f)
		}

		for _, fc := range fi.colours {
			for cName, c := range fc.cMap {
				if matchFunc(cName, c) {
					nc := NamedColour{
						name:   fc.f.Name() + ":" + cName,
						colour: c,
					}
					ncs = append(ncs, nc)
				}
			}
		}
	}

	slices.SortFunc(ncs, NamedColourCompare)

	return ncs, nil
}
