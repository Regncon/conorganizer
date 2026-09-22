package formsubmission

import "github.com/Regncon/conorganizer/components"

func dataSignals(values any) string {
	return DataSignals(values)
}

// DataSignals is kept as an alias for components.DataSignals so existing
// call sites in this package keep working. New code outside this package
// should call components.DataSignals directly.
func DataSignals(values any) string {
	return components.DataSignals(values)
}
