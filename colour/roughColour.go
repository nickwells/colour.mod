package colour

import (
	"image/color" //nolint:misspell
)

// RoughColour is used for classifying colour values
type RoughColour int

// The RoughColour values named here are used to roughly partition the hue
// values. The black, white and grey values cover the extremes of saturation
// and luminance.
const (
	RoughlyRed     RoughColour = iota // Red
	RoughlyGreen                      // Green
	RoughlyBlue                       // Blue
	RoughlyCyan                       // Cyan
	RoughlyMagenta                    // Magenta
	RoughlyYellow                     // Yellow
	RoughlyBlack                      // Black
	RoughlyGrey                       // Grey
	RoughlyWhite                      // White
	RoughlyOther                      // Other
)

const (
	band = colourInterval / 2
)

const (
	hueRed = iota * colourInterval
	hueYellow
	hueGreen
	hueCyan
	hueBlue
	hueMagenta

	bwThreshold       = 0.025
	lumThresholdBlack = 0 + bwThreshold
	lumThresholdWhite = 1 - bwThreshold

	greyThreshold = 0.125
)

// Roughly tries to clasify the given colour into one of the
// RoughColour values.
func Roughly(c color.RGBA) RoughColour { //nolint:misspell
	type hue2Name struct {
		hueVal float64
		rc     RoughColour
	}

	hues := []hue2Name{
		{hueRed, RoughlyRed},
		{hueGreen, RoughlyGreen},
		{hueBlue, RoughlyBlue},
		{hueCyan, RoughlyCyan},
		{hueMagenta, RoughlyMagenta},
		{hueYellow, RoughlyYellow},
	}

	hsl, _ := RGBA2HSLAndHSV(c)

	if hsl.Luminance < lumThresholdBlack {
		return RoughlyBlack
	}

	if hsl.Luminance > lumThresholdWhite {
		return RoughlyWhite
	}

	if hsl.Saturation < greyThreshold {
		return RoughlyGrey
	}

	for _, h := range hues {
		lb, ub := h.hueVal-band, h.hueVal+band
		if lb < 0 {
			lb += maxHue

			if hsl.Hue > lb || hsl.Hue < ub {
				return h.rc
			}
		}

		if hsl.Hue > lb && hsl.Hue < ub {
			return h.rc
		}
	}

	return RoughlyOther
}
