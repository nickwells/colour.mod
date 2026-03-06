//nolint:misspell
package colour

import (
	"image/color"
	"regexp"
)

// MakeColors - see [MakeColours]
func MakeColors(count int) ([]color.RGBA, error) {
	return MakeColours(count)
}

// MakeColorsBetween - see [MakeColoursBetween]
func MakeColorsBetween(count int, lower, upper color.RGBA) (
	[]color.RGBA,
	error,
) {
	return MakeColoursBetween(count, lower, upper)
}

// NamedColorAllowedValues - see [NamedColourAllowedValues]
func NamedColorAllowedValues(fl Families) string {
	return NamedColourAllowedValues(fl)
}

// Color - see [NamedColour.Colour]
func (nc NamedColour) Color() color.RGBA {
	return nc.Colour()
}

// MakeNamedColor - see [MakeNamedColour]
func MakeNamedColor(name string, c color.RGBA) NamedColour {
	return MakeNamedColour(name, c)
}

// MakeNamedColorFromFamilyColor - see [MakeNamedColourFromFamilyColour]
func MakeNamedColorFromFamilyColor(fc FamilyColour) NamedColour {
	return MakeNamedColourFromFamilyColour(fc)
}

// MakeDistinctNamedColorsFromFamilyColor see
// [MakeDistinctNamedColoursFromFamilyColour]
func MakeDistinctNamedColorsFromFamilyColor(fc FamilyColour) []NamedColour {
	return MakeDistinctNamedColoursFromFamilyColour(fc)
}

// ParseNamedColor - see [ParseNamedColour]
func ParseNamedColor(fl Families, s string) (NamedColour, error) {
	return ParseNamedColour(fl, s)
}

// ColorsMatchingByRegexp - see [ColoursMatchingByRegexp]
func ColorsMatchingByRegexp(fl Families, re *regexp.Regexp) (
	[]NamedColour, error,
) {
	return ColoursMatchingByRegexp(fl, re)
}

// ColorsMatchingByFunc - see [ColoursMatchingByFunc]
func ColorsMatchingByFunc(fl Families, matchFunc MatchingFunc) (
	[]NamedColour, error,
) {
	return ColoursMatchingByFunc(fl, matchFunc)
}

// IsAPotentialColorString - see [IsAPotentialColourString]
func IsAPotentialColorString(s string) bool {
	return IsAPotentialColourString(s)
}

// ParseColorDefinition - see [ParseColourDefinition]
func ParseColorDefinition(s string) (color.RGBA, error) {
	return ParseColourDefinition(s)
}

// ParseColorPart - see [ParseColourPart]
func ParseColorPart(val, partName string) (uint8, error) {
	return ParseColourPart(val, partName)
}

// AllColors - see [AllColours]
func (fl Families) AllColors() ([]color.RGBA, error) {
	return fl.AllColours()
}

// AllColorNames - see [AllColourNames]
func (fl Families) AllColorNames() ([]string, error) {
	return fl.AllColourNames()
}

// CGAColor - see [CGAColour]
func CGAColor(i int) (color.RGBA, error) {
	return CGAColour(i)
}

// ContrastColorful - see [ContrastColourful]
func ContrastColorful(c color.RGBA) color.RGBA {
	return ContrastColourful(c)
}

// AllColors - see [Family.AllColours]
func (f Family) AllColors() ([]color.RGBA, error) {
	return f.AllColours()
}

// ColorNames - see [Family.ColourNames]
func (f Family) ColorNames() ([]string, error) {
	return f.ColourNames()
}

// Color - see [Family.Colour]
func (f Family) Color(cName string) (color.RGBA, error) {
	return f.Colour(cName)
}

// ColorNameCount - see [Family.ColourNameCount]
func (f Family) ColorNameCount() (int, error) {
	return f.ColourNameCount()
}

// DistinctColorCount - see [Family.DistinctColourCount]
func (f Family) DistinctColorCount() (int, error) {
	return f.DistinctColourCount()
}

// ToGray - see [ToGrey]
func ToGray(c color.RGBA) color.RGBA {
	return ToGrey(c)
}

// ToGrayEqual - see [ToGreyEqual]
func ToGrayEqual(c color.RGBA) color.RGBA {
	return ToGreyEqual(c)
}

// ToGrayBT709 - see [ToGreyBT709]
func ToGrayBT709(c color.RGBA) color.RGBA {
	return ToGreyBT709(c)
}

// ToGrayBT2100 - see [ToGreyBT2100]
func ToGrayBT2100(c color.RGBA) color.RGBA {
	return ToGreyBT2100(c)
}

// ToGrayCustom - see [ToGreyCustom]
func ToGrayCustom(c color.RGBA,
	wtRed, wtGreen, wtBlue float64,
) (color.RGBA, error) {
	return ToGreyCustom(c, wtRed, wtGreen, wtBlue)
}

// MakeGray - see [MakeGrey]
func MakeGray(g uint8) color.RGBA {
	return MakeGrey(g)
}

// IsAColorAlias - see [IsAColourAlias]
func IsAColorAlias(s1, s2 string) (string, bool) {
	return IsAColourAlias(s1, s2)
}
