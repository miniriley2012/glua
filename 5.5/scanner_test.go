package glua

import (
	"strings"
	"testing"
)

func TestScan(t *testing.T) {
	const source = "const thing"

	scn := New(strings.NewReader(source))

	for scn.Scan() {
		if err := scn.Err(); err != nil {
			t.Fatalf("error: %v", err)
		}
		t.Logf("%s: %#v", scn.Token, source[scn.Start:scn.End])
	}
}
