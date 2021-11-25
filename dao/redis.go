package dao

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var (
	RedisDB *redis.Client
)

func InitRedis(addr string, dbNum string) (err error) {
	dbNumInt, err := strconv.Atoi(dbNum)
	if err != nil {
		return err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       dbNumInt,
	})
	var ctx = context.Background()
	_, err = rdb.Ping(ctx).Result()
	if err == nil {
		RedisDB = rdb
	}
	return err
}
