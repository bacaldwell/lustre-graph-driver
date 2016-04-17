// +build linux

package lustre

import (
	"github.com/bacaldwell/lustre-graph-driver/driver/graphtest"
	"testing"
)

// This avoids creating a new driver for each test if all tests are run
// Make sure to put new tests between TestLustreSetup and TestLustreTeardown
func TestLustreSetup(t *testing.T) {
	graphtest.GetDriver(t, "Lustre")
}

func TestLustreCreateEmpty(t *testing.T) {
	graphtest.DriverTestCreateEmpty(t, "Lustre")
}

func TestLustreCreateBase(t *testing.T) {
	graphtest.DriverTestCreateBase(t, "Lustre")
}

func TestLustreCreateSnap(t *testing.T) {
	graphtest.DriverTestCreateSnap(t, "Lustre")
}

func TestLustreTeardown(t *testing.T) {
	graphtest.PutDriver(t)
}
