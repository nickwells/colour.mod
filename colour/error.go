package colour

import "fmt"

// These error strings can be used with errors.Is to test the error returned
// to see if it is of one of these types and hence the meaning of the
// Error.Val. See [Error].
const (
	// BadFamily is the text of the error message when there is a problem
	// with the Family name
	BadFamily = "bad colour family"
	// BadColourName is the text of the error message when there is a problem
	// with the Colour name
	BadColourName = "bad colour name"
	// BadColourCount is the text of the error message when there is a problem
	// with the range count when generating a colour range
	BadColourCount = "bad colour count"
	// BadColourProximity is the text of the error message when there is a
	// problem with the colour proximity (used to determine closeness to a
	// reference colour)
	BadColourProximity = "bad colour proximity"
)

// Error is the type of an error from the colour package
type Error struct {
	Text       string
	Value      string
	Proximity  float64
	Family     Family
	ColourName string
	Count      int
}

// badFamilyErr returns an Error with Text set to BadFamily and Family set
// to the Family name
func badFamilyErr(f Family) Error {
	return Error{
		Text:   BadFamily,
		Value:  fmt.Sprintf("%q", f),
		Family: f,
	}
}

// badColourErr returns an Error with Text set to BadColour and ColourName
// set to the colour name
func badColourErr(c string) Error {
	return Error{
		Text:       BadColourName,
		Value:      fmt.Sprintf("%q", c),
		ColourName: c,
	}
}

// badColourCountErr returns an Error with Text set to BadColourCount and
// Count set to the bad value
func badColourCountErr(count int) Error {
	return Error{
		Text:  BadColourCount,
		Value: fmt.Sprintf("expecting 0 < %d", count),
		Count: count,
	}
}

// badProximityErr returns an Error with Text set to BadColourProximity and
// Value set to a string representation of the bad value
func badProximityErr(p float64) Error {
	return Error{
		Text:      BadColourCount,
		Value:     fmt.Sprintf("expecting 0 <= %f < %f", p, maxProximity),
		Proximity: p,
	}
}

// Error returns the Error formatted as a string
func (err Error) Error() string {
	return err.Text + ": " + err.Value
}

// Is returns true if the Error.Text matches the target error string
func (err Error) Is(target error) bool {
	if tErr, ok := target.(Error); ok {
		return tErr.Text == err.Text
	}

	return err.Text == target.Error()
}
