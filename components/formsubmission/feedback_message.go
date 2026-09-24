package formsubmission

const (
	userFeedbackMessage  = "Klarte ikke å lagre endringen. Prøv igjen. Kontakt styret dersom problemet vedvarer."
	adminFeedbackMessage = "Klarte ikke å lagre endringen. Prøv igjen. Sjekk logger dersom problemet vedvarer."
)

func FeedbackMessage(isAdmin bool) string {
	if isAdmin {
		return adminFeedbackMessage
	}
	return userFeedbackMessage
}
