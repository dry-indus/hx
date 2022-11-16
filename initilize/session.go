package initilize

import (
	"context"
	"fmt"
	"hx/global"
	"time"

	"github.com/gorilla/sessions"
	"github.com/rbcervilla/redisstore/v8"
)

func initSession() {
	var err error
	global.DL_U_SESSION_STORE, err = redisstore.NewRedisStore(context.Background(), global.DL_CORE_REDIS)
	if err != nil {
		panic(fmt.Sprint("failed to create DL_U_SESSION_STORE: ", err))
	}

	global.DL_U_SESSION_STORE.Options(sessions.Options{
		MaxAge:   int(24 * time.Hour / time.Second),
		HttpOnly: true,
	})

	global.DL_M_SESSION_STORE, err = redisstore.NewRedisStore(context.Background(), global.DL_CORE_REDIS)
	if err != nil {
		panic(fmt.Sprint("failed to create DL_M_SESSION_STORE: ", err))
	}

	global.DL_M_SESSION_STORE.Options(sessions.Options{
		MaxAge:   int(8 * time.Hour / time.Second),
		HttpOnly: true,
	})
}
