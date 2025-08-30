# Twitch-Spotify-Bot

This bot allows Twitch chat users to request songs via Spotify using the `!sr` command. When a song is requested, it is added to the Spotify playback queue and a confirmation message is sent in the Twitch chat.

## Setup Instructions

1. **Clone the repository and install dependencies:**
	 ```sh
	 git clone <repo-url>
	 cd Twitch-Spotify-Bot
	 go mod tidy
	 ```

2. **Obtain your Spotify and Twitch tokens:**
	 - To get your Spotify tokens, run:
		 ```sh
		 go run ./cmd/login-spotify
		 ```
	 - To get your Twitch tokens, run:
		 ```sh
		 go run ./cmd/login-twitch
		 ```
	 - Follow the instructions in the terminal and browser. Copy the generated refresh tokens.

3. **Configure your environment variables:**
	 - Copy `.env.template` to `.env` and fill in the required values:
		 - Spotify: `CLIENT_ID`, `CLIENT_SECRET`, `SCOPES`, `SPOTIFY_TOKEN` (refresh token)
		 - Twitch: `TWITCH_CLIENT_ID`, `TWITCH_CLIENT_SECRET`, `TWITCH_SCOPES`, `TWITCH_USERNAME`, `TWITCH_OAUTH` (refresh token), `TWITCH_CHANNEL`

4. **Run the bot:**
	 ```sh
	 go run ./src/main.go
	 ```

## Usage
- In your Twitch chat, type:
	```
	!sr <song name or artist>
	```
- The bot will search Spotify, add the first matching track to the queue, and confirm in chat.

---

**Note:**
- Make sure your Spotify account is actively playing on a device for the queue to work.
- The bot will automatically refresh both Spotify and Twitch tokens as needed.