package tgser

import (
	"fmt"
	"hx/global"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	Tg TgSer
)

type TgSer struct{}

func (TgSer) SendText(chatId int64, text string) (result tgbotapi.Message, err error) {
	msg := tgbotapi.NewMessage(chatId, text)
	return global.DL_HX_BOT.Send(msg)
}

func (TgSer) Reply(update tgbotapi.Update) tgbotapi.Chattable {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("邀请码：%v", update.Message.Chat.ID))
	msg.ReplyToMessageID = update.Message.MessageID
	return msg
}
