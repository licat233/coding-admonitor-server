package redisx

import (
	"admonitor-server/common/global"
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//NewRedis 创建redis客户端
func NewRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.pass"),
		DB:       viper.GetInt("redis.db"),
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.LOG.Error("redis connect ping failed, err:", zap.Error(err))
		panic("连接redis服务器失败")
	}
	return client
}
