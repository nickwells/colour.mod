package colour

import (
	"image/color" //nolint:misspell
	"math"
)

// ContrastColourful calculates a colour that has a high contrast with the
// supplied colour and is typically a brighter colour than is given by the
// [Constrast] func
func ContrastColourful(c color.RGBA) color.RGBA { //nolint:misspell
	return rgba{
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
		halfRange = math.MaxUint8 / 2
		minDiff   = 3 * math.MaxUint8 / 8
	)

	newV := math.MaxUint8 - v

	if newV > halfRange {
		if (newV - v) < minDiff {
			newV += (math.MaxUint8 - newV) / 2 //nolint:mnd
		}
	} else {
		if (v - newV) < minDiff {
			newV /= 2 //nolint:mnd
		}
	}

	return newV
}

// luminanceThresholdAdjustment returns the threshold and the associated
// adjustment to be made to the luminance. It selects different values
// according to the calculated RoughColour value of the supplied colour.
func luminanceThresholdAdjustment(c rgba) (
	float64, float64,
) {
	const (
		lumThresholdGrey    = 0.5
		lumThresholdRed     = 0.6
		lumThresholdGreen   = 0.3
		lumThresholdBlue    = 0.7
		lumThresholdCyan    = 0.35
		lumThresholdMagenta = 0.55
		lumThresholdYellow  = 0.3
		lumThresholdOther   = 0.6

		stdAdjustment   = 0.1
		maxAdjustment   = 0
		otherAdjustment = 0.05
	)

	switch Roughly(c) {
	case RoughlyRed:
		return lumThresholdRed, stdAdjustment
	case RoughlyGreen:
		return lumThresholdGreen, stdAdjustment
	case RoughlyBlue:
		return lumThresholdBlue, stdAdjustment
	case RoughlyCyan:
		return lumThresholdCyan, stdAdjustment
	case RoughlyMagenta:
		return lumThresholdMagenta, stdAdjustment
	case RoughlyYellow:
		return lumThresholdYellow, stdAdjustment
	case RoughlyGrey:
		return lumThresholdGrey, maxAdjustment
	case RoughlyBlack:
		return lumThresholdBlack, maxAdjustment
	case RoughlyWhite:
		return lumThresholdWhite, maxAdjustment
	}

	return lumThresholdOther, otherAdjustment
}

// Contrast calculates a colour that has a high contrast with the supplied
// colour. It adjusts the brightness of the colour generated according to the
// colour supplied.
func Contrast(c color.RGBA) color.RGBA { //nolint:misspell
	const oppositeHue = maxHue / 2

	hsl, _ := RGBA2HSLAndHSV(c)
	hsl.Hue = math.Mod(hsl.Hue+oppositeHue, maxHue)

	threshold, adjustment := luminanceThresholdAdjustment(c)

	if hsl.Luminance < threshold {
		hsl.Luminance = 1 - adjustment
	} else {
		hsl.Luminance = adjustment
	}

	return hsl.ToRGBA()
}
