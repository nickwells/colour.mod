package colour

import (
	"errors"
	"fmt"
	"image/color" //nolint:misspell
	"maps"
	"math"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/nickwells/english.mod/english"
)

var rgbIntroRE = regexp.MustCompile(
	`^[[:space:]]*[rR][gG][bB]([aA])?[[:space:]]*\{`)
var rgbOutroRE = regexp.MustCompile(`[[:space:]]*\}[[:space:]]*$`)

var rgbDfltComponents = map[string]uint8{
	"R": 0,
	"G": 0,
	"B": 0,
	"A": math.MaxUint8,
}

// IsAPotentialColourString returns true if the given string could be a
// string encoding a colour in the form of a RGBA value. It only checks the
// start of the string. The string must start with "rgb" or "rgba" followed
// by a bracket ("{"). Arbitrary amounts of white space are allowed around
// the "rgb" or "rgba" and upper and lower case variants are treated the
// same.
func IsAPotentialColourString(s string) bool {
	return rgbIntroRE.MatchString(s)
}

// makeColour constructs an RGBA colour from the components.
func makeColour(components map[string]uint8) rgba {
	var c rgba

	c.R = components["R"]
	c.G = components["G"]
	c.B = components["B"]
	c.A = components["A"]

	return c
}

// extractColourComponentParts extracts the component name and value from the
// passed string. It will return an error if the string is not properly
// formed or if the component name is not recognised.
func extractColourComponentParts(s string) (string, string, error) {
	name, val, ok := strings.Cut(s, ":")
	if !ok {
		return name, val, fmt.Errorf(
			"bad colour component: %q,"+
				" the name and value should be separated by a colon(:)",
			s)
	}

	name = strings.TrimSpace(name)
	name = strings.ToUpper(name)

	if _, ok := rgbDfltComponents[name]; !ok {
		aval := slices.Collect(maps.Keys(rgbDfltComponents))
		slices.Sort(aval)

		return name, val, fmt.Errorf(
			"unknown colour component: %q, allowed values: %s",
			name, english.JoinQuoted(aval, ", ", " or "))
	}

	val = strings.TrimSpace(val)

	return name, val, nil
}

// ParseColourDefinition parses the given string into an RGBA colour.
func ParseColourDefinition(s string) (color.RGBA, error) { //nolint:misspell
	var c rgba

	if !rgbIntroRE.MatchString(s) {
		return c, fmt.Errorf("the colour definition is not started correctly")
	}

	if !rgbOutroRE.MatchString(s) {
		return c, fmt.Errorf(
			"the colour definition starts with %q but has no trailing %q",
			rgbIntroRE.FindString(s), "}")
	}

	strippedVal := rgbIntroRE.ReplaceAllString(s, "")
	strippedVal = rgbOutroRE.ReplaceAllString(strippedVal, "")

	components := maps.Clone(rgbDfltComponents)

	for part := range strings.SplitSeq(strippedVal, ",") {
		compName, value, err := extractColourComponentParts(part)
		if err != nil {
			return c, err
		}

		parsedVal, err := ParseColourPart(value, compName)
		if err != nil {
			return c, err
		}

		components[compName] = parsedVal
	}

	c = makeColour(components)

	return c, nil
}

// ParseColourPart takes the named part of a colour value (the Red, Green, Blue
// or Alpha component) as a string and converts it into an appropriate 8-bit
// value. It returns a non-nil error if the value cannot be converted.
func ParseColourPart(val, partName string) (uint8, error) {
	rVal, err := strconv.ParseUint(val, 0, 8)
	if err != nil {
		errIntro := fmt.Sprintf(
			"cannot convert the %q value (%q) to a valid number", partName, val)
		if errors.Is(err, strconv.ErrRange) {
			return 0, fmt.Errorf("%s: %w", errIntro, strconv.ErrRange)
		}

		if errors.Is(err, strconv.ErrSyntax) {
			return 0, fmt.Errorf("%s: %w", errIntro, strconv.ErrSyntax)
		}

		return 0, fmt.Errorf("%s: %w", errIntro, err)
	}

	return uint8(rVal), nil
}

// RGBAllowedValues returns a string describing an allowed RGB string
func RGBAllowedValues() string {
	aVal := strings.Builder{}
	aVal.WriteString("a string" +
		" giving the Red/Green/Blue/Alpha values as follows:" +
		" RGB{R: #, G: #, B: #, A: #}")

	defaults := map[string][]string{}

	for k, v := range rgbDfltComponents {
		dVal := fmt.Sprintf("%#4.2x", v)
		defaults[dVal] = append(defaults[dVal], k)
	}

	keys := slices.Collect(maps.Keys(defaults))
	sort.Strings(keys)
	aVal.WriteString(" (defaults: ")

	sep := ""
	for _, k := range keys {
		aVal.WriteString(sep)
		sep = ", "

		sort.Strings(defaults[k])
		aVal.WriteString(strings.Join(defaults[k], " / ") + ": " + k)
	}

	aVal.WriteString(").")
	aVal.WriteString(" Upper and lowercase values are treated equally")
	aVal.WriteString(" and whitespace is allowed anywhere.")

	return aVal.String()
}
