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
		"`@karere attachment` to see a Slack attachment message.\n" +
		"`@karere help` to see this again."
)

var greetingPrefixes = []string{"Hi", "Hello", "Hey", "Kia ora"}

func main() {
	bot := slackbot.New(os.Getenv("SLACK_TOKEN"))

	toMe := bot.Messages(slackbot.DirectMessage, slackbot.DirectMention).Subrouter()

	hi := "hi|hello|hey|kia ora|tena koe"
	toMe.Hear(hi).MessageHandler(HelloHandler)
	toMe.Hear("help").MessageHandler(HelpHandler)
	toMe.Hear("attachment").MessageHandler(AttachmentsHandler)
	bot.Hear(`<@([a-zA-z0-9]+)?>`).MessageHandler(MentionHandler)
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
	txt := "Karere (literally, _messenger_) is the EDA Slack bot."
	attachment := slack.Attachment{
		Pretext:   "Karere. :sunglasses: :thumbsup:",
		Title:     "Keep track of your progress",
		TitleLink: "https://devacademy.co.nz",
		Text:      txt,
		Fallback:  txt,
		ImageURL:  "https://drive.google.com/open?id=0ByHQe6U7_e5aTHFWSk83WkZUMXM",
		Color:     "#7CD197",
	}

	// supports multiple attachments
	attachments := []slack.Attachment{attachment}
	bot.ReplyWithAttachments(evt, attachments, WithTyping)
}
