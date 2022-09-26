package initilize

import (
	"hx/global"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func initTgBot(handle TGHandle) {
	hxBot, err := tgbotapi.NewBotAPI(global.Telegram.HXBotToken)
	if err != nil {
		global.DL_LOGGER.Warnf("new hx tg bot failed! err: %v", err)
		return
	}
	hxBot.Debug = global.Telegram.HXBotDebug

	global.DL_HX_BOT = hxBot

	go TgListen(hxBot, handle)
}

func TgListen(bot *tgbotapi.BotAPI, handle TGHandle) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			global.DL_LOGGER.Debugf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if handle == nil {
				continue
			}

			bot.Send(handle.Reply(update))
		}
	}
}

type TGHandle interface {
	Reply(update tgbotapi.Update) tgbotapi.Chattable
}
