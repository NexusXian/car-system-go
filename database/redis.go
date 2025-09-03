package database

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var RDB *redis.Client
var CTX = context.Background()

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.Addr"),
		Password: viper.GetString("redis.Password"),
		DB:       viper.GetInt("redis.DB"), // 默认DB
	})

	// 测试连接
	_, err := RDB.Ping(CTX).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}
