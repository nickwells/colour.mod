package colour

import (
	"fmt"
	"image/color" //nolint:misspell
	"slices"
	"strings"

	"github.com/nickwells/english.mod/english"
)

// qualifiedColourName holds a colour name and the families having that name
// for the colour.
type qualifiedColourName struct {
	families Families
	cName    string
}

// Describe returns a string representation of the colour. If an exact match
// is not found then the RGB values are shown. Otherwise the shortest name
// for the colour in each family is used and if only one name is found then
// that is returned without any Family-qualification. It searches in the
// standard families. See [Families.Describe] .
func Describe(c color.RGBA) string { //nolint:misspell
	return standardFamilies.Describe(c)
}

// Describe returns a string representation of the colour. If an exact match
// is not found then the RGB value is shown. Otherwise the shortest name for
// the colour in each family is used and if only one name is found then that
// is returned without any Family-qualification
func (fl Families) Describe(c color.RGBA) string { //nolint:misspell
	colours, err := fl.ClosestWithin(c, 0)
	if err != nil {
		return err.Error()
	}

	if len(colours) == 0 {
		return fmt.Sprintf("%#4.2v", c)
	}

	if len(colours) == 1 {
		cName := colours[0].CNames[0]
		for _, altCName := range colours[0].CNames[1:] {
			cName = preferredName(cName, altCName)
		}

		return cName
	}

	qcns := getFamilyNames(colours)

	if len(qcns) == 1 {
		return qcns[0].cName
	}

	desc := []string{}

	for _, qcn := range qcns {
		for _, f := range qcn.families {
			val := fmt.Sprintf("%q", f.String()+":"+qcn.cName)
			desc = append(desc, val)
		}
	}

	return english.Join(desc, ", ", " or ")
}

// getFamilyNames returns a list of qualifiedColourName values in order of,
// firstly number of families having the colour name, secondly by the length
// of the name and finally by the lexical order of the colour name.
func getFamilyNames(colours []FamilyColour) []qualifiedColourName {
	familiesByColour := map[string]Families{}
	qcns := []qualifiedColourName{}

	for _, cd := range colours {
		cName := cd.CNames[0]
		for _, altCName := range cd.CNames[1:] {
			cName = preferredName(cName, altCName)
		}

		familiesByColour[cName] = append(familiesByColour[cName], cd.Family)
	}

	for cn, fl := range familiesByColour {
		qcns = append(qcns, qualifiedColourName{cName: cn, families: fl})
	}

	slices.SortFunc(qcns, func(a, b qualifiedColourName) int {
		rval := len(b.families) - len(a.families)
		if rval == 0 {
			rval = len(a.cName) - len(b.cName)
		}

		if rval == 0 {
			if a.cName < b.cName {
				return 1
			}

			return -1
		}

		return rval
	})

	return qcns
}

// preferredName returns the preferred name. We prefer shorter names over
// longer; names without digits over names with; names coming earlier in
// lexical order over names coming later.
func preferredName(cName1, cName2 string) string {
	const digits = "1234567890"

	if len(cName2) < len(cName1) {
		return cName2 // prefer shorter names
	}

	if len(cName2) > len(cName1) {
		return cName1 // prefer shorter names
	}

	var (
		cn1HasDigits = strings.ContainsAny(cName1, digits)
		cn2HasDigits = strings.ContainsAny(cName2, digits)
	)

	if cn1HasDigits && !cn2HasDigits {
		return cName2 // prefer names without digits
	}

	if !cn1HasDigits && cn2HasDigits {
		return cName1 // prefer names without digits
	}

	if cName2 < cName1 {
		return cName2 // prefer name coming 1st in lexical order
	}

	return cName1
}
