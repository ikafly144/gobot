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
package util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/dlclark/regexp2"
	"github.com/ikafly144/gobot/pkg/product"
	"github.com/ikafly144/gobot/pkg/translate"
)

func ErrorMessage(locale discordgo.Locale, err error) (res *discordgo.InteractionResponseData) {
	res = &discordgo.InteractionResponseData{}
	res.Content = ""
	res.Embeds = append(res.Embeds, &discordgo.MessageEmbed{
		Title:       translate.Message(locale, "error.message"),
		Description: err.Error(),
		Color:       0xff0000,
	})
	res.Flags = discordgo.MessageFlagsEphemeral
	return
}

func DeepcopyJson(src interface{}, dst interface{}) (err error) {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, dst)
	if err != nil {
		return err
	}
	return nil
}

func LogResp(resp *http.Response) {
}

func MessageResp(resp *http.Response) string {
	defer resp.Body.Close()
	byteArray, _ := io.ReadAll(resp.Body)
	jsonBytes := ([]byte)(byteArray)
	str := fmt.Sprintf("succeed %v %v ```json\r%v```", resp.Request.Method, resp.StatusCode, string(jsonBytes))
	return str
}

func ToEmojiA(i int) string {
	return string(rune('🇦' - 1 + i))
}

func GetCustomEmojis(s string) []*discordgo.ComponentEmoji {
	var toReturn []*discordgo.ComponentEmoji
	emojis := discordgo.EmojiRegex.FindAllString(s, -1)
	if len(emojis) < 1 {
		return toReturn
	}
	for _, em := range emojis {
		parts := strings.Split(em, ":")
		toReturn = append(toReturn, &discordgo.ComponentEmoji{
			ID:       parts[2][:len(parts[2])-1],
			Name:     parts[1],
			Animated: strings.HasPrefix(em, "<a:"),
		})
	}
	return toReturn
}

func EmojiFormat(e *discordgo.ComponentEmoji) string {
	if e.ID != "" && e.Name != "" {
		if e.Animated {
			return "<a:" + APIName(e) + ">"
		}

		return "<:" + APIName(e) + ">"
	}

	return APIName(e)
}

func APIName(e *discordgo.ComponentEmoji) string {
	if e.ID != "" && e.Name != "" {
		return e.Name + ":" + e.ID
	}
	if e.Name != "" {
		return e.Name
	}
	return e.ID
}

func Regexp2FindAllString(re *regexp2.Regexp, s string) []string {
	var matches []string
	m, _ := re.FindStringMatch(s)
	for m != nil {
		matches = append(matches, m.String())
		m, _ = re.FindNextMatch(m)
	}
	return matches
}

func StatusString(status discordgo.Status) (str string) {
	switch status {
	case discordgo.StatusOnline:
		return "<:online:1055430359363354644>"
	case discordgo.StatusDoNotDisturb:
		return "<:dnd:1055434290629980220>"
	case discordgo.StatusIdle:
		return "<:idle:1055433789020586035> "
	case discordgo.StatusInvisible:
		return "<:offline:1055434315514785792>"
	case discordgo.StatusOffline:
		return "<:offline:1055434315514785792>"
	}
	return
}

func DeferDeleteInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	time.Sleep(time.Second * 3)
	ErrorCatch("", s.InteractionResponseDelete(i.Interaction))
}

var logf *os.File

func init() {
	os.MkdirAll("logs", os.ModeDir)
	var err error
	var slash = "/"
	if runtime.GOOS == "windows" {
		slash = "\\"
	}
	logf, err = os.Create("logs" + slash + time.Now().Local().Format("2006-01-02-15-04-05") + ".log")
	stats()
	if err != nil {
		log.Printf("[ERROR] cannot open log file! : %v", err)
		logf = os.Stdout
	}
	go close()
}

func close() {
	defer logf.Close()
	defer fmt.Fprintf(logf, "\n[%v] Closed", time.Now().Local().Format(time.RFC3339))

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	<-sigCh
}

func stats() {
	fmt.Fprintf(logf,
		"[%v] ==Profiles==\n                runtime: %v %v (%v)\n                thread: max %v routine: %v\n                goroot: %v\n",
		time.Now().Local().Format(time.RFC3339),
		runtime.Version(),
		runtime.GOOS+"/"+runtime.GOARCH,
		runtime.Compiler,
		runtime.NumCPU(),
		runtime.NumGoroutine(),
		runtime.GOROOT(),
	)
}

func ErrorCatch[T any](data T, err error) (T, error) {
	if err != nil {
		fmt.Fprintf(logf, "[%v] %v\n", time.Now().Local().Format(time.RFC3339), err)
		stats()
		logf.Write(debug.Stack())
		fmt.Fprint(logf, "\n================\n\n")
	}
	return data, err
}

func WebhookExec(s *discordgo.Session, channelID string) (webhookID string, webhookToken string) {
	ws, _ := ErrorCatch(s.ChannelWebhooks(channelID))
	for _, w := range ws {
		if w.User.ID == s.State.User.ID {
			webhookID = w.ID
			webhookToken = w.Token
			return
		}
	}
	w, err := ErrorCatch(s.WebhookCreate(channelID, product.ProductName+"-webhook", s.State.User.AvatarURL("1024")))
	if err != nil {
		return
	}
	webhookID = w.ID
	webhookToken = w.Token
	return
}
