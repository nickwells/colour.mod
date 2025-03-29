package colour

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/nickwells/english.mod/english"
)

// Family labels a collection of colour names
type Family uint8

// These represent the various families of colour names.
const (
	AnyColours Family = iota // means: use the standard search list
	WebColours
	CGAColours
	X11Colours
	HTMLColours
	PantoneColours
	FarrowAndBallColours // not in the standard search list
	CrayolaColours       // not in the standard search list
	maxFamily
)

// Some aliases for people who use Merriam-Webster rather than the OED
const (
	AnyColors           = AnyColours
	WebColors           = WebColours
	CGAColors           = CGAColours
	X11Colors           = X11Colours
	HTMLColors          = HTMLColours
	PantoneColors       = PantoneColours
	FarrowAndBallColors = FarrowAndBallColours
	CrayolaColors       = CrayolaColours
)

// the colour family map is the way to access the various colour maps
var cFamMap = map[Family]map[string]color.RGBA{
	WebColours:           webColours,
	CGAColours:           cgaColours,
	X11Colours:           x11Colours,
	HTMLColours:          htmlColours,
	PantoneColours:       pantoneColours,
	FarrowAndBallColours: farrowAndBallColours,
	CrayolaColours:       crayolaColours,
}

// The searchOrder provides the set of colours to search when the AnyColours
// family has been selected. Note that not all the available colour maps are
// present.
//
// excluded colour families:
//   - FarrowAndBallColours
//   - CrayolaColours
var searchOrder = []Family{
	WebColours,
	CGAColours,
	HTMLColours,
	X11Colours,
	PantoneColours,
}

var (
	// ErrBadFamily is the standard error returned when a bad colour family
	// value is passed
	ErrBadFamily = errors.New("bad colour family")
	// ErrBadColour is the standard error returned when a bad colour name is
	// passsed
	ErrBadColour = errors.New("bad colour name")
	// ErrBadColor is an alternative value for people who use Merriam-Webster
	// rather than the OED
	ErrBadColor = ErrBadColour
)

// String returns a string representing the given Family or a value
// indicating that the Family is not recognised.
func (f Family) String() string {
	switch f {
	case AnyColours:
		return "Any"
	case WebColours:
		return "Web"
	case CGAColours:
		return "CGA"
	case X11Colours:
		return "X11"
	case HTMLColours:
		return "HTML"
	case PantoneColours:
		return "Pantone"
	case FarrowAndBallColours:
		return "FarrowAndBall"
	case CrayolaColours:
		return "Crayola"
	}

	return fmt.Sprintf("BadFamily:%d", f)
}

// Literal returns a string equal to the given Family name const symbol. It
// will panic if the value is not recognised.
func (f Family) Literal() string {
	switch f {
	case AnyColours:
		return "AnyColours"
	case WebColours:
		return "WebColours"
	case CGAColours:
		return "CGAColours"
	case X11Colours:
		return "X11Colours"
	case HTMLColours:
		return "HTMLColours"
	case PantoneColours:
		return "PantoneColours"
	case FarrowAndBallColours:
		return "FarrowAndBallColours"
	case CrayolaColours:
		return "CrayolaColours"
	}

	panic(fmt.Errorf("BadFamily:%d", f))
}

// familyList returns the list of families as a string
func familyList(fl []Family) string {
	fNames := []string{}
	for _, f := range fl {
		fNames = append(fNames, f.String())
	}

	return english.Join(fNames, ", ", " and ")
}

// IsValid returns true if f is a valid colour Family, false otherwise
func (f Family) IsValid() bool {
	return f < maxFamily
}

// ColourNames returns the names of colours in the given Family.
func (f Family) ColourNames() []string {
	var names []string

	nameMap := map[string]bool{}

	if f == AnyColours {
		for _, acFam := range searchOrder {
			for n := range cFamMap[acFam] {
				nameMap[n] = true
			}
		}
	} else {
		if m, ok := cFamMap[f]; ok {
			for n := range m {
				nameMap[n] = true
			}
		}
	}

	for n := range nameMap {
		names = append(names, n)
	}

	return names
}

// ColorNames - see ColourNames.
//
// This is an alias for people who follow Merriam-Webster rather than the
// OED. Note that there is a (very) small performance advantage from using
// this aliased form
func (f Family) ColorNames() []string {
	return f.ColourNames()
}

// Colour returns the RGBA colour of the given name in the given Family.
func (f Family) Colour(cName string) (color.RGBA, error) {
	if f == AnyColours {
		for _, cf := range searchOrder {
			cMap := cFamMap[cf]
			if cVal, ok := cMap[cName]; ok {
				return cVal, nil
			}
		}

		return color.RGBA{}, ErrBadColour
	}

	cMap, ok := cFamMap[f]
	if !ok {
		return color.RGBA{}, ErrBadFamily
	}

	if cVal, ok := cMap[cName]; ok {
		return cVal, nil
	}

	return color.RGBA{}, ErrBadColour
}

// Color - see Colour
//
// This is an alias for people who follow Merriam-Webster rather than the
// OED. Note that there is a (very) small performance advantage from using
// this aliased form
func (f Family) Color(cName string) (color.RGBA, error) {
	return f.Colour(cName)
}
