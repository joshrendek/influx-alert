package main

/*
Generate slack token: https://api.slack.com/web
*/

import (
	"fmt"
	"github.com/bluele/slack"
	"github.com/fatih/color"
	"github.com/tbruyelle/hipchat-go/hipchat"
	"net/url"
	"os"
)

func (n *Notifier) Run(message string) {
	switch n.Name {
	case "slack":
		if slack_api == nil {
			color.Red("[!] Slack used as a notifier, but not configured with ENV vars.")
			return
		}
		err = slack_api.ChatPostMessage(slack_channel.Id, message, &slack.ChatPostMessageOpt{IconEmoji: ":fire:"})
		if err != nil {
			color.Red(fmt.Sprintf("[!] Error posting to Slack: %s", err))
		}
	case "hipchat":
		if hipchat_api == nil {
			color.Red("[!] HipChat used as a notifier, but not configured with ENV vars.")
			return
		}
		_, err = hipchat_api.Room.Notification(os.Getenv("HIPCHAT_ROOM_ID"), &hipchat.NotificationRequest{Message: "Testing", Color: "red"})
		if err != nil {
			color.Red(fmt.Sprintf("[!] Error posting to HipChat: %s", err))
		}
	case n.Name: // default
		color.Yellow(fmt.Sprintf("[>] Unknown notifier: %s", n.Name))
	}

}

func setupSlack() {
	if len(os.Getenv("SLACK_API_TOKEN")) == 0 ||
		len(os.Getenv("SLACK_ROOM")) == 0 {
		color.Yellow("[>] Skipping Slack setup, missing SLACK_API_TOKEN and SLACK_ROOM")
		return
	}
	slack_api = slack.New(os.Getenv("SLACK_API_TOKEN"))

	slack_channel, err = slack_api.FindChannelByName(os.Getenv("SLACK_ROOM"))
	if err != nil {
		panic(err)
	}
}

func setupHipchat() {
	hipchat_api = hipchat.NewClient(os.Getenv("HIPCHAT_API_TOKEN"))
	if len(os.Getenv("HIPCHAT_API_TOKEN")) == 0 ||
		len(os.Getenv("HIPCHAT_ROOM_ID")) == 0 {
		color.Yellow("[>] Skipping Hipchat setup, missing HIPCHAT_API_TOKEN and HIPCHAT_ROOM_ID")
		return
	}
	if os.Getenv("HIPCHAT_SERVER") != "" {
		hipchat_api.BaseURL, err = url.Parse(os.Getenv("HIPCHAT_SERVER"))
		if err != nil {
			color.Red("Error connecting to private hipchat server: ", err)
			panic(err)
		}
	}
}
