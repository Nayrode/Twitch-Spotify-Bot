package main

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"

	handler "twitch-spotify-bot/src/handler"
	"twitch-spotify-bot/src/observer"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables.")
	}

	twitchUsername := os.Getenv("TWITCH_USERNAME")
	twitchRefreshToken := os.Getenv("TWITCH_OAUTH")
	twitchChannel := os.Getenv("TWITCH_CHANNEL")
	spotifyRefreshToken := os.Getenv("SPOTIFY_TOKEN")

	if twitchUsername == "" || twitchRefreshToken == "" || twitchChannel == "" || spotifyRefreshToken == "" {
		log.Fatal("Missing required environment variables: TWITCH_USERNAME, TWITCH_OAUTH, TWITCH_CHANNEL, SPOTIFY_TOKEN")
	}

	twitchAccessToken, err := handler.GetTwitchAccessToken(twitchRefreshToken)
	if err != nil {
		log.Fatalf("Could not refresh Twitch access token: %v", err)
	}

	obs := observer.NewTwitchObserver(twitchUsername, twitchAccessToken, twitchChannel)

	twitchClientID := os.Getenv("TWITCH_CLIENT_ID")
	err = obs.Start(func(user, message string) {
		if strings.HasPrefix(message, "!sr ") {
			// Check if channel is live
			isLive, liveErr := handler.IsTwitchChannelLive(twitchChannel, twitchAccessToken, twitchClientID)
			if liveErr != nil {
				log.Printf("[ERROR] Could not check if channel is live: %v", liveErr)
				return
			}
			if !isLive {
				log.Printf("[INFO] Ignoring request, channel is not live.")
				return
			}

			query := strings.TrimSpace(strings.TrimPrefix(message, "!sr "))
			if query == "" {
				log.Printf("[WARN] Empty search request from %s", user)
				return
			}

			// Remove URLs from query (if user pastes a link)
			urlRegex := regexp.MustCompile(`https?://\\S+`)
			query = urlRegex.ReplaceAllString(query, "")
			query = strings.TrimSpace(query)

			accessToken, err := handler.GetSpotifyAccessToken(spotifyRefreshToken)
			if err != nil {
				log.Printf("[ERROR] Could not refresh Spotify access token: %v", err)
				return
			}

			uri, err := handler.SearchSpotifyTrack(query, accessToken)
			if err != nil {
				log.Printf("[ERROR] Spotify search failed for '%s': %v", query, err)
				return
			}

			err = handler.AddToSpotifyQueue(uri, accessToken)
			if err != nil {
				log.Printf("[ERROR] Failed to add to queue: %v", err)
				return
			}

			log.Printf("[INFO] Added to queue: %s (requested by %s)", uri, user)
			// Send a message in Twitch chat
			msg := "Added to queue: " + query + " (requested by " + user + ")"
			if sendErr := obs.SendMessage(msg); sendErr != nil {
				log.Printf("[WARN] Could not send Twitch chat message: %v", sendErr)
			}
		}
	})
	if err != nil {
		log.Fatalf("Twitch observer error: %v", err)
	}
}
