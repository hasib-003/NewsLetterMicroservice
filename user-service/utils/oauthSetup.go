package utils

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"log"
	"os"
)

func SetupOauth() {
	clientID := os.Getenv("OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("OAUTH_CLIENT_SECRET")
	callback := os.Getenv("OAUTH_CALLBACK")
	if clientID == "" || clientSecret == "" || callback == "" {
		log.Fatalf("Google auth credentials missing")
	}
	goth.UseProviders(
		google.New(clientID, clientSecret, callback))

	store := sessions.NewCookieStore([]byte("session-secret"))
	gothic.Store = store
}
