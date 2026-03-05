package colour

import (
	"regexp"
)

// repl represents a regular expression and an associated substitution to be
// performed
type repl struct {
	re  *regexp.Regexp
	sub string
}

// replPair represents a pair of transformations of the spelling between UK
// and USA forms of words
type replPair struct {
	uk2usa repl
	usa2uk repl
}

var (
	apostropheRE = regexp.MustCompile(`'`)
	spaceRE      = regexp.MustCompile(`\s+`)
)

// makeSpellings constructs a slice of replPairs from an internal list
// of word pairs representing UK and US alternative spellings.
func makeSpellings() []replPair {
	spellings := []replPair{}

	//nolint:misspell
	wordPairs := []struct{ uk, usa string }{
		{uk: "colour", usa: "color"},
		{uk: "grey", usa: "gray"},
		{uk: "aluminium", usa: "aluminum"},
		{uk: "harbour", usa: "harbor"},
		{uk: "odour", usa: "odor"},
		{uk: "lustre", usa: "luster"},
	}
	for _, wp := range wordPairs {
		rp := replPair{
			uk2usa: repl{re: regexp.MustCompile(`\b` + wp.uk + `\b`), sub: wp.usa},
			usa2uk: repl{re: regexp.MustCompile(`\b` + wp.usa + `\b`), sub: wp.uk},
		}
		spellings = append(spellings, rp)
	}

	return spellings
}

// withAliases takes the passed map of colour names to RGBA values and
// generates aliases for some common values. For instance it maps:
//
// UK spellings to US spellings
// US spellings to UK spellings
//
//	xxx's   -> xxxs
//	xxx yyy -> xxx-yyy
//
// it then returns the original map with the alias values added.
func withAliases(m colourNameToRGBA) colourNameToRGBA {
	for name, c := range m {
		for alias := range transformColourNames(generateAltSpellings(name)) {
			if name == alias {
				continue
			}

			if _, ok := m[alias]; !ok {
				m[alias] = c
			}
		}
	}

	return m
}

// generateAltSpellings takes the colour name and generates extra names with
// alternative spellings applied. It will always return at least the supplied
// name.
func generateAltSpellings(name string) []string {
	alts := []string{name}
	spellings := makeSpellings()
	changed := name

	for _, rp := range spellings {
		changed = rp.uk2usa.re.ReplaceAllString(changed, rp.uk2usa.sub)
	}

	if changed != name {
		alts = append(alts, changed)
	}

	changed = name
	for _, rp := range spellings {
		changed = rp.usa2uk.re.ReplaceAllString(changed, rp.usa2uk.sub)
	}

	if changed != name {
		alts = append(alts, changed)
	}

	return alts
}

// transformColourNames takes the set of colour names and generates
// alternative values with a set of transformations applied
func transformColourNames(names []string) map[string]bool {
	transformed := map[string]bool{}
	replacements := []repl{
		{re: apostropheRE, sub: ""}, // delete apostrophes
		{re: spaceRE, sub: "-"},     // replace runs of space with a dash
	}

	for _, name := range names {
		transformed[name] = true

		for _, repl := range replacements {
			name = repl.re.ReplaceAllString(name, repl.sub)
		}

		transformed[name] = true
	}

	return transformed
}

// firstIsUK2ndIsUSA returns true if the first string matches the UK Regexp
// and the second matches the USA Regexp, false otherwise
func firstIsUK2ndIsUSA(ukRE, usaRE *regexp.Regexp, s1, s2 string) bool {
	return ukRE.MatchString(s1) && usaRE.MatchString(s2)
}

// firstMatches2ndDoesNot returns true if the first string matches the
// supplied Regexp and the second doesn't, false otherwise
func firstMatches2ndDoesNot(re *regexp.Regexp, s1, s2 string) bool {
	return re.MatchString(s1) && !re.MatchString(s2)
}

// choosePreferredAlias chooses the preferred alternative between the two
// supplied strings.
func choosePreferredAlias(s1, s2 string) string {
	transformations := []*regexp.Regexp{
		spaceRE,
		apostropheRE,
	}

	for _, re := range transformations {
		if firstMatches2ndDoesNot(re, s1, s2) {
			return s1
		}

		if firstMatches2ndDoesNot(re, s2, s1) {
			return s2
		}
	}

	spellings := makeSpellings()
	for _, sp := range spellings {
		if firstIsUK2ndIsUSA(sp.uk2usa.re, sp.usa2uk.re, s1, s2) {
			return s1
		}

		if firstIsUK2ndIsUSA(sp.uk2usa.re, sp.usa2uk.re, s2, s1) {
			return s2
		}
	}

	return s1 // we shouldn't get here but if we do just pick an alias
}

// IsAnAlias returns true and the preferred alias if name1 is an alias of
// name2 or vice versa. Otherwise it will return an empty string and false.
func IsAnAlias(name1, name2 string) (string, bool) {
	if name1 == name2 {
		return name1, true
	}

	name1Aliases := transformColourNames(generateAltSpellings(name1))
	name2Aliases := transformColourNames(generateAltSpellings(name2))

	if name1Aliases[name2] || // name2 is an alias of name1
		name2Aliases[name1] { // name1 is an alias of name2
		return choosePreferredAlias(name1, name2), true
	}

	return "", false
}

// stripAliasesOf returns a slice containing the strings from names with any
// aliases of s replaced by a single copy of the preferred alias.
//
// Note that the resulting slice always has at least one element: the target
// string itself.
func stripAliasesOf(s string, names []string) []string {
	stripped := []string{}

	for _, t := range names {
		if pref, alias := IsAnAlias(s, t); alias {
			s = pref
			continue
		}

		stripped = append(stripped, t)
	}

	stripped = append(stripped, s)

	return stripped
}

// StripAliases removes aliases from names and returns the resulting slice.
func StripAliases(names []string) []string {
	if len(names) == 0 {
		return names
	}

	checked := map[string]bool{}

	for len(names) > 1 {
		s := names[0]
		if checked[s] {
			break
		}

		stripped := stripAliasesOf(s, names[1:])
		checked[stripped[len(stripped)-1]] = true
		names = stripped
	}

	return names
}
