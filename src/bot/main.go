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

	HelpText = "Karere uses git-style command syntax:\n" +
		"\u2022 `about`: access _all_, implemented _yes_.\n" +
		"\u0149 `add <github-user> <slack-user> [cohort-repo]`: access _restrict_, implemented _no_.\n" +
		"* `blocks [cohort-repo]`: access _restrict_, implemented _no_.\n" +
		"* `gist [add|rm] [gist-name]`: access _all_, implemented _no_.\n" +
		"* `help`: access _all_, implemented _partial_.\n" +
		"* `init [cohort-name]`: access _restrict_, implemented _no_.\n" +
		"* `log [cohort-name]`: access _restrict_, implemented _no_.\n" +
		"* `mv <github-user> <old-cohort-repo> <new-cohort-repo>: access _restrict_, implemented _no_.`\n" +
		"* `push u2w3 [cohort-repo]`: access _restrict_, implemented _no_.\n" +
		"* `reset u1w1 [cohort-repo]`: access _restrict_, implemented _no_.\n" +
		"* `rm <github-user> [cohort-repo]`: access _restrict_, implemented _no_.\n" +
		"* `snip [add|rm] [snippet-name]`: access _all_, implemented _no_.\n" +
		"* `version`: access _restrict_, implemented _no_.\n" +
		"In all cases where _cohort-repo_ is omitted, Karere will attempt to use the Slack channel name instead."
)

var greetingPrefixes = []string{"Hi", "Hello", "Hey", "Kia ora"}

func main() {
	bot := slackbot.New(os.Getenv("SLACK_TOKEN"))

	toMe := bot.Messages(slackbot.DirectMessage, slackbot.DirectMention).Subrouter()

	hi := "hi|hello|hey|kia ora|tena koe"
	toMe.Hear(hi).MessageHandler(HelloHandler)
	toMe.Hear("help").MessageHandler(HelpHandler)
	toMe.Hear("about").MessageHandler(AttachmentsHandler)
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
	msg := fmt.Sprintf("I'm sorry, I don't know how to: `%s`.\n", evt.Text)
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
	txt := `Karere can help keep you up to date with your course work, and lets us know when you need more help.
Licensed under the AGPL v3: http://www.gnu.org/licenses/agpl-3.0.html .
GitHub: https://github.com/richchurcher/karere`
	fields := []slack.AttachmentField{
		slack.AttachmentField{
			Title: "Foo",
			Value: "Bar",
			Short: true,
		},
	}

	attachment := slack.Attachment{
		Pretext:   "Karere (_messenger_) is the EDA Slack bot.",
		Title:     "Keep track of your progress",
		TitleLink: "https://devacademy.co.nz",
		Text:      txt,
		Fallback:  txt,
		ImageURL:  "http://i.imgur.com/4PA5eqt.jpg",
		Color:     "#7CD197",
		Fields:    fields,
	}

	// supports multiple attachments
	attachments := []slack.Attachment{attachment}
	bot.ReplyWithAttachments(evt, attachments, WithTyping)
}
