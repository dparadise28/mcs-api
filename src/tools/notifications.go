package tools

import (
	"github.com/nlopes/slack"
	"log"
)

var (
	SLACK_TOKEN = ""
	MODE        = ""
)

func SendToSlack(destinationChannel, pretext, text string) {
	log.Println(SLACK_TOKEN)
	if len(SLACK_TOKEN) > 0 {
		api := slack.New(SLACK_TOKEN)
		channels, err := api.GetChannels(false)
		log.Println(err)
		/*
			// get back to later (log errors and monitor logs)
			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
		*/
		for _, channel := range channels {
			log.Println(channel.Name)
			if channel.Name == destinationChannel {
				params := slack.PostMessageParameters{}
				attachment := slack.Attachment{
					Pretext: pretext,
					Text:    text,
				}
				params.Attachments = []slack.Attachment{attachment}
				// channelID, timestamp, err (should track these; come back to later)
				_, _, err := api.PostMessage(channel.ID, text, params)
				log.Println(err)
				return
			}
		}
	}
}
