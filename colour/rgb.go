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

var (
	rgbIntroRE = regexp.MustCompile(
		`^[[:space:]]*[rR][gG][bB]([aA])?[[:space:]]*\{`)
	rgbOutroRE = regexp.MustCompile(`[[:space:]]*\}[[:space:]]*$`)

	rgbAlt6RE = regexp.MustCompile(
		`^[[:space:]]*` +
			`#` +
			`([[:xdigit:]][[:xdigit:]])` +
			`([[:xdigit:]][[:xdigit:]])` +
			`([[:xdigit:]][[:xdigit:]])` +
			`[[:space:]]*$`)
	rgbAlt3RE = regexp.MustCompile(
		`^[[:space:]]*` +
			`#` +
			`([[:xdigit:]])` +
			`([[:xdigit:]])` +
			`([[:xdigit:]])` +
			`[[:space:]]*$`)
)

var rgbDfltComponents = map[string]uint8{
	"R": 0,
	"G": 0,
	"B": 0,
	"A": math.MaxUint8,
}

// IsAPotentialColourString returns true if the given string could be a
// string encoding a colour in the form of a RGBA value. It only checks the
// start of the string. The string must either be a hash ("#") followed by
// either 3 or six hexadecimal digits (0-9, a-f) or else must start with "rgb"
// or "rgba" followed by a bracket ("{"). Arbitrary amounts of white space
// are allowed around the "rgb" or "rgba" and upper and lower case variants
// are treated the same.
func IsAPotentialColourString(s string) bool {
	if rgbAlt3RE.MatchString(s) {
		return true
	}

	if rgbAlt6RE.MatchString(s) {
		return true
	}

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

// Parse3DigitColour takes a string of the form "#xxx" where each "x" is a
// hexadecimal digit and returns a colour value and error.  Each digit is
// doubled so that a string of "#05f" will yield a colour with a red value of
// 0x0, a green of 0x55 and a blue of 0xff.
//
// A non-nil error is returned if all the digits cannot be extracted from the
// string (3 parts are expected) or if they cannot be parsed into an unsigned
// 8-bit number.
func Parse3DigitColour(s string) (color.RGBA, error) { //nolint:misspell
	return parseNDigitColour(s, "3-digit",
		rgbAlt3RE,
		func(xd string) string { return xd + xd })
}

// Parse6DigitColour takes a string of the form "#xxxxxx" where each "x" is a
// hexadecimal digit and returns a colour value and error. The digits are
// taken in pairs so that a string of "#0145ef" will yield a colour with a
// red value of 0x01, a green of 0x45 and a blue of 0xef,
//
// A non-nil error is returned if all the digits cannot be
// extracted from the string (3 parts are expected) or if they cannot be
// parsed into an unsigned 8-bit number.
func Parse6DigitColour(s string) (color.RGBA, error) { //nolint:misspell
	return parseNDigitColour(s, "6-digit",
		rgbAlt6RE,
		func(xd string) string { return xd })
}

// parseNDigitColour takes a string of N hexadecimal digit and returns a
// colour value and error. A non-nil error is returned if all the digits
// cannot be extracted from the string (3 parts are expected) or if they
// cannot be parsed into an unsigned 8-bit number.
func parseNDigitColour(
	s, name string,
	re *regexp.Regexp,
	mkDigits func(string) string,
) (
	color.RGBA, //nolint:misspell
	error,
) {
	idx2Component := []string{
		"R",
		"G",
		"B",
	}

	errIntro := fmt.Sprintf("the %s colour (%q) is badly formed", name, s)

	// FindStringSubmatch returns a slice of strings, the 0th entry is the
	// whole string and the remaining parts, of which there should be 3, are
	// expected to be the hex digits for the Red, Green and Blue components
	// respectively. Note that there is no Alpha value
	parts := re.FindStringSubmatch((s))
	if len(parts) == 0 {
		return rgba{}, errors.New(errIntro)
	}

	// strip the first entry (the match of the whole string)
	parts = parts[1:]
	if len(parts) != len(idx2Component) {
		return rgba{}, fmt.Errorf("%s: %d parts expected, %d found",
			errIntro,
			len(idx2Component),
			len(parts))
	}

	components := maps.Clone(rgbDfltComponents)

	for i, xd := range parts {
		v, err := strconv.ParseUint(mkDigits(xd), 16, 8)
		if err != nil {
			return rgba{},
				fmt.Errorf("%s: digit %d(%s) cannot be converted to a number",
					errIntro,
					i+1, xd)
		}

		components[idx2Component[i]] = uint8(v)
	}

	return makeColour(components), nil
}

// ParseColourDefinition parses the given string into an RGBA colour.
func ParseColourDefinition(s string) (color.RGBA, error) { //nolint:misspell
	var c rgba

	if rgbAlt3RE.MatchString(s) {
		return Parse3DigitColour(s)
	}

	if rgbAlt6RE.MatchString(s) {
		return Parse6DigitColour(s)
	}

	if !rgbIntroRE.MatchString(s) {
		return c, fmt.Errorf("the colour definition (%q) is invalid", s)
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
	aVal.WriteString("\n\n")
	aVal.WriteString(`Or a literal hash ("#")`)
	aVal.WriteString(" immediately followed by")
	aVal.WriteString(" precisely 3 or 6 hexadecimal digits")

	return aVal.String()
}
