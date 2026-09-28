package feedback

// Category is what a piece of feedback is about.
type Category string

const (
	CategoryWebsite    Category = "website"
	CategoryConvention Category = "convention"
	CategoryOther      Category = "other"
)

// Categories lists every valid category in display order.
var Categories = []Category{CategoryWebsite, CategoryConvention, CategoryOther}

// Valid reports whether c is one of the known categories.
func (c Category) Valid() bool {
	switch c {
	case CategoryWebsite, CategoryConvention, CategoryOther:
		return true
	}
	return false
}

// Label is the user-facing Norwegian name of the category.
func (c Category) Label() string {
	switch c {
	case CategoryWebsite:
		return "Nettsiden"
	case CategoryConvention:
		return "Festivalen"
	case CategoryOther:
		return "Noe annet"
	}
	return ""
}

// Slug is the Norwegian URL word for the category, used in ?om=<slug>.
func (c Category) Slug() string {
	switch c {
	case CategoryWebsite:
		return "nettsiden"
	case CategoryConvention:
		return "festivalen"
	case CategoryOther:
		return "annet"
	}
	return ""
}

// CategoryFromSlug maps a URL word back to its category.
func CategoryFromSlug(slug string) (Category, bool) {
	for _, category := range Categories {
		if category.Slug() == slug {
			return category, true
		}
	}
	return "", false
}
