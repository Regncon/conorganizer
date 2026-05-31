package formsubmission

import "encoding/json"

func dataSignals(values any) string {
	return DataSignals(values)
}

func DataSignals(values any) string {
	b, err := json.Marshal(values)
	if err != nil {
		return "{}"
	}
	return string(b)
}
