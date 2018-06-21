package paircommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"strings"
	"github.com/prizarena/prizarena-public/pabot"
	"github.com/prizarena/prizarena-public/pamodels"
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/app"
		"github.com/prizarena/rock-paper-scissors/server-go/rpstrans"
	"github.com/prizarena/rock-paper-scissors/server-go/rpssecrets"
			)

var inlineQueryCommand = bots.NewInlineQueryCommand(
	"inline-query",
	func(whc bots.WebhookContext) (m bots.MessageFromBot, err error) {
		tgInlineQuery := whc.Input().(telegram.TgWebhookInlineQuery)
		inlineQuery := pabot.InlineQueryContext{
			ID:   tgInlineQuery.GetInlineQueryID(),
			Text: strings.TrimSpace(tgInlineQuery.TgUpdate().InlineQuery.Query),
		}

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

var newBoardSizesKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	[]tgbotapi.InlineKeyboardButton{
		{Text: "4x3", CallbackData: newBoardCallbackData(4,3)},
		{Text: "4x4", CallbackData: newBoardCallbackData(4,4)},
		{Text: "5x4", CallbackData: newBoardCallbackData(5,4)},
	},
	[]tgbotapi.InlineKeyboardButton{
		{Text: "6x4", CallbackData: newBoardCallbackData(6,4)},
		{Text: "6x5", CallbackData: newBoardCallbackData(6,5)},
		{Text: "6x6", CallbackData: newBoardCallbackData(6,6)},
		{Text: "7x6", CallbackData: newBoardCallbackData(7,6)},
	},
	[]tgbotapi.InlineKeyboardButton{
		{Text: "8x6", CallbackData: newBoardCallbackData(8,6)},
		{Text: "8x7", CallbackData: newBoardCallbackData(8,7)},
		{Text: "8x8", CallbackData: newBoardCallbackData(8,8)},
		{Text: "8x9", CallbackData: newBoardCallbackData(8,9)},
	},
	[]tgbotapi.InlineKeyboardButton{
		{Text: "8x10", CallbackData: newBoardCallbackData(8,6)},
		{Text: "8x11", CallbackData: newBoardCallbackData(8,7)},
		{Text: "8x12", CallbackData: newBoardCallbackData(8,8)},
	},
)

func inlineQueryPlay(whc bots.WebhookContext, inlineQuery pabot.InlineQueryContext) (m bots.MessageFromBot, err error) {
	return pabot.ProcessInlineQueryTournament(whc, inlineQuery, rpssecrets.RpsPrizarenaGameID, "tournament",
		func(tournament pamodels.Tournament) (m bots.MessageFromBot, err error) {
			c := whc.Context()

			translator := whc.BotAppContext().GetTranslator(c)

			newGameOption := func(lang string) tgbotapi.InlineQueryResultArticle {
				t := strongo.NewSingleMapTranslator(strongo.LocalesByCode5[lang], translator)

				articleID := "new_game?l=" + lang
				if tournament.ID != "" {
					articleID += "&t=" + tournament.ShortTournamentID()
				}
				return tgbotapi.InlineQueryResultArticle{
					ID:          articleID,
					Type:        "article",
					Title:       t.Translate(rpstrans.NewGameInlineTitle),
					Description: t.Translate(rpstrans.NewGameInlineDescription),
					InputMessageContent: tgbotapi.InputTextMessageContent{
						Text:                  m.Text,
						ParseMode:             "HTML",
						DisableWebPagePreview: m.DisableWebPagePreview,
					},
					ReplyMarkup: newBoardSizesKeyboard,
				}
			}

			m.BotMessage = telegram.InlineBotMessage(tgbotapi.InlineConfig{
				InlineQueryID: inlineQuery.ID,
				Results: []interface{}{
					newGameOption("en-US"),
					// newGameOption("ru-RU"),
				},
				CacheTime: 10,
			})
			return
		})
	return
}

