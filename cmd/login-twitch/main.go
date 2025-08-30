// Standalone Twitch login script
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
	twitchOauthConfig *oauth2.Config
	twitchState       = "randomtwitchstate123"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	scopes := os.Getenv("TWITCH_SCOPES")

	redirectURL := "http://localhost:8081/twitch_callback"

	twitchOauthConfig = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{scopes},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://id.twitch.tv/oauth2/authorize",
			TokenURL: "https://id.twitch.tv/oauth2/token",
		},
	}

	http.HandleFunc("/twitch_callback", handleTwitchCallback)

	authURL := twitchOauthConfig.AuthCodeURL(twitchState, oauth2.AccessTypeOffline)
	fmt.Println("Go to the following link in your browser for Twitch login:\n", authURL)

	log.Println("Listening on :8081 ...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleTwitchCallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != twitchState {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := twitchOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Token exchange failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Access Token:", token.AccessToken)
	fmt.Println("Refresh Token:", token.RefreshToken)
	fmt.Println("Expiry:", token.Expiry)

	fmt.Fprintf(w, "Twitch authorization complete! You can close this tab.")
}
