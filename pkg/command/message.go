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
package command

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/ikafly144/gobot/pkg/product"
	"github.com/ikafly144/gobot/pkg/session"
	"github.com/ikafly144/gobot/pkg/translate"
	"github.com/ikafly144/gobot/pkg/types"
	"github.com/ikafly144/gobot/pkg/util"
)

func MModify(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := &discordgo.ApplicationCommandInteractionData{}
	byte, _ := json.Marshal(i.Interaction.Data)
	json.Unmarshal(byte, data)
	mes, err := s.ChannelMessage(i.ChannelID, data.TargetID)
	if err != nil {
		log.Print(err)
	}
	if mes.Author.ID == s.State.User.ID {
		if len(mes.Components) != 0 {
			for _, v := range mes.Components {
				if v.Type() == discordgo.ActionsRowComponent {
					byte, _ := v.MarshalJSON()
					data := &discordgo.ActionsRow{}
					json.Unmarshal(byte, data)
					for _, v := range data.Components {
						if v.Type() == discordgo.SelectMenuComponent {
							byte, _ := v.MarshalJSON()
							data := &discordgo.SelectMenu{}
							json.Unmarshal(byte, data)
							switch data.CustomID {
							case product.CommandPanelRole:
								panelRole(s, i, mes)
							case product.CommandPanelMinecraft:
								panelMinecraft(s, i, mes)
							}
						}
					}
				}
			}
		}
	} else {
		str := translate.Message(i.Locale, "message_modify_cant_use_this_message_error")
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Content: &str,
		})
	}
}

func panelRole(s *discordgo.Session, i *discordgo.InteractionCreate, mes *discordgo.Message) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	one := 1
	str := translate.Message(i.Locale, "message_modify_role_add_message")
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &str,
		Embeds: &[]*discordgo.MessageEmbed{
			{
				Title: mes.ID,
			},
		},
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						MenuType:  discordgo.RoleSelectMenu,
						CustomID:  product.CommandPanelAdd + ":" + session.InteractionSave(i),
						MinValues: &one,
						MaxValues: 25,
					},
				},
			},
		},
	})
}

func panelMinecraft(s *discordgo.Session, i *discordgo.InteractionCreate, mes *discordgo.Message) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    product.CommandPanelMinecraftAddServerName,
							Label:       translate.Message(i.Locale, "panel_minecraft_display_name"),
							Placeholder: translate.Message(i.Locale, "panel_minecraft_my_server"),
							Style:       discordgo.TextInputShort,
							MaxLength:   16,
							Required:    true,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    product.CommandPanelMinecraftAddAddress,
							Label:       translate.Message(i.Locale, "panel_minecraft_address"),
							Placeholder: "example.com",
							Required:    true,
							MaxLength:   32,
							Style:       discordgo.TextInputShort,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    product.CommandPanelMinecraftAddPort,
							Label:       translate.Message(i.Locale, "panel_minecraft_port"),
							Placeholder: "25565",
							Style:       discordgo.TextInputShort,
							Value:       "25565",
							MinLength:   5,
							MaxLength:   5,
							Required:    true,
						},
					},
				},
			},
			CustomID: product.CommandPanelMinecraftAddModal + ":" + mes.ID,
			Title:    translate.Message(i.Locale, "panel_minecraft_add_server"),
			Flags:    discordgo.MessageFlagsEphemeral,
		},
	})
	log.Print(err)
}

var selects map[types.MessageSelect]*discordgo.Message = make(map[types.MessageSelect]*discordgo.Message)

func GetSelectingMessage(uid string, gid string) (mes *discordgo.Message, err error) {
	id := types.MessageSelect{
		MemberID: uid,
		GuildID:  gid,
	}
	if m, ok := selects[id]; ok {
		mes = m
		return
	} else {
		err = errors.New("no message is selected")
		return
	}
}

func MSelect(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
	data := &discordgo.ApplicationCommandInteractionData{}
	byte, _ := json.Marshal(i.Interaction.Data)
	json.Unmarshal(byte, data)
	mes, err := s.ChannelMessage(i.ChannelID, data.TargetID)
	if err != nil {
		log.Print(err)
	}
	id := types.MessageSelect{
		MemberID: i.Member.User.ID,
		GuildID:  i.GuildID,
	}
	selects[id] = mes
	str := translate.Message(i.Locale, "message_command_select_message")
	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &str,
	})
	util.DeferDeleteInteraction(s, i)
}

func RemoveSelect(uid string, gid string) {
	id := types.MessageSelect{
		MemberID: uid,
		GuildID:  gid,
	}
	delete(selects, id)
}
