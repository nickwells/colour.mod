package colour

import (
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestContrast(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		c rgba
	}{
		{
			ID: testhelper.MkID("black"),
			c:  rgba{A: 0xff},
		},
		{
			ID: testhelper.MkID("white"),
			c:  rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		},
		{
			ID: testhelper.MkID("white"),
			c:  rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		},
		{
			ID: testhelper.MkID("grey"),
			c:  rgba{R: 0x80, G: 0x80, B: 0x80, A: 0xff},
		},
		{
			ID: testhelper.MkID("red"),
			c:  rgba{R: 0xff, A: 0xff},
		},
		{
			ID: testhelper.MkID("green"),
			c:  rgba{G: 0xff, A: 0xff},
		},
		{
			ID: testhelper.MkID("blue"),
			c:  rgba{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
		},
		{
			ID: testhelper.MkID("cyan"),
			c:  rgba{G: 0xff, B: 0xff, A: 0xff},
		},
		{
			ID: testhelper.MkID("magenta"),
			c:  rgba{R: 0xff, B: 0xff, A: 0xff},
		},
		{
			ID: testhelper.MkID("yellow"),
			c:  rgba{R: 0xff, G: 0xff, A: 0xff},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			for _, generator := range []struct {
				name string
				f    func(rgba) rgba
			}{
				{name: "Contrast", f: Contrast},
				{name: "ContrastColourful", f: ContrastColourful},
			} {
				cc := generator.f(tc.c)

				if rc, rcc := Roughly(tc.c), Roughly(cc); rc == rcc {
					t.Log(tc.IDStr() + " : " + generator.name)
					t.Logf("\t:          colour %#v (roughly %s)", tc.c, rc)
					t.Logf("\t: contrast colour %#v (roughly %s)", cc, rcc)
					t.Error("\t: bad contrast")
				}
			}
		})
	}
}
