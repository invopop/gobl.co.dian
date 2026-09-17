package dian_test

import (
	"flag"
	"testing"

	// Register the Colombia DIAN addon so example documents declaring the
	// co-dian-v2 addon normalize and validate.
	_ "github.com/invopop/gobl.co.dian/addon"

	"github.com/invopop/gobl/pkg/examples"
)

var update = flag.Bool("update", false, "update the example golden files")

// TestExamples converts every document under examples/ to a calculated,
// validated JSON envelope and compares it against its golden output, using the
// shared GOBL example helpers. Run with -update to (re)generate the goldens.
func TestExamples(t *testing.T) {
	examples.Run(t, "examples", *update)
}
