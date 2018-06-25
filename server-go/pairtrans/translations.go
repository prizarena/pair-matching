package pairtrans

import (
	"github.com/strongo/bots-framework/core"
	"github.com/prizarena/prizarena-public/patrans"
)

func init() {
	patrans.RegisterTranslations(TRANS)
}

var TRANS = map[string]map[string]string{
	bots.MessageTextOopsSomethingWentWrong: {
		"ru-RU": "Упс, что-то пошло не так... \xF0\x9F\x98\xB3",
		"en-US": "Oops, something went wrong... \xF0\x9F\x98\xB3",
		"fa-IR": "اوه، یک جای کار مشکل دارد ...  \xF0\x9F\x98\xB3",
		"it-IT": "Ops, qualcosa e' andato storto... \xF0\x9F\x98\xB3",
	},
	MT_START_SELECT_LANG: {
		"en-US": "<b>Please select your language</b>\nПожалуйста выберите язык общения",
		"ru-RU": "<b>Пожалуйста выберите язык общения</b>\nPlease select your language",
	},
	ChallengeFriendCommandText: {
		"en-US": "🤺 Challenge Telegram friend",
		"ru-RU": "🤺 Новая игра в Telegram",
	},
	NewGameInlineTitle: {
		"en-US": "🀄 Pair matching - new game",
		"ru-RU": "🀄 Найди пары - новая игра",
	},
	NewGameInlineDescription: {
		"en-US": "Starts new Pair-Matching game",
		"ru-RU": "Создать новую игру",
	},
	GameCardTitle: {
		"en-US": "Pair-Matching game",
		"ru-RU": "Игра: Найди пару",
	},
	ChooseSizeOfNextBoard: {
		"en-US": "Choose size of next board:",
		"ru-RU": "Выберите размер следующей доски:",
	},
	SinglePlayer: {
		"en-US": "⚔ Single-player",
		"ru-RU": "⚔ Играть одному",
	},
	MultiPlayer: {
		"en-US": "⚔ Multi-player",
		"ru-RU": "⚔ Играть с противником",
	},
	Board: {
		"en-US": "Board",
		"ru-RU": "Доска",
	},
	Tournaments: {
		"en-US": "🏆 Tournaments",
		"ru-RU": "🏆 Турниры",
	},
	FirstMoveDoneAwaitingSecond: {
		"en-US": "Player <b>%v</b> made choice, awaiting another player...",
		"ru-RU": "Игрок <b>%v</b> сделал свой ход, ожидается ход второго игрока...",
	},
	FindFast: {
		"en-US": "Find matching pairs as fast as you can.",
		"ru-RU": "Найдите совпадающие пары настолько быстро как можете.",
	},
	RulesShort: {
		"en-US": `<pre></pre>`,
	},
	NewGameText: {
		"en-US": `🀄 <b>Pair matching game</b>

Please choose board size.`,
		"ru-RU": `🀄 Игра: <b>Найди пары</b>

Выберите размер доски.`,
	},
	MT_HOW_TO_START_NEW_GAME: {
		"en-US": `<b>To begin new game:</b>
  1. Open Telegram chat with your friend
  2. Type into the text input @BiddingTicTacToeBot and a space
  3. Wait for a popup menu and choose "New game"

<i>First 2 steps can be replaced by clicking the button below!</i>`,
		//
		"ru-RU": `<b>Чтобы начать игру:</b>
  1. Откройте чат с вашим другом
  2. Наберите в поле ввода @BiddingTicTacToeBot и пробел
  3. Дождитесь появлению меню и выберите "Новая игра"

<i>Два первых шага могут быть заменены одним кликом на кнопку ниже!</i>`,
	},
	MT_NEW_GAME_WELCOME: {
		"en-US": `To start the game please choose board size.`,
		"ru-RU": `Чтобы начать игру выберите размер доски.`,
	},
	MT_HOW_TO_INLINE: {
		"en-US": `To begin the game and to make first move:
  * choose a cell
  * click Start at bottom of the screen`,
		//
		"ru-RU": `Чтобы начать игру и сделать первый ход:
  * выберите клетку
  * нажмите Start внизу экрана`,
	},
	MT_PLAYER: {
		"en-US": "Player <b>%v</b>:",
		"ru-RU": "Игрок <b>%v</b>:",
	},
	MT_AWAITING_PLAYER: {
		"en-US": "awaiting...",
		"ru-RU": "ожидается...",
	},
	MT_PLAYER_BALANCE: {
		"en-US": "Balance: %d",
		"ru-RU": "Баланс: %d",
	},
	MT_ASK_TO_RATE: {
		"en-US": `🙋 <b>Did you like the game?</b> We'll appreciate if you <a href="https://t.me/storebot?start=BiddingTicTacToeBot">rate us</a> with 5⭐s!'`,
		"ru-RU": `🙋 <b>Понравилась игра?</b> Будем признательны если <a href="https://t.me/storebot?start=BiddingTicTacToeBot">оцените нас</a> на 5⭐!`,
	},
	// MT_FREE_GAME_SPONSORED_BY: {
	// 	"en-US": "🙏 Free game sponsored by %v",
	// 	"ru-RU": "🙏 Бесплатная игра при поддержке %v - бот для учёта долгов",
	// },
}
