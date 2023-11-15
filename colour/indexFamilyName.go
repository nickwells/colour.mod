package colour

import (
	"image/color"
)

// IndexFamilyName holds the index value calculated from the RGBA value and
// the corresponding family name and colour name.
type IndexFamilyName struct {
	Idx        uint32
	Family     Family
	ColourName string
}

// NamesByIndex returns a slice containing all the entries in the colour
// maps with the RGBA value converted to an index and the colour name and
// family.
func NamesByIndex() []IndexFamilyName {
	coloursByIdx := []IndexFamilyName{}
	for _, f := range searchOrder {
		for cName, c := range cFamMap[f] {
			coloursByIdx = append(coloursByIdx,
				IndexFamilyName{
					Idx:        colourIndex(c),
					Family:     f,
					ColourName: cName,
				})
		}
	}
	return coloursByIdx
}

// colourIndex constructs an index value from the R/G/B fields of the
// color.RGBA value. This is used internally to find the names of a colour
// from its RGBA value. Note that the Alpha component of the colour is
// ignored.
func colourIndex(c color.RGBA) uint32 {
	return (((uint32(c.R) << 8) +
		uint32(c.G)) << 8) +
		uint32(c.B)
}
