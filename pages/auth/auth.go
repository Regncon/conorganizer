package auth

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/descope/go-sdk/descope/client"
	"github.com/go-chi/chi/v5"
)

func SetupAuthRoute(router chi.Router, logger *slog.Logger) error {
	projectID := "P2ufzqahlYUHDIprVXtkuCx8MH5C"
	descopeClient, err := client.NewWithConfig(&client.Config{ProjectID: projectID})
	if err != nil {
		logger.Error("Failed to create Descope client", slog.String("projectID", projectID), err)
		return err
	}
	if descopeClient != nil {
		logger.Info("Descope client created successfully")
	}
	// extract from request authorization header. The above sample code sends the the session token in authorization header.
	// sessionToken := r.Header.Get("Authorization")

	router.Route("/auth", func(authRouter chi.Router) {
		authRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			loginForm().Render(r.Context(), w)
		})

		authRouter.Get("/test", func(w http.ResponseWriter, r *http.Request) {

			cookieSession, err := r.Cookie("session_token")
			if err != nil {
				if err == http.ErrNoCookie {
					http.Error(w, "Cookie ikke funnet", http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, "Cookie value: %s", cookieSession.Value)

			cookieRefresh, err := r.Cookie("session_token")
			if err != nil {
				if err == http.ErrNoCookie {
					http.Error(w, "Cookie ikke funnet", http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, "Cookie value: %s", cookieRefresh.Value)

			authorized, userToken, err := descopeClient.Auth.ValidateAndRefreshSessionWithTokens(r.Context(), cookieSession.Value, cookieRefresh.Value)

			cookie := &http.Cookie{
				Name:    "session_token",
				Value:   userToken.JWT,
				Path:    "/",
				Expires: time.Unix(userToken.Expiration, 0),
			}
			// Skriv cookien til responsen
			http.SetCookie(w, cookie)
			if err != nil {
				fmt.Println("Could not validate user session: ", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			} else {
				fmt.Printf("Successfully validated user session: %v\n", userToken.Claims["email"])
			}
			w.Write([]byte(fmt.Sprintf("test, is authrs: %v", authorized)))
		})

	})

	return nil
}
