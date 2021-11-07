package main

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Messages []*discordgo.Message

const OCM_DISCORD_GUILD_ID = "897425025827766322"
const MAX_MESSAGE_HISTORY_COUNT = 100
const BEST_POSTS_COUNT = 10

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Command for finding the best posts in the channel
	if m.Content == "!bestposts" {
		channelId := m.ChannelID
		var messages, _ = s.ChannelMessages(channelId, MAX_MESSAGE_HISTORY_COUNT, "", "", "")

		filteredMessages := Messages(messages).ExtractTopPosts(BEST_POSTS_COUNT)
		messageContent := BuildTopPostsDirectMessage(channelId, filteredMessages)
		dmChannel, _ := s.UserChannelCreate(m.Author.ID)
		s.ChannelMessageSend(dmChannel.ID, messageContent)
	}

}

func (msgs Messages) ExtractTopPosts(count int) []*discordgo.Message {
	var msgsReactions []*discordgo.Message
	for _, m := range msgs {
		if len(m.Reactions) > 0 {
			msgsReactions = append(msgsReactions, m)
		}
	}

	sort.Slice(msgsReactions, func(i, j int) bool {
		return len(msgsReactions[i].Reactions) > len(msgsReactions[j].Reactions)
	})

	if len(msgsReactions) > count {
		msgsReactions = msgsReactions[0:count]
	}

	return msgsReactions
}

func BuildTopPostsDirectMessage(channelId string, msgs Messages) string {
	var sb strings.Builder
	sb.WriteString("Here are the top posts in the channel\n\n")
	for i, m := range msgs {
		sb.WriteString("#" + strconv.Itoa(i+1) + " by: " + m.Author.Username + " > ")
		sb.WriteString("https://discord.com/channels/" + OCM_DISCORD_GUILD_ID + "/" + channelId + "/" + m.ID)
		sb.WriteString("\n")
	}
	return sb.String()
}
