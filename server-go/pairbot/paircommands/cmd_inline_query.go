package paircommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"strings"
	"github.com/prizarena/prizarena-public/pabot"
	"github.com/prizarena/prizarena-public/pamodels"
	"github.com/strongo/bots-api-telegram"
	"bytes"
	"fmt"
	"strconv"
	"github.com/prizarena/pair-matching/server-go/pairtrans"
	"github.com/prizarena/pair-matching/server-go/pairsecrets"
)

var inlineQueryCommand = bots.NewInlineQueryCommand(
	"inline-query",
	func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		tgInlineQuery := whc.Input().(telegram.TgWebhookInlineQuery)
		inlineQuery := pabot.InlineQueryContext{
			ID:   tgInlineQuery.GetInlineQueryID(),
			Text: strings.TrimSpace(tgInlineQuery.TgUpdate().InlineQuery.Query),
		}
		words := strings.Split(inlineQuery.Text, " ")

		removeLang := func() {
			if len(words) == 1 {
				words = []string{}
			} else {
				words = words[1:]
			}
		}
		switch words[0] {
		case "ru":
			whc.SetLocale("ru-RU")
			removeLang()
		case "en":
			words = words[1:]
			removeLang()
		}

		inlineQuery.Text = strings.Join(words, " ")

		switch {
		case strings.HasPrefix(inlineQuery.Text, "tournament?id="):
			// return inlineQueryTournament(whc, inlineQuery)
		case inlineQuery.Text == "" || inlineQuery.Text == "play" || strings.HasPrefix(inlineQuery.Text, "play?tournament="):
			return inlineQueryPlay(whc, inlineQuery)
		}
		return
	},
)

// func inlineQueryDefault(whc bots.WebhookContext, inlineQuery pabot.InlineQueryContext) (m bots.MessageFromBot, err error) {
// 	return
// }

func inlineKbMarkup(lang string) *tgbotapi.InlineKeyboardMarkup {
	sizeButton := func(width, height int) tgbotapi.InlineKeyboardButton {
		return tgbotapi.InlineKeyboardButton{
			Text:         fmt.Sprintf(strconv.Itoa(width) + "x" + strconv.Itoa(height)),
			CallbackData: newBoardCallbackData(width, height, lang),
		}
	}
	return tgbotapi.NewInlineKeyboardMarkup(
		[]tgbotapi.InlineKeyboardButton{
			sizeButton(4, 2),
			sizeButton(4, 3),
			sizeButton(4, 4),
			sizeButton(5, 4),
		},
		[]tgbotapi.InlineKeyboardButton{
			sizeButton(6, 4),
			sizeButton(6, 5),
			sizeButton(6, 6),
		},
		[]tgbotapi.InlineKeyboardButton{
			sizeButton(7, 6),
			sizeButton(8, 6),
			sizeButton(8, 7),
			sizeButton(8, 8),
		},
		[]tgbotapi.InlineKeyboardButton{
			sizeButton(8, 9),
			sizeButton(8, 10),
			sizeButton(8, 11),
			sizeButton(8, 12),
		},
	)
}

var newBoardSizesKeyboard = map[string]*tgbotapi.InlineKeyboardMarkup{
	"en-US": inlineKbMarkup("en-US"),
	"ru-RU": inlineKbMarkup("ru-RU"),
}

func inlineQueryPlay(whc bots.WebhookContext, inlineQuery pabot.InlineQueryContext) (m bots.MessageFromBot, err error) {
	return pabot.ProcessInlineQueryTournament(whc, inlineQuery, pairsecrets.PrizarenaGameID, "tournament",
		func(tournament pamodels.Tournament) (m bots.MessageFromBot, err error) {
			// c := whc.Context()

			// translator := whc.BotAppContext().GetTranslator(c)

			newGameOption := func() tgbotapi.InlineQueryResultArticle {
				// t := strongo.NewSingleMapTranslator(strongo.LocalesByCode5[lang], translator)

				lang := whc.Locale().Code5

				articleID := "new_game?l=" + lang
				if tournament.ID != "" {
					articleID += "&t=" + tournament.ShortTournamentID()
				}
				text := new(bytes.Buffer)
				text.WriteString(whc.Translate(pairtrans.NewGameText))
				return tgbotapi.InlineQueryResultArticle{
					ID:          articleID,
					Type:        "article",
					Title:       whc.Translate(pairtrans.NewGameInlineTitle),
					Description: whc.Translate(pairtrans.NewGameInlineDescription),
					InputMessageContent: tgbotapi.InputTextMessageContent{
						Text:                  text.String(),
						ParseMode:             "HTML",
						DisableWebPagePreview: m.DisableWebPagePreview,
					},
					ReplyMarkup: newBoardSizesKeyboard[lang],
				}
			}

			m.BotMessage = telegram.InlineBotMessage(tgbotapi.InlineConfig{
				InlineQueryID: inlineQuery.ID,
				Results: []interface{}{
					newGameOption(),
					// newGameOption("ru-RU"),
				},
				CacheTime: 10,
			})
			return
		})
	return
}
