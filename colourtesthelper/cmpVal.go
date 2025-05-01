package colourtesthelper

import (
	"image/color" //nolint:misspell
	"strings"
	"testing"
)

// reportRGBADiff reports the difference between the two values
func reportRGBADiff(t *testing.T, prefix, part string, exp, act uint8) {
	t.Helper()

	if exp == act {
		return
	}

	t.Logf("\t: %s:   expected %s: 0x%02x\n", prefix, part, exp)
	t.Logf("\t: %s:     actual %s: 0x%02x\n", prefix, part, act)
}

// DiffRGBA compares the actual against the expected value and reports an
// error if they differ.
func DiffRGBA(t *testing.T,
	id, name string,
	act, exp color.RGBA, //nolint:misspell
) bool {
	t.Helper()

	if act != exp {
		t.Log(id)
		t.Logf("\t: expected %s: %#v\n", name, exp)
		t.Logf("\t:   actual %s: %#v\n", name, act)

		charCnt := len(name) + len("expected") + 1
		t.Logf("\t: %*s:\n", charCnt, "diffs")

		intro := strings.Repeat(" ", charCnt)
		reportRGBADiff(t, intro, "  red", exp.R, act.R)
		reportRGBADiff(t, intro, "green", exp.G, act.G)
		reportRGBADiff(t, intro, " blue", exp.B, act.B)
		reportRGBADiff(t, intro, "alpha", exp.A, act.A)

		t.Errorf("\t: %s is incorrect\n", name)

		return true
	}

	return false
}
