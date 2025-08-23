package colour

import (
	"errors"
	"fmt"
	"image/color" //nolint:misspell
	"strings"
)

// cmpColourPart checks that the difference between the parts is greater than
// prec and returns a non-empty string describing the problem if it is.
func cmpColourPart(p1, p2 uint8, precision int, name string) string {
	epsilon := precision * precision
	diff := int(p1) - int(p2)
	d2 := diff * diff

	if d2 > epsilon {
		return fmt.Sprintf(
			"c1.%s(%#02x) and c2.%s(%#02x) differ by more than %d",
			name, p1, name, p2, precision)
	}

	return ""
}

// Compare returns an error if any component of the two colours is not within
// precision of the other. Passing a precision of zero will check for an
// exact match.
func Compare(c1, c2 color.RGBA, precision uint8) error { //nolint:misspell
	probs := []string{}

	if p := cmpColourPart(c1.R, c2.R, int(precision), "R"); p != "" {
		probs = append(probs, p)
	}

	if p := cmpColourPart(c1.G, c2.G, int(precision), "G"); p != "" {
		probs = append(probs, p)
	}

	if p := cmpColourPart(c1.B, c2.B, int(precision), "B"); p != "" {
		probs = append(probs, p)
	}

	if p := cmpColourPart(c1.A, c2.A, int(precision), "A"); p != "" {
		probs = append(probs, p)
	}

	if len(probs) == 0 {
		return nil
	}

	return errors.New("the colours differ:" + strings.Join(probs, ", "))
}

// WithinDist returns true if the two colours have a Euclidean distance less
// than the limit value.
func WithinDist(c1, c2 color.RGBA, limit float64) bool { //nolint:misspell
	limit *= limit

	return float64(dist(c1, c2)) < limit
}
