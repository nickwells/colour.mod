package colour

import (
	"fmt"
	"image/color" //nolint:misspell
	"math"

	"github.com/nickwells/mathutil.mod/v2/mathutil"
)

// Saturation returns a colour with the same Hue and Luminance as the
// supplied colour but with the saturation set to the supplied value. The
// supplied saturation must be between zero and one inclusive otherwise an
// error will be returned. Supplying a saturation equal to (or very close to)
// that of the original colour will have no effect.
func Saturation(
	c color.RGBA, saturation float64, //nolint:misspell
) (
	color.RGBA, error, //nolint:misspell
) {
	if saturation < 0 {
		return c,
			fmt.Errorf("the saturation (%.2f) must be >= 0", saturation)
	}

	if saturation > 1 {
		return c,
			fmt.Errorf("the saturation (%.2f) must be <= 1", saturation)
	}

	hsl, _ := RGBA2HSLAndHSV(c)

	const epsilon = 0.00001
	if mathutil.AlmostEqual(hsl.Saturation, saturation, epsilon) {
		return c, nil
	}

	hsl.Saturation = saturation

	return hsl.ToRGBA(), nil
}

// Luminance returns a colour with the same Hue and Saturation as the
// supplied colour but with the luminance set to the supplied value. The
// supplied luminance must be between zero and one inclusive otherwise an
// error will be returned. Supplying a luminance equal to (or very close to)
// that of the original colour will have no effect.
func Luminance(
	c color.RGBA, luminance float64, //nolint:misspell
) (
	color.RGBA, error, //nolint:misspell
) {
	if luminance < 0 {
		return c,
			fmt.Errorf("the luminance (%.2f) must be >= 0", luminance)
	}

	if luminance > 1 {
		return c,
			fmt.Errorf("the luminance (%.2f) must be <= 1", luminance)
	}

	hsl, _ := RGBA2HSLAndHSV(c)

	const epsilon = 0.00001
	if mathutil.AlmostEqual(hsl.Luminance, luminance, epsilon) {
		return c, nil
	}

	hsl.Luminance = luminance

	return hsl.ToRGBA(), nil
}

// Invert returns the inverted value of the colour. Each of the red, green
// and blue components are subtracted from the max value and the resulting
// colour is generated from these values.
func Invert(c color.RGBA) color.RGBA { //nolint:misspell
	c.R = math.MaxUint8 - c.R
	c.G = math.MaxUint8 - c.G
	c.B = math.MaxUint8 - c.B

	return c
}

// Complement returns the complementary colour. This is that colour having
// the same Luminance and Saturation but the 'opposite' Hue - it sits at the
// opposite side of the colour wheel. Note that shades of grey (colours with
// zero saturation) from black to white are unchanged - there is no
// complementary colour.
func Complement(c color.RGBA) color.RGBA { //nolint:misspell
	const (
		maxDegrees     = 360
		halfMaxDegrees = 180
	)

	hsl, _ := RGBA2HSLAndHSV(c)
	if hsl.Saturation == 0 {
		return c
	}

	hsl.Hue = (hsl.Hue + halfMaxDegrees)
	if hsl.Hue > maxDegrees {
		hsl.Hue -= maxDegrees
	}

	return hsl.ToRGBA()
}
