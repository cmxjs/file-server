package models

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/cmxjs/file-server/dao"
)

// set key
func GetLog(key string) (string, error) {
	if dao.RedisDB == nil {
		return "", errors.New("redis connect not exists")
	}

	var ctx = context.Background()

	value, err := dao.RedisDB.Get(ctx, key).Result()
	if err == nil {
		log.Printf("Success to get key from redis, key: %v\n", key)
	}
	return value, err
}

// get key
func SetLog(key string, value string, expired time.Duration) error {
	if dao.RedisDB == nil {
		return errors.New("redis connect not exists")
	}

	var ctx = context.Background()

	err := dao.RedisDB.Set(ctx, key, value, expired).Err()
	if err == nil {
		log.Printf("Success to set key in redis, key: %v\n", key)
	} else {
		log.Printf("Failed to set key in redis. key: %v. err: %v/n", key, err)
	}
	return err
}
