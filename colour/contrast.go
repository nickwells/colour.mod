package colour

import (
	"image/color" //nolint:misspell
	"math"
)

// ContrastColourful calculates a colour that has a high contrast with the
// supplied colour and is typically a brighter colour than is given by the
// [Constrast] func
func ContrastColourful(c color.RGBA) color.RGBA { //nolint:misspell
	cc := rgba{
		R: colourCounterpart(c.R),
		G: colourCounterpart(c.G),
		B: colourCounterpart(c.B),
		A: c.A,
	}

	// note that the largest possible minimum is for the colour at the centre
	// of the cube (where the R, G and B values are all set to 0x80). This is
	// a little under 222 so we choose (arbitrarily) 150 as a reasonable target.
	const minDistSquared = 150 * 150
	for distSquared(c, cc) < minDistSquared {
		cc = adjustColour(c, cc)
	}

	return cc
}

// adjustColourComponent returns the component (comp) shifted in the given
// direction. If direction is less than 0 then towards the maximum value,
// otherwise towards the minimum. It moves in steps of half the difference
// between the value and the target value (max or zero) but with a magnitude
// of at least 10.
func adjustColourComponent(comp uint8, direction int) uint8 {
	const (
		divisor = 2
		minStep = 10
	)

	if direction < 0 {
		diff := math.MaxUint8 - comp
		step := diff / divisor
		step = max(step, minStep)
		step = min(step, diff)

		return comp + step
	}

	step := comp / divisor
	step = max(step, minStep)
	step = min(step, comp)

	return comp - step
}

// adjustColour takes the two colours and finds the Red, Green or Blue
// component for which the contrasting colour value (cc) is closest to the
// original colour value (c)
func adjustColour(c, cc rgba) rgba {
	rDiff := int(c.R) - int(cc.R)
	gDiff := int(c.G) - int(cc.G)
	bDiff := int(c.B) - int(cc.B)

	rDiff2 := rDiff * rDiff
	gDiff2 := gDiff * gDiff
	bDiff2 := bDiff * bDiff

	// if the red component has the smallest distance
	if rDiff2 <= gDiff2 && rDiff2 <= bDiff2 {
		cc.R = adjustColourComponent(cc.R, rDiff)

		return cc
	}

	// if the green component has the smallest distance
	if gDiff2 <= bDiff2 {
		cc.G = adjustColourComponent(cc.G, gDiff)

		return cc
	}

	// the blue component has the smallest distance
	cc.B = adjustColourComponent(cc.B, bDiff)

	return cc
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
