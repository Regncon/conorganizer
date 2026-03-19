package ticketholder

import "hash/fnv"

var accentColors = []string{
	"var(--color-accent-blue)",
	"var(--color-accent-purple)",
	"var(--color-accent-pink)",
	"var(--color-accent-red)",
	"var(--color-accent-orange)",
	"var(--color-accent-yellow)",
	"var(--color-accent-green)",
	"var(--color-accent-teal)",
	"var(--color-accent-cyan)",
}

func ColorForName(name string) string {
	h := fnv.New32a()
	h.Write([]byte(name))
	idx := int(h.Sum32()) % len(accentColors)
	if idx < 0 {
		idx += len(accentColors)
	}
	return accentColors[idx]
}
