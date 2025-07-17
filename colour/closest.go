package colour

import (
	"errors"
	"image/color" //nolint:misspell
	"math"
	"sort"
)

var maxProximity = math.Sqrt(3) * math.MaxUint8 //nolint:mnd

// FamilyColour represents a colour in a Family
//
//nolint:misspell
type FamilyColour struct {
	// dist is a measure of the distance between this colour and the search
	// colour. It is the square of the Euclidean distance from the colour in
	// the RGB colour cube. For display purposes the square root of the dist
	// should be used.
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

// dist returns the distance metric for the two colours. This is the sum of
// the square of the differences between each of the red, green and blue
// values. It is the square of the Euclidean distance in the RGB colour cube.
func dist(sc, c rgba) int {
	rd, gd, bd := int(c.R)-int(sc.R), int(c.G)-int(sc.G), int(c.B)-int(sc.B)

	return rd*rd + gd*gd + bd*bd
}

// ClosestWithin returns those colours with Euclidean distance less than
// proximity from the given colour amongst the Families. The notion of
// 'closeness' is those colours having the smallest sum of squares of
// differences between the red, green and blue components.
//
// If prox is set to zero only exact matches will be returned. If it is set
// to 442 (255 times root 3) or more then every colour will be returned.
func (fl Families) ClosestWithin(
	sc color.RGBA, //nolint:misspell
	proximity float64,
) (
	[]FamilyColour, error,
) {
	if err := fl.check(); err != nil {
		return []FamilyColour{}, err
	}

	if proximity < 0 {
		return []FamilyColour{}, badProximityErr(proximity)
	}

	if proximity >= maxProximity {
		return []FamilyColour{}, badProximityErr(proximity)
	}

	cds := fl.getSortedDists(sc)

	// square the proximity so we don't have to take square roots
	return coloursWithin(cds, int(proximity*proximity)), nil
}

// coloursWithin returns all entries in pc with a Dist <= dist.
func coloursWithin(cds []FamilyColour, dist int) []FamilyColour {
	results := []FamilyColour{}
	lastCD := FamilyColour{}

	for _, cd := range cds {
		if cd.dist > dist {
			break
		}

		if lastCD.Family == "" { // we have the first of a new set
			lastCD = cd

			continue
		}

		if lastCD.Family == cd.Family &&
			lastCD.Colour == cd.Colour {
			// same Family and colour so add the names
			lastCD.CNames = append(lastCD.CNames, cd.CNames...)

			continue
		}

		results = append(results, lastCD)
		lastCD = cd
	}

	if lastCD.dist <= dist && // not strictly needed but for clarity and safety
		lastCD.Family != "" {
		results = append(results, lastCD)
	}

	return results
}

// ClosestN returns up to n colours closest to the given colour amongst the
// Families. The notion of 'closeness' is those colours having the smallest
// sum of squares of differences between the red, green and blue components
// of that colour and the search colour 'sc'.
//
// The resulting slice may contain fewer than n entries if there are fewer
// than n distinct colours in the collection of Families.
func (fl Families) ClosestN(sc color.RGBA, n int) ( //nolint:misspell
	[]FamilyColour, error,
) {
	if err := fl.check(); err != nil {
		return []FamilyColour{}, err
	}

	if n < 0 {
		return []FamilyColour{}, errors.New("bad colour count, must be >= 0")
	}

	if n == 0 {
		return []FamilyColour{}, nil
	}

	cds := fl.getSortedDists(sc)

	return nClosestColours(cds, n), nil
}

// nClosestColours returns the n entries in pc with the lowest distance from
// the search colour (see Closest).
func nClosestColours(cds []FamilyColour, n int) []FamilyColour {
	results := []FamilyColour{}
	lastCD := FamilyColour{}

	for _, cd := range cds {
		if len(results) >= n {
			break
		}

		if lastCD.Family == "" {
			lastCD = cd

			continue
		}

		if lastCD.Family == cd.Family &&
			lastCD.Colour == cd.Colour {
			lastCD.CNames = append(lastCD.CNames, cd.CNames...)

			continue
		}

		results = append(results, lastCD)
		lastCD = FamilyColour{}
	}

	if len(results) < n && // not strictly needed but for clarity and safety
		lastCD.Family != "" {
		results = append(results, lastCD)
	}

	return results
}

// check checks that each of the members of Families is a valid Family.
func (fl Families) check() error {
	for _, f := range fl {
		if !f.IsValid() {
			return badFamilyErr(f)
		}
	}

	return nil
}

// generateDists generates the proximities for all the colours in all
// the Families. The Families should have already been checked for validity.
func (fl Families) generateDists(sc rgba) []FamilyColour {
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
						dist:   dist(sc, c),
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
// Families and then sorts them, firstly by distance, then by Family in the
// order given in the Families slice and lastly by colour.
func (fl Families) getSortedDists(sc rgba) []FamilyColour {
	cds := fl.generateDists(sc)

	familyIdx := map[Family]int{}
	for i, f := range fl {
		familyIdx[f] = i
	}

	sort.Slice(cds, func(i, j int) bool {
		if cds[i].dist != cds[j].dist {
			return cds[i].dist < cds[j].dist
		}

		iIdx, jIdx := familyIdx[cds[i].Family], familyIdx[cds[j].Family]
		if iIdx != jIdx {
			return iIdx < jIdx
		}

		iRed, jRed := cds[i].Colour.R, cds[j].Colour.R
		if iRed != jRed {
			return iRed < jRed
		}

		iGreen, jGreen := cds[i].Colour.G, cds[j].Colour.G
		if iGreen != jGreen {
			return iGreen < jGreen
		}

		iBlue, jBlue := cds[i].Colour.B, cds[j].Colour.B
		if iBlue != jBlue {
			return iBlue < jBlue
		}

		return false
	})

	return cds
}
