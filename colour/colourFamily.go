package colour

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/nickwells/english.mod/english"
)

// Family represents the collection of colour names
type Family uint8

const (
	AnyColours Family = iota
	WebColours
	CGAColours
	X11Colours
	HTMLColours
	PantoneColours
	maxFamily
)

var cFamMap = map[Family]map[string]color.RGBA{
	WebColours:     webColours,
	CGAColours:     cgaColours,
	X11Colours:     x11Colours,
	HTMLColours:    htmlColours,
	PantoneColours: pantoneColours,
}

var searchOrder = []Family{
	WebColours,
	CGAColours,
	HTMLColours,
	X11Colours,
	PantoneColours,
}

var (
	ErrBadFamily = errors.New("bad colour family")
	ErrBadColour = errors.New("bad colour name")
	ErrBadColor  = ErrBadColour
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
		for _, m := range cFamMap {
			for n := range m {
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
