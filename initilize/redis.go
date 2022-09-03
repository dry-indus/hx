package initilize

import (
	"context"
	"hx/global"

	"github.com/go-redis/redis/v8"
)

func initRedis() {
	c := global.CoreRedis
	client := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Username: c.Username,
		Password: c.Password,
	})

	err := client.Ping(context.TODO()).Err()
	if err != nil {
		panic(err)
	}

	global.DL_CORE_REDIS = client
}
