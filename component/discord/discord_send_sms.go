package discord

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/teoit/gosctx"
)

var (
	ErrDiscordTokenEmpty   = errors.New("error Token Discord Empty")
	ErrDiscordToken        = errors.New("error Token Discord Auth")
	ErrConnectDiscord      = errors.New("error Connect Discord")
	ErrChannelDiscordEmpty = errors.New("error Channel Discord Empty")
)

type SendMessageDiscordSVC interface {
	SendMessageDev(msg string) error
	SendMessage(channelId string, msg string) error
}

type discordClient struct {
	id            string
	failOnError   bool
	token         string
	devChannelIds string
	channelIdDev  []string
	logger        gosctx.Logger
	client        *discordgo.Session
}

func NewDiscordClient(id string) *discordClient {
	return &discordClient{id: id}
}

func (d *discordClient) InitFlags() {
	flag.StringVar(&d.token, "discord-token", "", "token discord")
	flag.StringVar(&d.devChannelIds, "discord-dev-channel-ids", "", "channel ids discord, comma separated")
	flag.BoolVar(&d.failOnError, "discord-fail-on-error", false, "fail on error")
}

func (d *discordClient) ID() string {
	return d.id
}

func (d *discordClient) Activate(sc gosctx.ServiceContext) error {
	d.logger = sc.Logger(d.id)

	if d.token == "" {
		err := ErrDiscordTokenEmpty
		if d.failOnError {
			return err
		}
		d.logger.Error(err.Error())
		return nil
	}

	d.token = fmt.Sprintf("Bot %s", d.token)

	var channelIds []string
	if d.devChannelIds != "" {
		list := strings.Split(d.devChannelIds, ",")
		for _, v := range list {
			v = strings.TrimSpace(v)
			if v == "" {
				continue
			}
			channelIds = append(channelIds, v)
		}
		d.channelIdDev = channelIds
	}

	if len(d.channelIdDev) == 0 {
		err := ErrChannelDiscordEmpty
		if d.failOnError {
			return err
		}
		d.logger.Error(err.Error())
		return nil
	}

	dg, err := discordgo.New(d.token)
	if err != nil {
		err = ErrDiscordToken
		if d.failOnError {
			return err
		}
		d.logger.Error(err.Error())
		return nil
	}

	err = dg.Open()
	if err != nil {
		err = ErrConnectDiscord
		if d.failOnError {
			return err
		}
		d.logger.Error(err.Error())
		return nil
	}

	d.logger.Infof("Discord bot is now running time %s", time.Now().Format(time.RFC3339))
	d.client = dg

	return nil
}

func (d *discordClient) Stop() error {
	if d.client != nil {
		return d.client.Close()
	}
	return nil
}

func (d *discordClient) SendMessageDev(msg string) error {
	if len(d.channelIdDev) == 0 {
		return ErrChannelDiscordEmpty
	}

	for _, channelId := range d.channelIdDev {
		if channelId == "" {
			continue
		}
		if _, err := d.client.ChannelMessageSend(channelId, msg); err != nil {
			if d.failOnError {
				return err
			}
			d.logger.Errorf("error send message to discord channel %s: %v", channelId, err)

			return nil
		}
	}

	return nil
}

func (d *discordClient) SendMessage(channelId string, msg string) error {
	if channelId == "" {
		return ErrChannelDiscordEmpty
	}

	if _, err := d.client.ChannelMessageSend(channelId, msg); err != nil {
		if d.failOnError {
			return err
		}
		d.logger.Errorf("error send message to discord channel %s: %v", channelId, err)
		return nil
	}

	return nil
}
