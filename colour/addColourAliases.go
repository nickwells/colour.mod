package colour

import "regexp"

// repl represents a regular expression and an associated substitution to be
// performed
type repl struct {
	re  *regexp.Regexp
	sub string
}

// withAliases takes the passed map of colour names to RGBA values and
// generates aliases for some common values. For instance it maps:
//
//	grey    -> gray
//	gray    -> grey
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

// generateAltSpellings takes the colour name and generates a set of names
// with alternative spellings applied.
func generateAltSpellings(name string) []string {
	alts := []string{name}
	altSpellings := []repl{
		{re: regexp.MustCompile(`\bgrey\b`), sub: "gray"},
		{re: regexp.MustCompile(`\bgray\b`), sub: "grey"},
	}

	for _, repl := range altSpellings {
		changed := repl.re.ReplaceAllString(name, repl.sub)
		if changed != name {
			alts = append(alts, changed)
		}
	}

	return alts
}

// transformColourNames takes the set of colour names and generates
// alternative values with a set of transformations applied
func transformColourNames(names []string) map[string]bool {
	transformed := map[string]bool{}
	replacements := []repl{
		// delete apostrophes
		{re: regexp.MustCompile(`'`), sub: ""},
		// replace runs of space with a dash
		{re: regexp.MustCompile(`\s+`), sub: "-"},
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
