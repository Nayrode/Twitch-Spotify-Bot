package observer

import (
	"fmt"
	"log"

	twitch "github.com/gempir/go-twitch-irc/v4"
)

// TwitchObserver observes messages in a Twitch chat channel.
type TwitchObserver struct {
	client  *twitch.Client
	channel string
}

// SendMessage sends a message to the Twitch chat channel.
func (o *TwitchObserver) SendMessage(message string) error {
	if o.client == nil || o.channel == "" {
		return fmt.Errorf("twitch client or channel not initialized")
	}
	o.client.Say(o.channel, message)
	return nil
}

// NewTwitchObserver creates a new TwitchObserver for the given channel.
func NewTwitchObserver(username, oauthToken, channel string) *TwitchObserver {
	client := twitch.NewClient(username, "oauth:"+oauthToken)
	return &TwitchObserver{
		client:  client,
		channel: channel,
	}
}

// Start begins observing the Twitch chat.
func (o *TwitchObserver) Start(onMessage func(user, message string)) error {
	o.client.OnPrivateMessage(func(msg twitch.PrivateMessage) {
		onMessage(msg.User.Name, msg.Message)
	})

	o.client.OnConnect(func() {
		log.Printf("Connected to Twitch chat. Joining #%s", o.channel)
		o.client.Join(o.channel)
	})

	err := o.client.Connect()
	if err != nil {
		return err
	}
	return nil
}

// Example usage:
// observer := NewTwitchObserver("your_username", "your_oauth_token", "channel_name")
// err := observer.Start(func(user, message string) {
//     log.Printf("[%s]: %s", user, message)
// })
// if err != nil {
//     log.Fatal(err)
// }