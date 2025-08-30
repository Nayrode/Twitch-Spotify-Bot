package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// SearchSpotifyTrack searches for a track using the Spotify Web API and returns the URI of the first track found.
func SearchSpotifyTrack(query string, authToken string) (string, error) {
	url := "https://api.spotify.com/v1/search?q=" + strings.ReplaceAll(query, " ", "+") + "&type=track&limit=1"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("spotify API error %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Tracks struct {
			Items []struct {
				URI string `json:"uri"`
			} `json:"items"`
		} `json:"tracks"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}
	if len(result.Tracks.Items) == 0 {
		return "", fmt.Errorf("no tracks found for query: %s", query)
	}
	return result.Tracks.Items[0].URI, nil
}

// AddToSpotifyQueue adds a track or item to the user's Spotify queue using the Web API.
func AddToSpotifyQueue(uri string, authToken string) error {
	url := "https://api.spotify.com/v1/me/player/queue?uri=" + uri
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 && resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("spotify API error %d: %s", resp.StatusCode, string(body))
	}
	return nil
}
