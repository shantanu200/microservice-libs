package database

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shantanu200/microservice-libs/types"
)

type Redis struct {
	RDB *redis.Client
	ctx context.Context
}

func NewRedisPool() *Redis {
	return &Redis{}
}

func (r *Redis) Connect(options *types.RedisConnectionConfig) error {
	if options.Addr == "" {
		options.Addr = "redis://localhost:6379"
	}

	if options.PoolSize == 0 {
		options.PoolSize = 10
	}

	if options.MaxRetries == 0 {
		options.MaxRetries = 3
	}

	curr_tries := 1
	c := context.Background()

	for options.MaxRetries >= curr_tries {
		client := redis.NewClient(&redis.Options{
			Addr:     options.Addr,
			Password: options.Password,
			DB:       options.DB,
			PoolSize: options.PoolSize,
		})

		ctx, cancel := context.WithTimeout(c, 10*time.Second)
		defer cancel()

		if err := client.Ping(ctx).Err(); err != nil {
			log.Printf("[REDIS] Connection failed, retrying... (%d/%d)", curr_tries, options.MaxRetries)

			backoff := time.Duration(math.Pow(2, float64(curr_tries))) * time.Second

			time.Sleep(backoff)

			curr_tries++

			client.Close()
			continue
		}

		r.RDB = client
		r.ctx = c

		log.Printf("[REDIS] Connected to redis at %s", options.Addr)
		return nil
	}

	return fmt.Errorf("[REDIS] Unable to connect to redis after %d attempts", options.MaxRetries)
}

func (r *Redis) Close() error {
	if err := r.RDB.Close(); err != nil {
		log.Printf("[REDIS] Error closing redis connection: %s", err)
		return err
	}

	log.Printf("[REDIS] Connection closed")
	return nil
}

func (r *Redis) GetClient() *redis.Client {
	if r.RDB == nil {
		log.Printf("[REDIS] Redis not connected")
		if err := r.Connect(&types.RedisConnectionConfig{}); err != nil {
			log.Printf("[REDIS] Error connecting to redis: %s", err)
		}
	}

	return r.RDB
}
