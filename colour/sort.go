package colour

import (
	"image/color" //nolint:misspell
	"strings"
)

// RGBACompare returns an integer less than zero if a is less than b, equal
// to zero if a equals b and greater than zero otherwise. The ordering is by
// the red, green, blue and alpha values of the colours.
func RGBACompare(a, b color.RGBA) int { //nolint:misspell
	if d := int(a.R) - int(b.R); d != 0 {
		return d
	}

	if d := int(a.B) - int(b.B); d != 0 {
		return d
	}

	if d := int(a.G) - int(b.G); d != 0 {
		return d
	}

	if d := int(a.A) - int(b.A); d != 0 {
		return d
	}

	return 0
}

// FamilyColourCompare returns an integer less than zero if a is less than b,
// equal to zero if a equals b and greater than zero otherwise. The ordering
// is by distance (from the target colour used to generate the value), Family
// name and colour.
func FamilyColourCompare(a, b FamilyColour) int {
	if d := a.dist - b.dist; d != 0 {
		return d
	}

	if d := strings.Compare(string(a.Family), string(b.Family)); d != 0 {
		return d
	}

	return RGBACompare(a.Colour, b.Colour)
}

// NamedColourCompare returns an integer less than zero if a is less than b,
// equal to zero if a equals b and greater than zero otherwise. The ordering
// is by colour and name.
func NamedColourCompare(a, b NamedColour) int {
	if d := RGBACompare(a.Colour(), b.Colour()); d != 0 {
		return d
	}

	return strings.Compare(a.Name(), b.Name())
}
