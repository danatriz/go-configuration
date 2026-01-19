package test

import (
	"configurations/databases"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func TestRedis() (string, error) {
	// Setup configuration
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)

	// Initialize Redis
	err := databases.InitializeRedis()
	if err != nil {
		return "", fmt.Errorf("failed to initialize Redis: %w", err)
	}

	log.Println("Redis initialized successfully")

	client := databases.GetRedisClient()
	ctx := context.Background()

	// Test PING
	log.Println("=== PING ===")
	pingResult, err := client.Ping(ctx).Result()
	if err != nil {
		return "", fmt.Errorf("failed to ping Redis: %w", err)
	}
	log.Printf("PING response: %s", pingResult)

	// Get server info
	log.Println("\n=== Server Info ===")
	info, err := client.Info(ctx, "server").Result()
	if err != nil {
		return "", fmt.Errorf("failed to get server info: %w", err)
	}
	log.Printf("Server info:\n%s", info)

	// String operations
	log.Println("\n=== String Operations ===")
	err = client.Set(ctx, "name", "Alice", 1*time.Hour).Err()
	if err != nil {
		return "", fmt.Errorf("failed to set name: %w", err)
	}
	log.Println("SET name = Alice")

	name, err := client.Get(ctx, "name").Result()
	if err != nil {
		return "", fmt.Errorf("failed to get name: %w", err)
	}
	log.Printf("GET name = %s", name)

	// Set with expiration
	err = client.Set(ctx, "temp_key", "will_expire", 2*time.Second).Err()
	if err != nil {
		return "", fmt.Errorf("failed to set temp_key: %w", err)
	}
	log.Println("SET temp_key with 2s expiration")

	// MGET
	err = client.Set(ctx, "key1", "value1", 1*time.Hour).Err()
	if err != nil {
		return "", fmt.Errorf("failed to set key1: %w", err)
	}
	err = client.Set(ctx, "key2", "value2", 1*time.Hour).Err()
	if err != nil {
		return "", fmt.Errorf("failed to set key2: %w", err)
	}
	values, err := client.MGet(ctx, "key1", "key2", "nonexistent").Result()
	if err != nil {
		return "", fmt.Errorf("failed to mget: %w", err)
	}
	log.Printf("MGET key1, key2, nonexistent = %v", values)

	// Hash operations
	log.Println("\n=== Hash Operations ===")
	hashKey := "user:1001"
	err = client.HSet(ctx, hashKey, "name", "Bob").Err()
	if err != nil {
		return "", fmt.Errorf("failed to hset name: %w", err)
	}
	err = client.HSet(ctx, hashKey, "email", "bob@example.com").Err()
	if err != nil {
		return "", fmt.Errorf("failed to hset email: %w", err)
	}
	err = client.HSet(ctx, hashKey, "age", 30).Err()
	if err != nil {
		return "", fmt.Errorf("failed to hset age: %w", err)
	}
	log.Printf("HSET %s name, email, age", hashKey)

	hashName, _ := client.HGet(ctx, hashKey, "name").Result()
	hashEmail, _ := client.HGet(ctx, hashKey, "email").Result()
	hashAge, _ := client.HGet(ctx, hashKey, "age").Result()
	log.Printf("HGET %s: name=%s, email=%s, age=%s", hashKey, hashName, hashEmail, hashAge)

	// HGETALL
	allFields, err := client.HGetAll(ctx, hashKey).Result()
	if err != nil {
		return "", fmt.Errorf("failed to hgetall: %w", err)
	}
	log.Printf("HGETALL %s = %v", hashKey, allFields)

	// List operations
	log.Println("\n=== List Operations ===")
	listKey := "mylist"
	err = client.LPush(ctx, listKey, "item1").Err()
	if err != nil {
		return "", fmt.Errorf("failed to lpush item1: %w", err)
	}
	err = client.LPush(ctx, listKey, "item2").Err()
	if err != nil {
		return "", fmt.Errorf("failed to lpush item2: %w", err)
	}
	err = client.LPush(ctx, listKey, "item3").Err()
	if err != nil {
		return "", fmt.Errorf("failed to lpush item3: %w", err)
	}
	log.Printf("LPUSH %s item1, item2, item3", listKey)

	listLen, err := client.LLen(ctx, listKey).Result()
	if err != nil {
		return "", fmt.Errorf("failed to llen: %w", err)
	}
	log.Printf("LLEN %s = %d", listKey, listLen)

	listItems, err := client.LRange(ctx, listKey, 0, -1).Result()
	if err != nil {
		return "", fmt.Errorf("failed to lrange: %w", err)
	}
	log.Printf("LRANGE %s 0 -1 = %v", listKey, listItems)

	// Set operations
	log.Println("\n=== Set Operations ===")
	setKey := "myset"
	err = client.SAdd(ctx, setKey, "apple", "banana", "cherry").Err()
	if err != nil {
		return "", fmt.Errorf("failed to sadd: %w", err)
	}
	log.Printf("SADD %s apple, banana, cherry", setKey)

	setMembers, err := client.SMembers(ctx, setKey).Result()
	if err != nil {
		return "", fmt.Errorf("failed to smembers: %w", err)
	}
	log.Printf("SMEMBERS %s = %v", setKey, setMembers)

	setLen, err := client.SCard(ctx, setKey).Result()
	if err != nil {
		return "", fmt.Errorf("failed to scard: %w", err)
	}
	log.Printf("SCARD %s = %d", setKey, setLen)

	// Sorted Set operations
	log.Println("\n=== Sorted Set Operations ===")
	zsetKey := "leaderboard"
	err = client.ZAdd(ctx, zsetKey, redis.Z{Score: 100, Member: "player1"},
		redis.Z{Score: 200, Member: "player2"},
		redis.Z{Score: 150, Member: "player3"}).Err()
	if err != nil {
		return "", fmt.Errorf("failed to zadd: %w", err)
	}
	log.Printf("ZADD %s player1=100, player2=200, player3=150", zsetKey)

	zRange, err := client.ZRangeWithScores(ctx, zsetKey, 0, -1).Result()
	if err != nil {
		return "", fmt.Errorf("failed to zrange: %w", err)
	}
	log.Printf("ZRANGE %s WITHSCORES:", zsetKey)
	for _, z := range zRange {
		log.Printf("  %s: %.0f", z.Member, z.Score)
	}

	// Increment operations
	log.Println("\n=== Increment Operations ===")
	counterKey := "counter"
	err = client.Set(ctx, counterKey, 0, 0).Err()
	if err != nil {
		return "", fmt.Errorf("failed to set counter: %w", err)
	}

	val, err := client.Incr(ctx, counterKey).Result()
	if err != nil {
		return "", fmt.Errorf("failed to incr: %w", err)
	}
	log.Printf("INCR %s = %d", counterKey, val)

	val, err = client.IncrBy(ctx, counterKey, 5).Result()
	if err != nil {
		return "", fmt.Errorf("failed to incrby: %w", err)
	}
	log.Printf("INCRBY %s 5 = %d", counterKey, val)

	// Exists operation
	log.Println("\n=== Exists Operation ===")
	exists, err := client.Exists(ctx, "name", "key1", "nonexistent_key").Result()
	if err != nil {
		return "", fmt.Errorf("failed to exists: %w", err)
	}
	log.Printf("EXISTS name, key1, nonexistent_key = %d", exists)

	// Expire operation
	log.Println("\n=== Expire Operation ===")
	err = client.Expire(ctx, "name", 10*time.Minute).Err()
	if err != nil {
		return "", fmt.Errorf("failed to expire: %w", err)
	}
	ttl, err := client.TTL(ctx, "name").Result()
	if err != nil {
		return "", fmt.Errorf("failed to ttl: %w", err)
	}
	log.Printf("EXPIRE name 10m, TTL name = %v", ttl)

	// Delete operation
	log.Println("\n=== Delete Operation ===")
	deleted, err := client.Del(ctx, "key1", "key2").Result()
	if err != nil {
		return "", fmt.Errorf("failed to del: %w", err)
	}
	log.Printf("DEL key1, key2 = %d keys deleted", deleted)

	// Flush test keys
	testKeys := []string{"name", "temp_key", hashKey, listKey, setKey, zsetKey, counterKey}
	deleted, err = client.Del(ctx, testKeys...).Result()
	if err != nil {
		log.Printf("Warning: failed to delete test keys: %v", err)
	} else {
		log.Printf("\nCleaned up %d test keys", deleted)
	}

	// Get DB size
	dbSize, err := client.DBSize(ctx).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get dbsize: %w", err)
	}
	log.Printf("Current DB size: %d keys", dbSize)

	// Close connection
	err = databases.CloseRedis()
	if err != nil {
		return "", fmt.Errorf("failed to close Redis: %w", err)
	}

	result := "Redis test completed!\nAll operations successful!"
	log.Println(result)

	return result, nil
}
