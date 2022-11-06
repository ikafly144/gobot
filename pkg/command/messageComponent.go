package command

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Tnze/go-mc/chat"
	"github.com/bwmarrin/discordgo"
	"github.com/ikafly144/gobot/pkg/api"
	"github.com/millkhan/mcstatusgo/v2"
)

func MCpanelRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	component := i.Message.Components
	var content string
	bytes, _ := component[0].MarshalJSON()
	gid := i.GuildID
	uid := i.Member.User.ID
	if component[0].Type() == discordgo.ActionsRowComponent {
		data := &discordgo.ActionsRow{}
		json.Unmarshal(bytes, data)
		bytes, _ := data.Components[0].MarshalJSON()
		if data.Components[0].Type() == discordgo.SelectMenuComponent {
			data := &discordgo.SelectMenu{}
			json.Unmarshal(bytes, data)
			for _, v := range data.Options {
				for _, m := range i.Member.Roles {
					if v.Value == m {
						t := true
						for _, v2 := range i.MessageComponentData().Values {
							if v2 != v.Value {
								t = false
							}
						}
						if t {
							s.GuildMemberRoleRemove(gid, uid, v.Value)
							content += "はく奪 <@&" + v.Value + ">\r"
						}
					}
				}
			}
			for _, r := range i.MessageComponentData().Values {
				t := true
				for _, m := range i.Member.Roles {
					if r == m {
						t = false
					}
				}
				if t {
					s.GuildMemberRoleAdd(gid, uid, r)
					content += "付与 <@&" + r + ">\r"
				}
			}
		}
	}
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func MCpanelRoleAdd(s *discordgo.Session, i *discordgo.InteractionCreate) {
	mid := i.Message.Embeds[0].Title
	gid := i.GuildID
	cid := i.ChannelID
	mes, _ := s.ChannelMessage(cid, mid)
	rv := i.MessageComponentData().Values
	roles := []discordgo.Role{}
	for _, v := range rv {
		role, _ := s.State.Role(gid, v)
		roles = append(roles, *role)
	}
	options := []discordgo.SelectMenuOption{}
	for _, r := range roles {
		options = append(options, discordgo.SelectMenuOption{
			Label: r.Name,
			Value: r.ID,
		})
	}
	var fields string
	for _, r := range roles {
		fields += r.Mention() + "\r"
	}
	zero := 0
	content := discordgo.MessageEdit{
		ID:      mid,
		Channel: cid,
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: mes.Embeds[0].Title,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:  "roles",
						Value: fields,
					},
				},
			},
		},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:  "gobot_panel_role",
						MinValues: &zero,
						MaxValues: len(options),
						Options:   options,
					},
				},
			},
		},
	}
	s.ChannelMessageEditComplex(&content)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "OK",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func MCpanelMinecraft(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if fmt.Sprint(i.MessageComponentData().Values) == "[]" {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "OK",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		s.InteractionResponseDelete(i.Interaction)
		return
	}
	initialTimeOut := time.Second * 2
	ioTimeOut := time.Second * 2
	data := i.MessageComponentData()
	addresses := strings.Split(data.Values[0], ":")
	name := addresses[0]
	address := addresses[1]
	port, err := strconv.Atoi(addresses[2])
	if err != nil {
		log.Print(err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "ポート値が不正です",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	showIp, err := strconv.ParseBool(addresses[3])
	if err != nil {
		log.Print(err)
	}
	q, err := mcstatusgo.Status(address, uint16(port), initialTimeOut, ioTimeOut)
	if err != nil {
		log.Print(err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "サーバーの状況を取得できませんでした\r" + fmt.Sprint(err),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	message := chat.Message{}
	message.UnmarshalJSON([]byte(q.Description))
	hash := sha256.New()
	thumb := strings.ReplaceAll(q.Favicon, "data:image/png;base64,", "")
	io.WriteString(hash, thumb)
	str := hash.Sum(nil)
	code := hex.EncodeToString(str)
	bd := &api.ImagePngHash{
		Data: thumb,
		Hash: code,
	}
	b, _ := json.Marshal(bd)
	api.GetApi("/api/image/png/add", bytes.NewBuffer(b))
	log.Print(q.Version.Protocol)
	log.Print("https://sabafly.net/api/decode?s=" + code)
	embeds := []*discordgo.MessageEmbed{
		{
			Title:       name,
			Description: message.ClearString(),
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL:    "https://sabafly.net/api/decode?s=" + code,
				Width:  64,
				Height: 64,
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "プレイヤー",
					Value:  strconv.Itoa(q.Players.Online) + "/" + strconv.Itoa(q.Players.Max),
					Inline: true,
				},
				{
					Name:   "レイテンシ",
					Value:  strconv.Itoa(int(q.Latency.Abs().Milliseconds())) + "ms",
					Inline: true,
				},
				{
					Name:   "バージョン",
					Value:  q.Version.Name,
					Inline: true,
				},
			},
		},
	}
	if showIp {
		embeds[0].Fields = append(embeds[0].Fields, &discordgo.MessageEmbedField{
			Name:   "アドレス",
			Value:  address + ":" + strconv.Itoa(int(q.Port)),
			Inline: true,
		})
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Print(err)
	}
}
