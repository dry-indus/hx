package verifyser

import (
	"fmt"
	"hx/global"
	"hx/global/context"
	"hx/service/tgser"
	"hx/util"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	TgVerify TgVerifySer
)

type TgVerifySer struct{}

var ErrSenceNotExist = fmt.Errorf("sence not exist!")
var ErrTgId = fmt.Errorf("telegram id invalid!")

func (this TgVerifySer) SendCode(c context.ContextB, sence global.Sence, name string, tgId int64, length int) (code string, err error) {
	if tgId == 0 {
		return "", ErrTgId
	}

	switch sence {
	case global.RegisterSence:
		code, err = this.sendRegisterCode(c, tgId, length)
		if err != nil {
			return
		}
	default:
		return "", ErrSenceNotExist
	}

	key := fmt.Sprintf(global.VERIFY_CODE_FMT, sence, name, tgId)
	ttl := time.Duration(global.Application.VerifyCodeTTLMinutes) * time.Minute
	global.DL_CORE_REDIS.Set(c, key, code, ttl)

	return
}

func (TgVerifySer) sendRegisterCode(c context.ContextB, chatId int64, length int) (string, error) {
	code := util.RandString(length)
	text := fmt.Sprintf("【海鲜商户】尊敬的海鲜商户您好！您的注册验证码：%v, 切勿将验证码泄露给他人。", code)
	result, _ := tgser.Tg.SendText(chatId, text)
	chat := result.Chat
	if chat.ID != chatId || chat.ID == 0 {
		return "", ErrTgId
	}

	key := fmt.Sprintf(global.TG_CHAT_INFO_FMT, chat.ID)
	ttl := time.Duration(global.Application.VerifyCodeTTLMinutes) * time.Minute
	global.DL_CORE_REDIS.Set(c, key, util.MustMarshalToString(chat), ttl)

	return code, nil
}

var ErrCodeNotMatch = fmt.Errorf("code not match!")

func (this TgVerifySer) VerifyCode(c context.ContextB, sence global.Sence, name string, tgId int64, lcode string) error {
	key := fmt.Sprintf(global.VERIFY_CODE_FMT, sence, name, tgId)
	defer global.DL_CORE_REDIS.Del(c, key)
	rcode := global.DL_CORE_REDIS.Get(c, key).Val()

	if lcode != rcode {
		return ErrCodeNotMatch
	}

	return nil
}

var ErrTgNameNotMatch = fmt.Errorf("telegram name not match!")

func (TgVerifySer) VerifyTG(c context.ContextB, tgId int64, tgName string) error {
	if tgId == 0 {
		return ErrTgId
	}

	key := fmt.Sprintf(global.TG_CHAT_INFO_FMT, tgId)
	s := global.DL_CORE_REDIS.Get(c, key).Val()
	if len(s) == 0 {
		return ErrTgId
	}
	chat := tgbotapi.Chat{}
	util.JSON.UnmarshalFromString(s, &chat)
	if chat.UserName != tgName {
		return ErrTgNameNotMatch
	}

	return nil
}
