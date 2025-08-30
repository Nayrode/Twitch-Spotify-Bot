// Standalone Spotify login script
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var (
	oauthConfig *oauth2.Config
	state       = "randomstate123"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	scopes := os.Getenv("SCOPES") // e.g. "user-modify-playback-state"

	redirectURL := "http://127.0.0.1:8080/callback"

	// Configure OAuth2
	oauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{scopes},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.spotify.com/authorize",
			TokenURL: "https://accounts.spotify.com/api/token",
		},
	}

	// Start local server for callback
	http.HandleFunc("/callback", handleCallback)

	authURL := oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Println("Go to the following link in your browser:\n", authURL)

	log.Println("Listening on :8080 ...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != state {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Token exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Access Token:", token.AccessToken)
	fmt.Println("Refresh Token:", token.RefreshToken)
	fmt.Println("Expiry:", token.Expiry)

	fmt.Fprintf(w, "Authorization complete! You can close this tab.")
}
