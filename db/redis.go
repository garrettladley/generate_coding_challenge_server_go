package db

import (
	"github.com/garrettladley/generate_coding_challenge_server_go/config"
	"github.com/redis/go-redis/v9"
)

func CreateRedisConnection(env config.EnvVars) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     env.REDIS_ADDR,
		Password: env.REDIS_PASSWORD,
		DB:       env.REDIS_DB,
	})

	return rdb
}
