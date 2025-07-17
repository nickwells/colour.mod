package colour

import (
	"image/color" //nolint:misspell
	"math"
	"math/rand/v2"
)

// MakeColours generates up to count, randomly generated colours (without
// duplicates) across the range of all colours. It will return an error if
// count <= 0. The attempt to generate the set of colours will end
// prematurely if the attempt to generate the next colour yields a
// preexisting colour more than twice; this means the larger the value of
// count, the more chance that the set of colours will have fewer entries
// than count.
func MakeColours(count int) ([]color.RGBA, error) { //nolint:misspell
	const uint8Range = math.MaxUint8 + 1

	if count <= 0 {
		return nil, badColourCountErr(count)
	}

	colours := make([]rgba, 0, count)
	existingColours := map[rgba]bool{}

Loop:
	for range count {
		const maxAttempts = 3

		for range maxAttempts {
			//nolint:gosec
			c := rgba{
				R: uint8(rand.IntN(uint8Range)),
				G: uint8(rand.IntN(uint8Range)),
				B: uint8(rand.IntN(uint8Range)),
				A: math.MaxUint8,
			}
			if !existingColours[c] {
				colours = append(colours, c)
				existingColours[c] = true

				continue Loop
			}
		}

		break Loop
	}

	return colours, nil
}

// MakeColoursBetween generates up to count colours (without duplicates)
// across the range of colours between lower and upper (inclusive). It will
// return an error if count <= 0. If count == 1 then a single colour is
// generated mid-way between the lower and upper bound. Otherwise it will
// generate the colours in evenly-spaced bands.
func MakeColoursBetween(count int, lower, upper color.RGBA) ( //nolint:misspell
	[]color.RGBA, //nolint:misspell
	error,
) {
	switch {
	case count <= 0:
		return nil, badColourCountErr(count)
	case lower == upper:
		return []rgba{lower}, nil
	case count == 1:
		//nolint:gosec,mnd
		mid := rgba{
			R: uint8((int(lower.R) + int(upper.R)) / 2),
			G: uint8((int(lower.G) + int(upper.G)) / 2),
			B: uint8((int(lower.B) + int(upper.B)) / 2),
			A: uint8((int(lower.A) + int(upper.A)) / 2),
		}

		return []rgba{mid}, nil
	case count == 2: //nolint:mnd
		return []rgba{lower, upper}, nil
	}

	lowerHSL, _ := RGBA2HSLAndHSV(lower)
	upperHSL, _ := RGBA2HSLAndHSV(upper)
	div := float64(count - 1)

	startHue := lowerHSL.Hue
	hueInterval := (upperHSL.Hue - lowerHSL.Hue) / div

	startSat := lowerHSL.Saturation
	satInterval := (upperHSL.Saturation - lowerHSL.Saturation) / div

	startLum := lowerHSL.Luminance
	lumInterval := (upperHSL.Luminance - lowerHSL.Luminance) / div

	startA := lower.A
	aInterval := float64(upper.A-lower.A) / div

	colours := []rgba{lower}
	lastColour := lower

	for i := range count - 2 {
		scale := float64(i + 1)
		nextHSL := HSL{
			Hue:        startHue + scale*hueInterval,
			Saturation: startSat + scale*satInterval,
			Luminance:  startLum + scale*lumInterval,
		}

		nextColour := nextHSL.ToRGBA()
		nextColour.A = startA + uint8(scale*aInterval)

		if lastColour == nextColour {
			continue // the colour hasn't changed so skip it
		}

		lastColour = nextColour

		colours = append(colours, nextColour)
	}

	if lastColour != upper {
		colours = append(colours, upper)
	}

	return colours, nil
}
