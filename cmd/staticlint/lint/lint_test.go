package lint

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestCheckOSExit(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), CheckOSExit, "./...")
}
