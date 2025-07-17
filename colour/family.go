package colour

import (
	"fmt"
	"image/color" //nolint:misspell
	"strings"

	"github.com/nickwells/english.mod/english"
)

// rgba is a local type alias for color.RGBA so we don't have to have
// nolint comments littering the code everywhere
type rgba = color.RGBA //nolint:misspell

//nolint:misspell
// To add a new Family of colours you will need to:
//
// 1. add a new Family to the set of ...Colours consts
// 2. add a new alias to the ...Color = ...Colour consts
// 3. make a new entry in the allFamilies map
// 4. change the Families.familyColours() func to construct the correct
//    familyColours struct for the new Family value
// 5. if you want it to be one of the StandardFamilies you will need to
//    add it to the StandardFamilies slice

// Family identifies a collection of colour names
type Family string

// These represent the various families of colour names.
const (
	StandardColours      Family = "Standard"
	WebColours           Family = "Web"
	CGAColours           Family = "CGA"
	X11Colours           Family = "X11"
	HTMLColours          Family = "HTML"
	PantoneColours       Family = "Pantone"
	FarrowAndBallColours Family = "FarrowAndBall"
	CrayolaColours       Family = "Crayola"
	XKCDColours          Family = "xkcd"
)

// Some aliases for people who use Merriam-Webster rather than the OED
const (
	StandardColors      = StandardColours
	WebColors           = WebColours
	CGAColors           = CGAColours
	X11Colors           = X11Colours
	HTMLColors          = HTMLColours
	PantoneColors       = PantoneColours
	FarrowAndBallColors = FarrowAndBallColours
	CrayolaColors       = CrayolaColours
	XKCDColors          = XKCDColours
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

// colourNameToRGBA is the type of the structures mapping the text names to
// the RGBA values
type colourNameToRGBA map[string]rgba

// familyColours collects the Family with a map of colour names
type familyColours struct {
	f    Family
	cMap colourNameToRGBA
}

// familyInfo collects information about a colour family
type familyInfo struct {
	id          Family
	name        string
	description string
	colours     []familyColours
}

// allFamilies maps from a Family.Name to the associated familyInfo. There is
// one entry for each valid Family.
var allFamilies = map[string]familyInfo{
	StandardColours.Name(): {
		id:   StandardColours,
		name: StandardColours.Name(),
		description: "any of the 'standard' colour families: " +
			standardFamilies.TextByName(),
		colours: standardFamilies.familyColours(),
	},
	WebColours.Name(): {
		id:          WebColours,
		name:        WebColours.Name(),
		description: "colour names as defined in the HTML 4.01 specification",
		colours:     Families{WebColours}.familyColours(),
	},
	CGAColours.Name(): {
		id:          CGAColours,
		name:        CGAColours.Name(),
		description: "CGA colours - the Web colours but with different names",
		colours:     Families{CGAColours}.familyColours(),
	},
	HTMLColours.Name(): {
		id:          HTMLColours,
		name:        HTMLColours.Name(),
		description: "HTML colours - colours supported by all browsers",
		colours:     Families{HTMLColours}.familyColours(),
	},
	X11Colours.Name(): {
		id:          X11Colours,
		name:        X11Colours.Name(),
		description: "colour names from the X11 rgb.txt file",
		colours:     Families{X11Colours}.familyColours(),
	},
	PantoneColours.Name(): {
		id:   PantoneColours,
		name: PantoneColours.Name(),
		description: "colour names from the Pantone" +
			" Fashion Home + Interiors range as of 2023/Nov/11",
		colours: Families{PantoneColours}.familyColours(),
	},
	FarrowAndBallColours.Name(): {
		id:          FarrowAndBallColours,
		name:        FarrowAndBallColours.Name(),
		description: "Farrow And Ball paint colours",
		colours: []familyColours{
			{FarrowAndBallColours, farrowAndBallColours},
		},
	},
	CrayolaColours.Name(): {
		id:          CrayolaColours,
		name:        CrayolaColours.Name(),
		description: "colour names from the Crayola crayon range",
		colours:     Families{CrayolaColours}.familyColours(),
	},
	XKCDColours.Name(): {
		id:          XKCDColours,
		name:        XKCDColours.Name(),
		description: "colour names as surveyed by xkcd, some names are NSFW",
		colours:     Families{XKCDColours}.familyColours(),
	},
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

// info returns the familyInfo details for the given family
func (f Family) info() (familyInfo, bool) {
	key := f.Name()
	fi, ok := allFamilies[key]

	return fi, ok
}

// AllowedFamilies returns a map from colour family names to their associated
// descriptions. This is suitable for initialising a
// [github.com/nickwells/param.mod/v6/psetter.AllowedVals] map.
func AllowedFamilies() map[string]string {
	af := make(map[string]string)

	for name, fi := range allFamilies {
		colourCount := 0
		for _, fc := range fi.colours {
			colourCount += len(fc.cMap)
		}

		af[name] = fi.description + fmt.Sprintf(" (%d colours)", colourCount)
	}

	return af
}

// Text returns the list of family IDs as a string
func (fl Families) Text() string {
	fNames := []string{}
	for _, f := range fl {
		fNames = append(fNames, string(f))
	}

	return english.Join(fNames, ", ", " and ")
}

// TextByName returns the list of family names (the Family value cast to
// lowercase) as a string, see [Family.Name].
func (fl Families) TextByName() string {
	fNames := []string{}
	for _, f := range fl {
		fNames = append(fNames, f.Name())
	}

	return english.Join(fNames, ", ", " and ")
}

// Name returns the 'name'-form of the Family - this is the value converted
// to a string and mapped to lower case.
func (f Family) Name() string {
	return strings.ToLower(string(f))
}

// IsValid returns true if f is a recognised colour Family, false otherwise.
func (f Family) IsValid() bool {
	key := f.Name()
	_, ok := allFamilies[key]

	return ok
}

// ColourNames returns a slice containing the names of colours in the given
// Family. The returned value has distinct entries (no name appears twice)
// but in a random order.
func (f Family) ColourNames() ([]string, error) {
	var names []string

	fi, ok := f.info()
	if !ok {
		return names, badFamilyErr(f)
	}

	nameMap := map[string]bool{}

	for _, m := range fi.colours {
		for cn := range m.cMap {
			nameMap[cn] = true
		}
	}

	for n := range nameMap {
		names = append(names, n)
	}

	return names, nil
}

// AllColours returns a slice containing all the colours in the given
// Family. The returned value has distinct entries (no colour appears twice)
// but in a random order.
func (f Family) AllColours() ([]color.RGBA, error) { //nolint:misspell
	var colours []rgba

	fi, ok := f.info()
	if !ok {
		return colours, badFamilyErr(f)
	}

	colourMap := map[rgba]bool{}

	for _, m := range fi.colours {
		for _, c := range m.cMap {
			colourMap[c] = true
		}
	}

	for c := range colourMap {
		colours = append(colours, c)
	}

	return colours, nil
}

// ColorNames - see [Family.ColourNames]
//
// This is an alias for people who follow Merriam-Webster rather than the
// OED. Note that there is a (very) small performance disadvantage from using
// this aliased form which the compiler might optimise away.
func (f Family) ColorNames() ([]string, error) {
	return f.ColourNames()
}

// Colour returns the RGBA colour of the given name in the given Family. If a
// Family has multiple colour maps then the first matching colour is
// returned. A non-nil error is returned if the Family is not recognised or
// if the colour name is not found.
func (f Family) Colour(cName string) (color.RGBA, error) { //nolint:misspell
	fi, ok := f.info()
	if !ok {
		return rgba{}, badFamilyErr(f)
	}

	for _, m := range fi.colours {
		if c, ok := m.cMap[cName]; ok {
			return c, nil
		}
	}

	return rgba{}, badColourErr(cName)
}

// Color - see [Family.Colour]
//
// This is an alias for people who follow Merriam-Webster rather than the
// OED. Note that there is a (very) small performance disadvantage from using
// this aliased form which the compiler might optimise away.
//
//nolint:misspell
func (f Family) Color(cName string) (color.RGBA, error) {
	return f.Colour(cName)
}
