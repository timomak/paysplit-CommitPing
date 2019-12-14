package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/nlopes/slack"
	"gopkg.in/go-playground/webhooks.v5/github"

	_ "github.com/joho/godotenv/autoload"
)

// CreateSlackClient Create a client for the slack API
func CreateSlackClient(apiKey string) *slack.RTM {
	api := slack.New(apiKey)
	rtm := api.NewRTM()
	go rtm.ManageConnection() // goroutine!
	return rtm
}

// NotifySlackChannel sends a message to a Slack Channel using the Slack API
func NotifySlackChannel(slackClient *slack.RTM, message, channel string) {
	slackMsg := slack.MsgOptionText(message, false) // Not sure why the false.
	slackClient.PostMessage(channel, slackMsg)      // Channel name, message
}

// main is our entrypoint, where the application initializes the Slackbot.
func main() {
	hook, _ := github.New(github.Options.Secret(string(os.Getenv("WEBHOOK")))) // Secret for Webhook.
	e := echo.New()
	e.POST("/push", func(c echo.Context) error {
		fmt.Println("PUSH route called.")
		payload, err := hook.Parse(c.Request(), github.PushEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn't one of the ones asked to be parsed
			}
		}
		switch payload.(type) {

		case github.PushPayload:
			fmt.Println("Message received")
			// newMessage := "A commit has just been made to tsukudabuddha/paysplit"
			release := payload.(github.PushPayload)

			newMessage := string(release.Pusher.Name) + ": " + string(release.HeadCommit.Message) + "\nRepo: " + string(release.Repository.FullName) + "\nURL: " + string(release.Repository.HTMLURL)

			// fmt.Println("Release:", release)
			fmt.Printf("%+v", release.Ref)
			if strings.Contains(release.Ref, "master") == true {
				fmt.Printf("Is Master")
			}
			slackIt(newMessage, "paysplit-devs") // Message, Channel Name
		}

		return c.String(http.StatusOK, "Success.")
	})
	e.Logger.Fatal(e.Start(":3000"))
}

// slackIt is a function that initializes the Slackbot and sends a custom message to a specific channel.
func slackIt(message, channel string) {
	botToken := os.Getenv("BOT_OAUTH_ACCESS_TOKEN")
	slackClient := CreateSlackClient(botToken)
	// fmt.Println("SENDING MESSASSAGE TO SLACK CHANNEL:", message)
	NotifySlackChannel(slackClient, message, channel)
}

// Finished.
