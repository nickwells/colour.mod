package colour

import (
	"fmt"
	"image/color" //nolint:misspell
	"strings"
)

// cmpColourPart checks that the difference between the parts is greater than
// prec and returns a non-empty string describing the problem if it is.
func cmpColourPart(p1, p2 uint8, precision int, name string) string {
	diff := int(p1) - int(p2)
	if diff < 0 {
		diff *= -1
	}

	if diff > precision {
		return fmt.Sprintf("%s differ by %d", name, diff)
	}

	return ""
}

// Compare returns an error if any component of the two colours is not within
// precision of the other. Passing a precision of zero will check for an
// exact match.
func Compare(c1, c2 color.RGBA, precision uint8) error { //nolint:misspell
	probs := []string{}

	if p := cmpColourPart(c1.R, c2.R, int(precision), "R's"); p != "" {
		probs = append(probs, p)
	}

	if p := cmpColourPart(c1.G, c2.G, int(precision), "G's"); p != "" {
		probs = append(probs, p)
	}

	if p := cmpColourPart(c1.B, c2.B, int(precision), "B's"); p != "" {
		probs = append(probs, p)
	}

	if p := cmpColourPart(c1.A, c2.A, int(precision), "A's"); p != "" {
		probs = append(probs, p)
	}

	if len(probs) == 0 {
		return nil
	}

	return fmt.Errorf(
		"colours (%#02v and %#02v) differ (precision: %d): %s",
		c1, c2, precision, strings.Join(probs, ", "))
}

// WithinDist returns true if the two colours have a Euclidean distance less
// than or equal to the dist value. Note that the range of the dimensions is
// [0,255] and so the maximum distance is the square root of 3 times 255
// squared or slightly more than 441. Note that the comparison is only on the
// red, green and blue components of the colour; the alpha value is not
// considered.
func WithinDist(c1, c2 color.RGBA, dist float64) bool { //nolint:misspell
	dist *= dist

	return float64(distSquared(c1, c2)) <= dist
}
