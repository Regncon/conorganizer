package packagefailure

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("intentional package-level failure for reporter demo")
	os.Exit(1)
}
