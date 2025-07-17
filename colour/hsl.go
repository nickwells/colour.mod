package colour

import (
	"fmt"
	"image/color" //nolint:misspell
	"math"
)

// HSL represents a colour defined by Hue, Saturation and Luminance (or
// Lightness or Level).
type HSL struct {
	// Hue is a number in the range [0, 360). Zero represents red, 120
	// represents green and 240 represents blue, with yellow, cyan and
	// magenta in the intervals.
	Hue float64
	// Saturation is a value in the range [0, 1]. The lower the value, the
	// greyer the resulting colour.
	Saturation float64
	// Luminance is a value in the range [0, 1]. The larger this value, the
	// brighter the colour, the smaller the darker. For any Hue or
	// Saturation, a value of 0 results in the colour black and a value of 1
	// results in the colour white.
	Luminance float64
}

// String returns a string representation of the HSL value
func (hsl HSL) String() string {
	return fmt.Sprintf("{H:%3.0f S:%0.3f L:%0.3f}",
		hsl.Hue, hsl.Saturation, hsl.Luminance)
}

// ToRGBA converts an HSL colour value into an RGBA value.  The alpha value
// is forced to 0xfff. Note that the conversions between HSL and RGBA values
// are lossy ; that is, converting an RGBA value to an HSL value and back
// again is not guaranteed to generate the original colour.
func (hsl HSL) ToRGBA() color.RGBA { //nolint:misspell
	chroma := (1 - math.Abs(2*hsl.Luminance-1)) * //nolint:mnd
		hsl.Saturation

	h := math.Mod(hsl.Hue, maxHue) / colourInterval
	x := chroma * (1 - math.Abs(math.Mod(h, 2)-1)) //nolint:mnd
	m := hsl.Luminance - chroma/2                  //nolint:mnd

	var r, g, b float64

	if h <= hueYellow/colourInterval {
		r, g, b = chroma, x, 0
	} else if h <= hueGreen/colourInterval {
		r, g, b = x, chroma, 0
	} else if h <= hueCyan/colourInterval {
		r, g, b = 0, chroma, x
	} else if h <= hueBlue/colourInterval {
		r, g, b = 0, x, chroma
	} else if h <= hueMagenta/colourInterval {
		r, g, b = x, 0, chroma
	} else {
		r, g, b = chroma, 0, x
	}

	r = (r + m) * math.MaxUint8
	g = (g + m) * math.MaxUint8
	b = (b + m) * math.MaxUint8

	return rgba{
		R: uint8(r),
		G: uint8(g),
		B: uint8(b),
		A: math.MaxUint8,
	}
}

// RGBA satisfies the Color interface from the [image/color] package
//
//nolint:misspell
func (hsl HSL) RGBA() (r, g, b, a uint32) {
	c := hsl.ToRGBA()
	return c.RGBA()
}
