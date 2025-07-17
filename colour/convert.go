package colour

import (
	"image/color" //nolint:misspell
	"math"
)

const (
	// colourInterval is the distance (in degrees) around the hue circle
	// between the primary and complementary colours
	colourInterval = 60
	// maxHue is the upperbound on the hue values (the maximum angular
	// distance around the hue circle).
	maxHue = 360
)

// rgbNormalised generates red, green and blue values (in that order)
// normalised to the range [0, 1] from an RGBA colour value.
func rgbNormalised(c color.RGBA) ( //nolint:misspell
	float64, float64, float64,
) {
	return float64(c.R) / math.MaxUint8,
		float64(c.G) / math.MaxUint8,
		float64(c.B) / math.MaxUint8
}

// RGBA2HSLAndHSV converts an RGBA colour value into HSL and HSV colour
// values. The conversion is lossy and converting from an RGBA to an HSL
// colour and back again will not necessarily yield the original RGBA value.
func RGBA2HSLAndHSV(c color.RGBA) (HSL, HSV) { //nolint:misspell
	r, g, b := rgbNormalised(c)
	xMin, xMax := min(r, g, b), max(r, g, b)
	chroma := xMax - xMin
	value := xMax
	luminance := (xMax + xMin) / 2 //nolint:mnd

	var hue, saturationHSL, saturationHSV float64

	if chroma == 0 {
		hue = 0
	} else if value == r {
		hue = math.Mod((g-b)/chroma, maxHue/colourInterval) +
			(hueRed / colourInterval)
	} else if value == g {
		hue = math.Mod((b-r)/chroma, maxHue/colourInterval) +
			(hueGreen / colourInterval)
	} else if value == b {
		hue = math.Mod((r-g)/chroma, maxHue/colourInterval) +
			(hueBlue / colourInterval)
	}

	hue *= colourInterval
	if hue < 0 {
		hue += maxHue
	}

	if value == 0 {
		saturationHSV = 0
	} else {
		saturationHSV = chroma / value
	}

	if luminance == 0 || luminance == 1 {
		saturationHSL = 0
	} else {
		saturationHSL = (value - luminance) / min(luminance, 1-luminance)
	}

	return HSL{
			Hue:        hue,
			Saturation: saturationHSL,
			Luminance:  luminance,
		}, HSV{
			Hue:        hue,
			Saturation: saturationHSV,
			Value:      value,
		}
}
