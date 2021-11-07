package main

import (
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const MAX_MESSAGE_HISTORY_COUNT = 100
const BEST_POSTS_COUNT = 10

type Messages []*discordgo.Message

func HandleCommand(s *discordgo.Session, m *discordgo.MessageCreate, gid, cid, command string) {
	if command == "!bestposts" {
		var messages, _ = s.ChannelMessages(cid, MAX_MESSAGE_HISTORY_COUNT, "", "", "")

		filteredMessages := Messages(messages).extractTopPosts(BEST_POSTS_COUNT)
		messageContent := buildTopPostsDirectMessage(gid, cid, filteredMessages)
		dmChannel, _ := s.UserChannelCreate(m.Author.ID)
		s.ChannelMessageSend(dmChannel.ID, messageContent)
	}
}

func (msgs Messages) extractTopPosts(count int) []*discordgo.Message {
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

func buildTopPostsDirectMessage(guildId string, channelId string, msgs Messages) string {
	var sb strings.Builder
	sb.WriteString("Here are the top posts in the channel\n\n")
	for i, m := range msgs {
		sb.WriteString("#" + strconv.Itoa(i+1) + " by: " + m.Author.Username + " > ")
		sb.WriteString("https://discord.com/channels/" + guildId + "/" + channelId + "/" + m.ID)
		sb.WriteString("\n")
	}
	return sb.String()
}
