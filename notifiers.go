package main

/*
Generate slack token: https://api.slack.com/web
*/

import (
	"fmt"
	"net/url"
	"os"

	pagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/bluele/slack"
	"github.com/fatih/color"
	"github.com/tbruyelle/hipchat-go/hipchat"
)

func (n *Notifier) Run(message string, isAlert bool) {
	switch n.Name {
	case "slack":
		if slack_api == nil {
			color.Red("[!] Slack used as a notifier, but not configured with ENV vars.")
			return
		}
		if isAlert {
			err = slack_api.ChatPostMessage(slack_channel.Id, message, &slack.ChatPostMessageOpt{IconEmoji: ":fire:"})
		} else {
			err = slack_api.ChatPostMessage(slack_channel.Id, message, &slack.ChatPostMessageOpt{IconEmoji: ":success:"})
		}
		if err != nil {
			color.Red(fmt.Sprintf("[!] Error posting to Slack: %s", err))
		}
	case "pagerduty":
		if pagerduty_api_token == "" || pagerduty_service_key == "" {
			color.Red("[!] PagerDuty used as a notifier, but not configured with ENV vars.")
		}

		if isAlert {
			event := pagerduty.Event{
				Type:        "trigger",
				ServiceKey:  pagerduty_service_key,
				Description: message,
			}
			resp, err := pagerduty.CreateEvent(event)
			if err != nil {
				fmt.Println(resp)
				color.Red(fmt.Sprintf("[!] Error posting to PagerDuty: %s", err))
			}
		} else {
			color.Green("[>] PagerDuty incident should be resolved now.")
		}

	case "hipchat":
		if hipchat_api == nil {
			color.Red("[!] HipChat used as a notifier, but not configured with ENV vars.")
			return
		}
		if isAlert {
			_, err = hipchat_api.Room.Notification(os.Getenv("HIPCHAT_ROOM_ID"), &hipchat.NotificationRequest{Message: message, Color: "red"})
		} else {
			_, err = hipchat_api.Room.Notification(os.Getenv("HIPCHAT_ROOM_ID"), &hipchat.NotificationRequest{Message: message, Color: "green"})
		}
		if err != nil {
			color.Red(fmt.Sprintf("[!] Error posting to HipChat: %s", err))
		}
	case n.Name: // default
		color.Yellow(fmt.Sprintf("[>] Unknown notifier: %s", n.Name))
	}

}

func setupPagerduty() {
	if len(os.Getenv("PAGERDUTY_API_TOKEN")) == 0 ||
		len(os.Getenv("PAGERDUTY_SERVICE_KEY")) == 0 {
		color.Yellow("[>] Skipping Pagerduty setup, missing PAGERDUTY_API_TOKEN and PAGERDUTY_SERVICE_KEY")
		return
	}

	pagerduty_api_token = os.Getenv("PAGERDUTY_API_TOKEN")
	pagerduty_service_key = os.Getenv("PAGERDUTY_SERVICE_KEY")
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
