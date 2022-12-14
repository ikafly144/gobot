/*
	Copyright (C) 2022  ikafly144

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package interaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Tnze/go-mc/chat"
	"github.com/bwmarrin/discordgo"
	"github.com/ikafly144/gobot/pkg/product"
	"github.com/ikafly144/gobot/pkg/translate"
	"github.com/ikafly144/gobot/pkg/util"
	"github.com/millkhan/mcstatusgo/v2"
)

func ComponentPanelMinecraft(s *discordgo.Session, i *discordgo.InteractionCreate) {
	util.ErrorCatch("", s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}))
	if fmt.Sprint(i.MessageComponentData().Values) == "[]" {
		s.InteractionResponseDelete(i.Interaction)
		return
	}
	initialTimeOut := time.Second * 10
	ioTimeOut := time.Second * 30
	data := i.MessageComponentData()
	addresses := strings.Split(data.Values[0], ":")
	name := addresses[0]
	address := addresses[1]
	port, err := util.ErrorCatch(strconv.ParseUint(addresses[2], 10, 16))
	if err != nil {
		e := translate.ErrorEmbed(i.Locale, "error_invalid_port_value")
		util.ErrorCatch(s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &e,
		}))
		return
	}
	showIp, _ := util.ErrorCatch(strconv.ParseBool(addresses[3]))
	q, err := util.ErrorCatch(mcstatusgo.Status(address, uint16(port), initialTimeOut, ioTimeOut))
	if err != nil {
		e := translate.ErrorEmbed(i.Locale, "error_failed_to_ping_server")
		util.ErrorCatch(s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &e,
		}))
		return
	}
	message := chat.Message{}
	util.ErrorCatch("", message.UnmarshalJSON([]byte(q.Description)))
	hash := sha256.New()
	thumb := strings.ReplaceAll(q.Favicon, "data:image/png;base64,", "")
	res, _ := util.ErrorCatch(base64.RawStdEncoding.DecodeString(thumb))
	util.ErrorCatch(io.WriteString(hash, thumb))
	str := hash.Sum(nil)
	code := hex.EncodeToString(str)
	color := 0x00ff00
	if q.Version.Protocol == 46 {
		color = 0xff0000
	}
	var player string
	for _, v := range q.Players.Sample {
		player += v["name"] + "\r"
	}
	if player != "" {
		player = "```" + player + "```"
	}
	embeds := []*discordgo.MessageEmbed{
		{
			Title:       name,
			Description: "```ansi\r" + message.String() + "```",
			Color:       color,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: "attachment://" + code + ".png",
			},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			Footer:    &discordgo.MessageEmbedFooter{Text: product.ProductName},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   translate.Message(i.Locale, "players"),
					Value:  "```" + strconv.Itoa(q.Players.Online) + "/" + strconv.Itoa(q.Players.Max) + "```" + player,
					Inline: true,
				},
				{
					Name:   translate.Message(i.Locale, "latency"),
					Value:  "```" + strconv.Itoa(int(q.Latency.Abs().Milliseconds())) + "ms" + "```",
					Inline: true,
				},
				{
					Name:   translate.Message(i.Locale, "version"),
					Value:  "```ansi\r" + chat.Text(q.Version.Name).String() + "```",
					Inline: true,
				},
			},
		},
	}
	if showIp {
		embeds[0].Fields = append(embeds[0].Fields, &discordgo.MessageEmbedField{
			Name:   translate.Message(i.Locale, "address"),
			Value:  "```" + address + "```",
			Inline: true,
		},
			&discordgo.MessageEmbedField{
				Name:   translate.Message(i.Locale, "port"),
				Value:  "```" + strconv.Itoa(int(q.Port)) + "```",
				Inline: true,
			})
	}
	util.ErrorCatch(s.ChannelMessageEditComplex(discordgo.NewMessageEdit(i.ChannelID, i.Message.ID)))
	util.ErrorCatch(s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &embeds,
		Files: []*discordgo.File{
			{
				Name:        code + ".png",
				ContentType: "image/png",
				Reader:      bytes.NewReader(res),
			},
		},
	}))
}
