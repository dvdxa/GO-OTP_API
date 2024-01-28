package redis

import (
	"otp_api/conf"
	"otp_api/internal/app/pkg/logger"
	"time"

	"github.com/go-redis/redis"
)

func NewRedisClient(cfg *conf.Config, log *logger.Logger) *redis.Client {
	redisHost := cfg.Redis.RedisAddr

	if redisHost == "" {
		redisHost = ":6379"
	}

	// Wrap the client's methods with logging

	client := redis.NewClient(&redis.Options{
		Addr:         redisHost,
		MinIdleConns: cfg.Redis.MinIdleConns, //pool conns
		PoolSize:     cfg.Redis.PoolSize,
		PoolTimeout:  time.Duration(cfg.Redis.PoolTimeout) * time.Second,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
	})

	client.WrapProcess(func(oldProcess func(cmd redis.Cmder) error) func(cmd redis.Cmder) error {
		return func(cmd redis.Cmder) error {
			start := time.Now()
			err := oldProcess(cmd)
			elapsed := time.Since(start)

			log.Printf("COMMAND: %s, ELAPSED: %s, ERROR: %v", cmd.Name(), elapsed, err)

			return err
		}
	})

	return client
}
