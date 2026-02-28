package colour

import (
	"fmt"
	"image/color" //nolint:misspell
	"math"
	"slices"
	"sort"

	"github.com/nickwells/english.mod/english"
)

// MaxColourProximity is the maximum effective value of the proximity value.
// Note that we add a small epsilon value to counteract inevitable numerical
// inaccuracies in the math package
var MaxColourProximity = (math.Sqrt(3) * math.MaxUint8) + 0.000001 //nolint:mnd

// FamilyColour represents a colour in a Family obtained by reference to its
// Euclidian distance (in the RGB colour cube) from some target colour.
//
//nolint:misspell
type FamilyColour struct {
	// dist is a measure of the distance between this colour and the target
	// colour. It is the square of the Euclidean distance from the colour in
	// the RGB colour cube. For display purposes the square root of the dist
	// should be used. Note that the Alpha value does not contribute to this
	// distance.
	dist int
	// Family is the colour Family in which this RGB colour has the
	// associated names
	Family Family
	// CNames is the list of names by which this RGB colour is known in the
	// associated Family
	CNames []string
	// Colour is the RGB colour
	Colour color.RGBA
}

// FullNames returns a single string giving all the possible names for this
// colour quoted and prefixed by the Family name
func (fc FamilyColour) FullNames() string {
	colourNames := []string{}
	for _, cn := range fc.CNames {
		colourNames = append(colourNames, fc.Family.Name()+":"+cn)
	}

	return english.JoinQuoted(colourNames, ", ", " or ")
}

// distSquared returns a distance metric for the two colours. This is the sum
// of the squares of the differences between each of the red, green and blue
// values. It is the square of the Euclidean distance in the RGB colour
// cube. Note that the Alpha value does not contribute to this distance.
func distSquared(target, c rgba) int {
	rd := int(c.R) - int(target.R)
	gd := int(c.G) - int(target.G)
	bd := int(c.B) - int(target.B)

	return rd*rd + gd*gd + bd*bd
}

// ClosestWithin returns those colours with Euclidean distance less than or
// equal to the given proximity from the given colour amongst the
// Families. The notion of 'closeness' is those colours having the smallest
// sum of squares of differences between the red, green and blue components.
//
// If proximity is equal to zero only exact matches will be returned. If
// proximity is greater than or equal to MaxColourProximity all colours will
// be returned. If it is less than zero then an error will be returned.
//
// If no families are given then the standard families are used.
func (fl Families) ClosestWithin(
	target color.RGBA, //nolint:misspell
	proximity float64,
) (
	[]FamilyColour, error,
) {
	if len(fl) == 0 {
		fl = standardFamilies
	}

	if err := fl.Check(); err != nil {
		return []FamilyColour{}, err
	}

	if proximity < 0 {
		return []FamilyColour{}, badProximityErr(proximity)
	}

	familyColours := fl.getSortedDists(target)

	// square the proximity so we don't have to take square roots
	return coloursWithin(familyColours, int(proximity*proximity)), nil
}

// coloursWithin returns all entries in pc with a Dist <= dist.
func coloursWithin(familyColours []FamilyColour, dist int) []FamilyColour {
	results := []FamilyColour{}
	if len(familyColours) == 0 {
		return results
	}

	lastFC := FamilyColour{}

	for _, fc := range familyColours {
		if fc.dist > dist {
			break
		}

		if lastFC.Family == "" { // we have the first of a new set
			lastFC = fc

			continue
		}

		if lastFC.Family == fc.Family &&
			lastFC.Colour == fc.Colour {
			// same Family and colour so add the names
			lastFC.CNames = append(lastFC.CNames, fc.CNames...)

			continue
		}

		sort.Strings(lastFC.CNames)
		results = append(results, lastFC)
		lastFC = fc
	}

	if lastFC.dist <= dist && // not strictly needed but for clarity and safety
		lastFC.Family != "" {
		sort.Strings(lastFC.CNames)
		results = append(results, lastFC)
	}

	return results
}

// ClosestN returns up to n colours closest to the given colour amongst the
// Families. The notion of 'closeness' is those colours having the smallest
// sum of squares of differences between the red, green and blue components
// of that colour and the target colour 'target'.
//
// The resulting slice may contain fewer than n entries if there are fewer
// than n distinct colours in the collection of Families.
//
// If no families are given then the standard families are used.
func (fl Families) ClosestN(target color.RGBA, n int) ( //nolint:misspell
	[]FamilyColour, error,
) {
	if len(fl) == 0 {
		fl = standardFamilies
	}

	if err := fl.Check(); err != nil {
		return []FamilyColour{}, err
	}

	if n < 0 {
		return []FamilyColour{},
			fmt.Errorf("bad colour count: %d - must be >= 0", n)
	}

	if n == 0 {
		return []FamilyColour{}, nil
	}

	familyColours := fl.getSortedDists(target)

	return nClosestColours(familyColours, n), nil
}

// nClosestColours returns the n entries in pc with the lowest distance from
// the target colour (see ClosestN).
func nClosestColours(familyColours []FamilyColour, n int) []FamilyColour {
	results := []FamilyColour{}
	if len(familyColours) == 0 || n == 0 {
		return results
	}

	lastFC := familyColours[0]

	for _, fc := range familyColours[1:] {
		if lastFC.Family == fc.Family &&
			lastFC.Colour == fc.Colour {
			lastFC.CNames = append(lastFC.CNames, fc.CNames...)

			continue
		}

		sort.Strings(lastFC.CNames)
		results = append(results, lastFC)
		lastFC = fc
	}

	sort.Strings(lastFC.CNames)
	results = append(results, lastFC)

	n = min(n, len(results))

	return results[:n]
}

// generateDists generates the proximities from the target colour for all the
// colours in all the Families. The Families should have already been checked
// for validity.
func (fl Families) generateDists(target rgba) []FamilyColour {
	results := []FamilyColour{}

	for _, f := range fl {
		fi, ok := f.info()
		if !ok {
			continue
		}

		for _, fc := range fi.colours {
			for name, c := range fc.cMap {
				results = append(results,
					FamilyColour{
						dist:   distSquared(target, c),
						Family: fc.f,
						CNames: []string{name},
						Colour: c,
					})
			}
		}
	}

	return results
}

// getSortedDists generates the proximities for all the colours in all the
// Families and then sorts them, according to FamilyColourCompare.
func (fl Families) getSortedDists(target rgba) []FamilyColour {
	familyColours := fl.generateDists(target)

	slices.SortFunc(familyColours, FamilyColourCompare)

	return familyColours
}
