package colour

import (
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestNamesByIndex(t *testing.T) {
	colourCount := 0
	for _, f := range searchOrder {
		colourCount += len(cFamMap[f])
	}
	nbi := NamesByIndex()
	testhelper.DiffInt[int](t, "NamesByIndex", "colour count",
		len(nbi), colourCount)
}
