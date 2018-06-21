package paircommands

import (
	"github.com/strongo/bots-framework/core"
	"github.com/strongo/bots-framework/platforms/telegram"
	"strings"
	"github.com/prizarena/prizarena-public/pabot"
	"github.com/prizarena/prizarena-public/pamodels"
	"github.com/strongo/bots-api-telegram"
	"github.com/strongo/app"
	"github.com/prizarena/turn-based"
	"github.com/prizarena/rock-paper-scissors/server-go/rpstrans"
	"github.com/prizarena/rock-paper-scissors/server-go/rpssecrets"
	"github.com/prizarena/pair-matching/server-go/pairmodels"
	"time"
	"fmt"
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
			//return inlineQueryTournament(whc, inlineQuery)
		case inlineQuery.Text == "" || inlineQuery.Text == "play" || strings.HasPrefix(inlineQuery.Text, "play?tournament="):
			return inlineQueryPlay(whc, inlineQuery)
		}
		return
	},
)

// func inlineQueryDefault(whc bots.WebhookContext, inlineQuery pabot.InlineQueryContext) (m bots.MessageFromBot, err error) {
// 	return
// }



func inlineQueryPlay(whc bots.WebhookContext, inlineQuery pabot.InlineQueryContext) (m bots.MessageFromBot, err error) {
	return pabot.ProcessInlineQueryTournament(whc, inlineQuery, rpssecrets.RpsPrizarenaGameID, "tournament",
		func(tournament pamodels.Tournament) (m bots.MessageFromBot, err error) {
			c := whc.Context()

			translator := whc.BotAppContext().GetTranslator(c)

			newGameOption := func(lang string) tgbotapi.InlineQueryResultArticle {
				t := strongo.NewSingleMapTranslator(strongo.LocalesByCode5[lang], translator)
				newBoard := pairmodels.PairsBoard{
					PairsBoardEntity: &pairmodels.PairsBoardEntity{
						SizeY: 8,
						SizeX: 8,
						BoardEntityBase: turnbased.BoardEntityBase{Lang: lang, Round: 1},
						Cells: pairmodels.Shuffle(8, 8),
					},
				}
				newBoard.Created = time.Now()

				// Renders game board to a Telegram message to return as inline result
				if m, err = renderPairsBoardMessage(t, &tournament, newBoard); err != nil {
					panic(err)
				}

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
					ReplyMarkup: m.Keyboard.(*tgbotapi.InlineKeyboardMarkup),
				}
			}

			m.BotMessage = telegram.InlineBotMessage(tgbotapi.InlineConfig{
				InlineQueryID: inlineQuery.ID,
				Results: []interface{}{
					newGameOption("en-US"),
					// newGameOption("ru-RU"),
				},
			})
			return
		})
	return
}

func renderPairsBoardMessage(t strongo.SingleLocaleTranslator, tournament *pamodels.Tournament, board pairmodels.PairsBoard) (m bots.MessageFromBot, err error) {
	kbRows := make([][]tgbotapi.InlineKeyboardButton, board.SizeY)
	for y, row := range board.Rows() {
		kbRow := make([]tgbotapi.InlineKeyboardButton, board.SizeX)
		for x, cell := range row {
			kbRow[y] = tgbotapi.InlineKeyboardButton{Text: string(cell), CallbackData: fmt.Sprintf("open?x=%v&y%v", x+1, y+1)}
		}
		kbRows[y] = kbRow
	}
	m.Text = "<b>New game</b>\nFind matching pairs as quickly as you can."
	m.Keyboard = tgbotapi.NewInlineKeyboardMarkup(kbRows...)
	return
}
