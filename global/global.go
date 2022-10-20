package global

import (
	"fmt"

	osssdk "github.com/aliyun/aliyun-oss-go-sdk/oss"
	redis8 "github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/qiniu/qmgo"
	"github.com/rbcervilla/redisstore/v8"
	"github.com/sirupsen/logrus"
)

var (
	DL_CLOSE           []func() error
	DL_CORE_CLI        *qmgo.Client
	DL_CORE_MDB        *qmgo.Database
	DL_LOGGER          *logrus.Logger
	DL_CORE_REDIS      *redis8.Client
	DL_U_SESSION_STORE *redisstore.RedisStore
	DL_M_SESSION_STORE *redisstore.RedisStore
	DL_HX_BOT          *tgbotapi.BotAPI
	DL_OSS_BUCKET      *osssdk.Bucket
)

func Close() {
	for _, f := range DL_CLOSE {
		err := f()
		if err != nil {
			if DL_LOGGER != nil {
				DL_LOGGER.Error(err)
			} else {
				fmt.Println(err)
			}
		}
	}
}
