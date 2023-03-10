package cache

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	_redis *redis.Client
)

func init() {
	db := 0
	if _db, err := strconv.Atoi(os.Getenv("REDIS_DATABASE")); err == nil {
		db = _db
	}
	_redis = connect(db)
}

func connect(db int) *redis.Client {
	var option *redis.Options
	option = &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	}
	if redisUrl := os.Getenv("REDIS_URL"); redisUrl != "" {
		conf := strings.Split(strings.Split(redisUrl, "://:")[1], "@")
		option = &redis.Options{
			Addr:     conf[1],
			Password: conf[0],
			DB:       db,
		}
	}
	return redis.NewClient(option)
}

func Set(key string, val interface{}, expires time.Duration) error {
	ctx := context.Background()
	defer ctx.Done()
	_, err := _redis.Set(ctx, key, val, expires).Result()
	return err
}

func Get(key string) (string, bool) {
	ctx := context.Background()
	defer ctx.Done()
	s, err := _redis.Get(ctx, key).Result()
	return s, err == nil
}

func Pull(key string) (string, bool) {
	ctx := context.Background()
	defer ctx.Done()
	s, err := _redis.GetDel(ctx, key).Result()
	if err != redis.Nil {
		if strings.Contains(err.Error(), "getdel") {
			v, ok := Get(key)
			Forget(key)
			return v, ok
		}
		return "", err == nil
	}
	return s, err == nil
}

func Exist(key string) bool {
	ctx := context.Background()
	defer ctx.Done()
	_, err := _redis.Get(ctx, key).Result()
	return err != redis.Nil
}

func Forget(key string) error {
	ctx := context.Background()
	defer ctx.Done()
	_, err := _redis.Del(ctx, key).Result()
	return err
}

func Forever() {}

func Duration(unix int64) time.Duration {
	return time.Until(time.Unix(unix, 0))
}
