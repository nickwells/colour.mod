package colour

import (
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

// checkForDups reports any duplicate colours in the colours slice
func checkForDups(t *testing.T, id string, colours []rgba) {
	t.Helper()

	dupCount := map[rgba]int{}
	for _, c := range colours {
		dupCount[c]++
	}

	for c, cnt := range dupCount {
		if cnt > 1 {
			t.Log(id)
			t.Logf("\tcolour %#4.2v appears %d times", c, cnt)
			t.Error("\tDuplicate colour found")
		}
	}
}

func TestAllColours(t *testing.T) {
	var (
		webCount  = len(webColours)
		htmlCount = len(htmlColours)
	)

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		fl         Families
		maxColours int
	}{
		{
			ID:         testhelper.MkID("Web and HTML"),
			fl:         Families{WebColours, HTMLColours},
			maxColours: webCount + htmlCount,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			colours, err := tc.fl.AllColours()
			testhelper.CheckExpErr(t, err, tc)

			if len(colours) > tc.maxColours {
				t.Log(tc.IDStr())
				t.Logf("\tcolours returned: %d, max: %d",
					len(colours), tc.maxColours)
				t.Error("\ttoo many colours")
			}

			checkForDups(t, tc.IDStr(), colours)
		})
	}
}
