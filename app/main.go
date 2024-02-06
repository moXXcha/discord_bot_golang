package main

import (
	main_internal "invite/internal"
	"fmt"

	"os"

	"github.com/bwmarrin/discordgo"

	"github.com/joho/godotenv"
	"os/signal"
    "syscall"
)

type Session interface {
	ChannelMessageSendEmbed(channelID string, embed *discordgo.MessageEmbed) (*discordgo.Message, error)
	GuildChannels(guildID string) (st []*discordgo.Channel, err error)
	GuildMembers(guildID string, after string, limit int) (st []*discordgo.Member, err error)
	Channel(channelID string) (st *discordgo.Channel, err error)
	ChannelInviteCreate(channelID string, i discordgo.Invite) (st *discordgo.Invite, err error)
	ChannelMessageSend(channelID string, content string) (*discordgo.Message, error)
	ApplicationCommandCreate(appID string, guildID string, cmd *discordgo.ApplicationCommand) (*discordgo.ApplicationCommand, error)
	InteractionRespond(*discordgo.Interaction, *discordgo.InteractionResponse) error
}

func main() {
	authToken := loadEnv()
	
	discord, err := discordgo.New("Bot " + authToken)
	discord.AddHandler(onMessageCreate)

	err = discord.Open()

	stopBot := make(chan os.Signal, 1)

	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-stopBot

	err = discord.Close()

	if err != nil {
		fmt.Println(err)
	}

	return
}

func loadEnv() string {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("読み込み出来ませんでした: %v", err)
	}

	message := os.Getenv("CLIENT_SECRET")

	return message
}


func newSession(s *discordgo.Session, bot bool) Session {
	if bot {
		var mockSession main_internal.MockSession
		return mockSession
	}
	return s
}