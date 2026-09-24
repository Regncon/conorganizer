package components

import (
	"encoding/json"

	"github.com/a-h/templ"
)

func KVPairsAttrs(kvPairs ...string) templ.Attributes {
	if len(kvPairs)%2 != 0 {
		panic("kvPairs must be a multiple of 2")
	}
	attrs := templ.Attributes{}
	for i := 0; i < len(kvPairs); i += 2 {
		attrs[kvPairs[i]] = kvPairs[i+1]
	}
	return attrs
}

// DataSignals renders several Datastar signals as one data-signals attribute.
// Prefer it over a stack of data-signals:<name> attributes: values are encoded
// as JSON, so quoting and escaping are handled for you.
//
//	<dialog data-signals={ components.DataSignals(map[string]any{
//		"puljeId":     puljeId,
//		"puljeCanEdit": state.CanEdit,
//	}) }>
//
// Keys are used verbatim, unlike the data-signals:<name> form, which converts
// kebab-case to camelCase. Write the signal name as the browser reads it:
// "puljeCanEdit" for $puljeCanEdit. Nested maps become grouped signals, so
// {"crop": map[string]any{"ready": false}} is read as $crop.ready.
// A value that cannot be marshalled yields "{}", leaving the signals unset.
func DataSignals(values any) string {
	b, err := json.Marshal(values)
	if err != nil {
		return "{}"
	}
	return string(b)
}
