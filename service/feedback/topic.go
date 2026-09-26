package feedback

// Topic is an optional tag narrowing down what feedback in a category is about.
type Topic string

const (
	TopicSignup  Topic = "signup"
	TopicProgram Topic = "program"
	TopicMyPage  Topic = "my-page"
	TopicMobile  Topic = "mobile"

	TopicEvents      Topic = "events"
	TopicVenue       Topic = "venue"
	TopicFood        Topic = "food"
	TopicInformation Topic = "information"
)

var topicsByCategory = map[Category][]Topic{
	CategoryWebsite:    {TopicSignup, TopicProgram, TopicMyPage, TopicMobile},
	CategoryConvention: {TopicEvents, TopicVenue, TopicFood, TopicInformation},
	CategoryOther:      {},
}

// TopicsFor lists the topics that can be picked for a category, in display order.
// It returns an empty slice for categories without topics or unknown categories.
func TopicsFor(c Category) []Topic {
	topics := topicsByCategory[c]
	result := make([]Topic, len(topics))
	copy(result, topics)
	return result
}

// Label is the user-facing Norwegian name of the topic.
func (t Topic) Label() string {
	switch t {
	case TopicSignup:
		return "Påmelding"
	case TopicProgram:
		return "Program"
	case TopicMyPage:
		return "Min side"
	case TopicMobile:
		return "Mobil"
	case TopicEvents:
		return "Arrangementer"
	case TopicVenue:
		return "Lokaler"
	case TopicFood:
		return "Mat"
	case TopicInformation:
		return "Informasjon"
	}
	return ""
}

func (t Topic) belongsTo(c Category) bool {
	for _, topic := range topicsByCategory[c] {
		if topic == t {
			return true
		}
	}
	return false
}
