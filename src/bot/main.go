package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	slackbot "github.com/BeepBoopHQ/go-slackbot"
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
)

const (
	WithTyping    = slackbot.WithTyping
	WithoutTyping = slackbot.WithoutTyping

	HelpText = "I will respond to the following commands:\n" +
		"`@karere hi` for a simple message.\n" +
		"`@karere attachment` to see a Slack attachment message.\n" +
		"`@karere hi` to demonstrate detecting a mention.\n" +
		"`@karere help` to see this again."
)

var greetingPrefixes = []string{"Hi", "Hello", "Hey", "Kia ora"}

func main() {
	bot := slackbot.New(os.Getenv("SLACK_TOKEN"))

	toMe := bot.Messages(slackbot.DirectMessage, slackbot.DirectMention).Subrouter()

	hi := "karere hi|karere hello"
	toMe.Hear(hi).MessageHandler(HelloHandler)
	toMe.Hear("help").MessageHandler(HelpHandler)
	toMe.Hear("attachment").MessageHandler(AttachmentsHandler)
	toMe.Hear(`<@([a-zA-z0-9]+)?>`).MessageHandler(MentionHandler)
	toMe.Hear("(karere ).*").MessageHandler(CatchAllHandler)
	bot.Run()
}

func HelloHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	rand.Seed(time.Now().UnixNano())
	msg := greetingPrefixes[rand.Intn(len(greetingPrefixes))] + " <@" + evt.User + ">!"
	bot.Reply(evt, msg, WithTyping)

	if slackbot.IsDirectMessage(evt) {
		dmMsg := "It's nice to talk to you directly."
		bot.Reply(evt, dmMsg, WithoutTyping)
	}

	bot.Reply(evt, "If you'd like to talk some more, "+HelpText, WithTyping)
}

func CatchAllHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	msg := fmt.Sprintf("I'm sorry, I don't know how to: `%s`.\n%s", evt.Text, HelpText)
	bot.Reply(evt, msg, WithTyping)
}

func MentionHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	if slackbot.IsMentioned(evt, bot.BotUserID()) {
		bot.Reply(evt, "You really do care about me. :heart:", WithTyping)
	}
}

func HelpHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	bot.Reply(evt, HelpText, WithTyping)
}

func AttachmentsHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	txt := "Beep Beep Boop is a ridiculously simple hosting platform for your Slackbots."
	attachment := slack.Attachment{
		Pretext:   "We bring bots to life. :sunglasses: :thumbsup:",
		Title:     "Host, deploy and share your bot in seconds.",
		TitleLink: "https://beepboophq.com/",
		Text:      txt,
		Fallback:  txt,
		ImageURL:  "https://storage.googleapis.com/beepboophq/_assets/bot-1.22f6fb.png",
		Color:     "#7CD197",
	}

	// supports multiple attachments
	attachments := []slack.Attachment{attachment}
	bot.ReplyWithAttachments(evt, attachments, WithTyping)
}
