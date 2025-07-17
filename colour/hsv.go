package colour

// HSV represents a colour defined by Hue, Saturation and Value.
type HSV struct {
	// Hue is a number in the range [0, 360). Zero represents red, 120
	// represents green and 240 represents blue, with yellow, cyan and
	// magenta in the intervals.
	Hue float64
	// Saturation is a value in the range [0, 1]
	Saturation float64
	// Value is a value in the range [0, 1]
	Value float64
}
