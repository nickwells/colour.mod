package colour

import (
	"cmp"
	"fmt"
	"image/color"
	"maps"
	"slices"
	"sort"
	"strings"

	"github.com/nickwells/english.mod/english"
)

// Describe returns a string representation of the colour. If it is not found
// in the colour to name map then the RGB values are shown. Otherwise the
// shortest name for the colour in each family is used and if only one name
// is found then that is returned without any Family-qualification.
func Describe(c color.RGBA) string {
	qNames, ok := colourToNameMap[colourIndex(c)]
	if !ok {
		return fmt.Sprintf("%#v", c)
	}

	if len(qNames) == 1 {
		return qNames[0].ColourName
	}

	qNames = getUniqueNames(qNames)

	if len(qNames) == 1 {
		return qNames[0].ColourName
	}

	distinctNames, keys := getDistinctNames(qNames)

	if len(distinctNames) == 1 {
		return qNames[0].ColourName
	}

	desc := []string{}

	for _, k := range keys {
		val := k + " ("
		val += familyList(distinctNames[k])
		val += ")"
		desc = append(desc, val)
	}

	return english.Join(desc, ", ", " or ")
}

// getDistinctNames gathers all the qualified names with the same name and
// associates them with the list of Families they belong to. Then it
// generates a list of keys which is sorted so that the list with the most
// Families come first. duplicates are sorted by length of colour name.
func getDistinctNames(qNames []QualifiedColourName) (
	map[string][]Family, []string,
) {
	distinctNames := map[string][]Family{}
	for _, qn := range qNames {
		v := distinctNames[qn.ColourName]
		v = append(v, qn.Family)
		distinctNames[qn.ColourName] = v
	}

	for cName, families := range distinctNames {
		sort.Slice(families,
			func(i, j int) bool {
				familyNameI := families[i].String()
				familyNameJ := families[j].String()
				lenFamilyNameI := len(familyNameI)
				lenFamilyNameJ := len(familyNameJ)

				if lenFamilyNameI < lenFamilyNameJ {
					return true
				}

				if lenFamilyNameI == lenFamilyNameJ {
					return familyNameI < familyNameJ
				}

				return false
			})

		distinctNames[cName] = families
	}

	keys := slices.SortedFunc(maps.Keys(distinctNames), func(a, b string) int {
		familyCountA := len(distinctNames[a])
		familyCountB := len(distinctNames[b])

		if familyCountA == familyCountB {
			return cmp.Compare(len(a), len(b))
		}

		return cmp.Compare(familyCountB, familyCountA)
	})

	return distinctNames, keys
}

// getUniqueNames takes the set of colour names and removes synonyms -
// multiple names from the same Family for the same colour. It prefers
// shorter names over longer ones and, for X11 names, it prefers names
// without digits over names with (so black over grey0).
func getUniqueNames(qNames []QualifiedColourName) []QualifiedColourName {
	uniqueNames := map[Family]QualifiedColourName{}

	for _, qn := range qNames {
		crnt, ok := uniqueNames[qn.Family]
		if ok {
			qn = preferredName(qn, crnt)
		}

		uniqueNames[qn.Family] = qn
	}

	if len(uniqueNames) < len(qNames) {
		qNames = []QualifiedColourName{}
		for _, qn := range uniqueNames {
			qNames = append(qNames, qn)
		}
	}

	return qNames
}

// preferredName returns the preferred name. We prefer shorter names over
// longer; names without digits over names with; names coming earlier in
// lexical order over names coming later.
func preferredName(newQCN, currentQCN QualifiedColourName) QualifiedColourName {
	const digits = "1234567890"

	var (
		newName  = newQCN.ColourName
		crntName = currentQCN.ColourName
	)

	if len(crntName) < len(newName) {
		return currentQCN // prefer shorter names
	}

	if len(crntName) > len(newName) {
		return newQCN // prefer shorter names
	}

	var (
		newHasDigits  = strings.ContainsAny(newName, digits)
		crntHasDigits = strings.ContainsAny(crntName, digits)
	)

	if newHasDigits && !crntHasDigits {
		return currentQCN // prefer names without digits
	}

	if !newHasDigits && crntHasDigits {
		return newQCN // prefer names without digits
	}

	if crntName < newName {
		return currentQCN // prefer name coming 1st in lexical order
	}

	return newQCN
}
