package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis(config *Config) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         config.RedisAddr,
		Password:     config.RedisPassword,
		DB:           config.RedisDB,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 2,
	})

	ctx, cancel := context.WithTimeout(Ctx, 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("❌ Redis connection failed:", err)
	}

	log.Println("✅ Connected to Redis!")
}

func GetUserFromCacheOrDB(userID int) (*User, error) {
	start := time.Now() // ⏱️ Start timer

	key := fmt.Sprintf("user:%d", userID)

	// 1. Try Redis cache
	val, err := RedisClient.Get(context.Background(), key).Result()
	if err == nil {
		log.Printf("🔁 Cache hit for %s, took %v", key, time.Since(start))
		var user User
		json.Unmarshal([]byte(val), &user)
		return &user, nil
	}

	// 2. Cache miss → DB
	log.Printf("💾 Cache miss for %s, querying DB...", key)
	var user User
	err = DB.QueryRow("SELECT id, name, email FROM users WHERE id = $1", userID).
		Scan(&user.ID, &user.Name, &user.Email)

	if err != nil {
		return nil, err
	}

	// 3. Save to Redis
	userJSON, _ := json.Marshal(user)
	RedisClient.Set(context.Background(), key, userJSON, 10*time.Minute)

	// 🕐 Log how long DB + cache took
	log.Printf("✅ User %d loaded from DB and cached, took %v", userID, time.Since(start))

	return &user, nil
}
