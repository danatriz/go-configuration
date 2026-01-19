package databases

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var RedisClient *redis.Client

func InitializeRedis() error {
	host := viper.GetString("redis.host")
	port := viper.GetInt("redis.port")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")
	poolSize := viper.GetInt("redis.pool_size")
	minIdleConns := viper.GetInt("redis.min_idle_conns")
	maxRetries := viper.GetInt("redis.max_retries")
	dialTimeout := viper.GetDuration("redis.dial_timeout")
	readTimeout := viper.GetDuration("redis.read_timeout")
	writeTimeout := viper.GetDuration("redis.write_timeout")
	poolTimeout := viper.GetDuration("redis.pool_timeout")

	if port == 0 {
		port = 6379
	}
	if db == 0 {
		db = 0
	}
	if poolSize == 0 {
		poolSize = 10
	}
	if minIdleConns == 0 {
		minIdleConns = 5
	}
	if maxRetries == 0 {
		maxRetries = 3
	}
	if dialTimeout == 0 {
		dialTimeout = 5 * time.Second
	}
	if readTimeout == 0 {
		readTimeout = 3 * time.Second
	}
	if writeTimeout == 0 {
		writeTimeout = 3 * time.Second
	}
	if poolTimeout == 0 {
		poolTimeout = 4 * time.Second
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Password:     password,
		DB:           db,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
		MaxRetries:   maxRetries,
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		PoolTimeout:  poolTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	RedisClient = client
	log.Printf("Connected to Redis at %s:%d, DB: %d", host, port, db)
	return nil
}

func GetRedisClient() *redis.Client {
	return RedisClient
}

func CloseRedis() error {
	if RedisClient != nil {
		err := RedisClient.Close()
		if err != nil {
			return err
		}
		log.Println("Redis connection closed")
	}
	return nil
}

func NewRedisClient() (*redis.Client, error) {
	host := viper.GetString("redis.host")
	port := viper.GetInt("redis.port")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")
	poolSize := viper.GetInt("redis.pool_size")
	minIdleConns := viper.GetInt("redis.min_idle_conns")
	maxRetries := viper.GetInt("redis.max_retries")
	dialTimeout := viper.GetDuration("redis.dial_timeout")
	readTimeout := viper.GetDuration("redis.read_timeout")
	writeTimeout := viper.GetDuration("redis.write_timeout")
	poolTimeout := viper.GetDuration("redis.pool_timeout")

	if port == 0 {
		port = 6379
	}
	if db == 0 {
		db = 0
	}
	if poolSize == 0 {
		poolSize = 10
	}
	if minIdleConns == 0 {
		minIdleConns = 5
	}
	if maxRetries == 0 {
		maxRetries = 3
	}
	if dialTimeout == 0 {
		dialTimeout = 5 * time.Second
	}
	if readTimeout == 0 {
		readTimeout = 3 * time.Second
	}
	if writeTimeout == 0 {
		writeTimeout = 3 * time.Second
	}
	if poolTimeout == 0 {
		poolTimeout = 4 * time.Second
	}

	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Password:     password,
		DB:           db,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
		MaxRetries:   maxRetries,
		DialTimeout:  dialTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		PoolTimeout:  poolTimeout,
	})

	return client, nil
}
