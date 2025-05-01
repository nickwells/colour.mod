package colour

import "image/color" //nolint:misspell

// Contrast calculates a colour that has a high contrast with the supplied
// colour
func Contrast(c color.RGBA) color.RGBA { //nolint:misspell
	return color.RGBA{ //nolint:misspell
		R: colourCounterpart(c.R),
		G: colourCounterpart(c.G),
		B: colourCounterpart(c.B),
		A: c.A,
	}
}

// colourCounterpart takes a colour component (one of the red, green
// or blue values in an RGBA colour) and calculates a value that will
// generate a contrasting colour.
func colourCounterpart(v uint8) uint8 {
	const (
		colourRange = 255
		halfRange   = colourRange / 2
		minDiff     = 3 * colourRange / 8
	)

	newV := colourRange - v

	if newV > halfRange {
		if (newV - v) < minDiff {
			newV += minDiff
		}
	} else {
		if (v - newV) < minDiff {
			newV -= minDiff
		}
	}

	return newV
}
