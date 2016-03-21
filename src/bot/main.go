package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	slackbot "github.com/BeepBoopHQ/go-slackbot"
	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

const (
	WithTyping    = slackbot.WithTyping
	WithoutTyping = slackbot.WithoutTyping

	HelpText = "Karere uses git-style command syntax:\n" +
		"\u2022 `about`: access _all_, implemented _yes_.\n" +
		"\u2022 `add <github-user> <slack-user> [cohort-repo]`: access _restrict_, implemented _no_.\n" +
		"\u2022 `blocks [cohort-repo]`: access _restrict_, implemented _no_.\n" +
		"\u2022 `gist [add|rm] [gist-name]`: access _all_, implemented _no_.\n" +
		"\u2022 `help`: access _all_, implemented _partial_.\n" +
		"\u2022 `init [cohort-name]`: access _restrict_, implemented _no_.\n" +
		"\u2022 `log [cohort-name]`: access _restrict_, implemented _no_.\n" +
		"\u2022 `mv <github-user> <old-cohort-repo> <new-cohort-repo>: access _restrict_, implemented _no_.`\n" +
		"\u2022 `push u2w3 [cohort-repo]`: access _restrict_, implemented _no_.\n" +
		"\u2022 `reset u1w1 [cohort-repo]`: access _restrict_, implemented _no_.\n" +
		"\u2022 `rm <github-user> [cohort-repo]`: access _restrict_, implemented _no_.\n" +
		"\u2022 `snip [add|rm] [snippet-name]`: access _all_, implemented _no_.\n" +
		"\u2022 `version`: access _restrict_, implemented _no_.\n" +
		"In all cases where _cohort-repo_ is omitted, Karere will attempt to use the Slack channel name instead."

	token = os.Getenv("GITHUB_ACCESS_TOKEN")
)

var greetingPrefixes = []string{"Hi", "Hello", "Hey", "Kia ora"}

func main() {
	bot := slackbot.New(os.Getenv("SLACK_TOKEN"))
	client := github.NewClient(nil)

	toMe := bot.Messages(slackbot.DirectMessage, slackbot.DirectMention).Subrouter()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)
	// list all repositories for the authenticated user
	//repos, _, err := client.Repositories.List("", nil)

	hi := "hi|hello|hey|kia ora|tena koe"
	toMe.Hear(hi).MessageHandler(HelloHandler)
	toMe.Hear("help").MessageHandler(HelpHandler)
	toMe.Hear("about").MessageHandler(AboutHandler)
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

func AboutHandler(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	txt := `Karere can help keep you up to date with your course work, and lets us know when you need more help.
Licensed under the AGPL v3: http://www.gnu.org/licenses/agpl-3.0.html .
GitHub: https://github.com/richchurcher/karere`
	attachment := slack.Attachment{
		Pretext:   "Karere (_messenger_) is the EDA Slack bot.",
		Title:     "Keep track of your progress",
		TitleLink: "https://devacademy.co.nz",
		Text:      txt,
		Fallback:  txt,
		ImageURL:  "http://i.imgur.com/4PA5eqt.jpg",
		Color:     "#7CD197",
	}

	attachments := []slack.Attachment{attachment}
	bot.ReplyWithAttachments(evt, attachments, WithTyping)
}
