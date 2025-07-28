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
		t.Logf("%#v", source[scn.Start:scn.End])
	}
}
