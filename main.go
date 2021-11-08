package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Guild struct {
		Id string `envconfig:"DISCORD_GUILD_ID"`
	}
	Discord struct {
		Token string `envconfig:"DISCORD_TOKEN"`
	}
}

var guildId string

func main() {

	// Load config and env variables
	cfg := setupConfig()
	guildId = cfg.Guild.Id

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(handleMessage)

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

	dg.Close()
}

func setupConfig() (cfg Config) {
	err := envconfig.Process("", &cfg)
	if err != nil {
		fmt.Println("error decoding env variables", err)
	}
	return cfg
}

//Callback for MessageCreate events
func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	HandleCommand(s, m, guildId, m.ChannelID, m.Content)
}
