package formsubmission

const (
	userFeedbackMessage  = "Klarte ikkje å lagre endringa. Prøv igjen. Kontakt styret dersom problemet held fram."
	adminFeedbackMessage = "Klarte ikkje å lagre endringa. Prøv igjen. Sjekk logger dersom problemet held fram."
)

func FeedbackMessage(isAdmin bool) string {
	if isAdmin {
		return adminFeedbackMessage
	}
	return userFeedbackMessage
}
