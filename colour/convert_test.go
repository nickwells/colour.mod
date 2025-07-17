package colour_test

import (
	"testing"

	"github.com/nickwells/colour.mod/v2/colour"
)

func TestHSL(t *testing.T) {
	colours, err := colour.MakeColours(9)
	if err != nil {
		t.Fatal("couldn't generate the colours:", err)
		return
	}

	for i, c := range colours {
		hsl, _ := colour.RGBA2HSLAndHSV(c)
		cFromHSL := hsl.ToRGBA()

		if err := colour.Compare(cFromHSL, c, 1); err != nil {
			t.Logf("colour: %d", i)
			t.Logf("conversion to and from HSL of %+v", c)
			t.Logf("                    generated %+v", cFromHSL)
			t.Error("\t", err)
		}
	}
}
