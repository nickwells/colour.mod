package colour

import (
	"fmt"
	"image/color" //nolint:misspell
	"math"
)

// The following constants give the red, green and blue weightings used for
// converting RGBA colours to grey-scale equivalents. There are various sets
// given corresponding to various stnadards.In each case they are intended to
// reflect the differing ways that various colours are perceived by the human
// eye. Green is typically given the highest weighting and blue the lowest.

// Colour weighting constants used in PAL video systems (also SECAM and
// NTSC).
const (
	wtRedPAL   = 0.299
	wtGreenPAL = 0.587
	wtBluePAL  = 0.114
)

// Colour weighting constants used in the ITU-R BT.709 standard
const (
	wtRedBT709   = 0.2126
	wtGreenBT709 = 0.7152
	wtBlueBT709  = 0.0722
)

// Colour weighting constants used in the ITU-R BT.2100 standard
const (
	wtRedBT2100   = 0.2627
	wtGreenBT2100 = 0.6780
	wtBlueBT2100  = 0.0593
)

// ToGrey returns a grey-scale equivalent to the given colour. It uses
// weighted conversions of the colour components (the red, green and blue
// values) rather than equal weightings.
//
// It uses the PAL weights to convert the colour component values.
//
// See also [ToGreyEqual], [ToGreyBT709], [ToGreyBT2100] and [ToGreyCustom].
func ToGrey(c color.RGBA) color.RGBA { //nolint:misspell
	g, err := ToGreyCustom(c, wtRedPAL, wtGreenPAL, wtBluePAL)
	if err != nil {
		panic(fmt.Errorf("unexpected error (PAL): %w", err))
	}

	return g
}

// ToGreyBT709 returns a grey-scale equivalent to the given colour. It
// uses weighted conversions of the colour components (the red, green and
// blue values) rather than equal weightings.
//
// It uses weights from the ITU-R BT.709 standard to convert the colour
// component values.
//
// See also [ToGreyEqual], [ToGrey], [ToGreyBT2100] and [ToGreyCustom].
func ToGreyBT709(c color.RGBA) color.RGBA { //nolint:misspell
	g, err := ToGreyCustom(c, wtRedBT709, wtGreenBT709, wtBlueBT709)
	if err != nil {
		panic(fmt.Errorf("unexpected error (BT.709): %w", err))
	}

	return g
}

// ToGreyBT2100 returns a grey-scale equivalent to the given colour. It
// uses weighted conversions of the colour components (the red, green and
// blue values) rather than equal weightings.
//
// It uses weights from the ITU-R BT.2100 standard to convert the colour
// component values.
//
// See also [ToGreyEqual], [ToGrey], [ToGreyBT709]  and [ToGreyCustom].
func ToGreyBT2100(c color.RGBA) color.RGBA { //nolint:misspell
	g, err := ToGreyCustom(c,
		wtRedBT2100, wtGreenBT2100, wtBlueBT2100)
	if err != nil {
		panic(fmt.Errorf("unexpected error (BT.2100): %w", err))
	}

	return g
}

// ToGreyEqual returns a grey-scale equivalent to the given colour. It uses
// equal weightings of the colour components (the red, green and blue
// values).
// See also [ToGrey], [ToGreyBT709], [ToGreyBT2100]  and [ToGreyCustom].
func ToGreyEqual(c color.RGBA) color.RGBA { //nolint:misspell
	g, err := ToGreyCustom(c, 1, 1, 1)
	if err != nil {
		panic(fmt.Errorf("unexpected error (equal): %w", err))
	}

	return g
}

// ToGreyCustom returns a grey-scale equivalent to the given colour. It uses
// the supplied weights of the colour components (the red, green and blue
// values). It returns a non-nil error if the weights sum to zero or if any
// of them is less than zero.
//
// See also [ToGrey], [ToGreyBT709], [ToGreyBT2100] and [ToGreyEqual].
func ToGreyCustom(c color.RGBA, //nolint:misspell
	wtRed,
	wtGreen,
	wtBlue float64,
) (color.RGBA, error) { //nolint:misspell
	if wtRed < 0 {
		return c,
			fmt.Errorf("the red weight (%f) is less than zero", wtRed)
	}

	if wtGreen < 0 {
		return c,
			fmt.Errorf("the green weight (%f) is less than zero", wtGreen)
	}

	if wtBlue < 0 {
		return c,
			fmt.Errorf("the blue weight (%f) is less than zero", wtBlue)
	}

	sumW := wtRed + wtGreen + wtBlue

	if sumW == 0 {
		return c, fmt.Errorf("the sum of the weights is zero")
	}

	wtRed /= sumW
	wtGreen /= sumW
	wtBlue /= sumW

	greyVal := toUint8(
		0 +
			(float64(c.R) * wtRed) +
			(float64(c.G) * wtGreen) +
			(float64(c.B) * wtBlue))

	return MakeGrey(greyVal), nil
}

// MakeGrey constructs an RGBA value with each of the red, green and blue
// values set to the supplied greyVal and the alpha set to the maximum value
func MakeGrey(greyVal uint8) color.RGBA { //nolint:misspell
	return color.RGBA{ //nolint:misspell
		R: greyVal,
		G: greyVal,
		B: greyVal,
		A: math.MaxUint8,
	}
}

// toUint8 returns the value rounded to the nearest integer and then forced
// to a valid uint8 value.
func toUint8(v float64) uint8 {
	v = math.Round(v)

	if v > math.MaxUint8 {
		return math.MaxUint8
	}

	if v < 0 {
		return 0
	}

	return uint8(v)
}
