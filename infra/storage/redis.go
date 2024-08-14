package storage

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func InitRedis(addr string, port string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr + ":" + port,
	})
	return &Redis{
		client: client,
	}, nil
}

func (s *Redis) RegisterAccess(ctx context.Context, key string, value string, limit int) (bool, int64, error) {
	now := time.Now()
	expiry := now.Add(-time.Second)
	// PIPELINE
	pipeline := s.client.Pipeline()
	redisKey := formatRedisKey("access", key, value)
	pipeline.ZRemRangeByScore(ctx, redisKey, "0", strconv.FormatInt(expiry.Unix(), 10))
	count := pipeline.ZCard(ctx, redisKey)
	_, err := pipeline.Exec(ctx)
	if err != nil {
		log.Println("redis: error on pipeline", err)
		return false, 0, err
	}
	if count.Val() >= int64(limit) {
		return false, count.Val(), nil
	}
	pipeline = s.client.Pipeline()
	pipeline.ZAdd(ctx, redisKey, redis.Z{Score: float64(now.Unix()), Member: now.Format(time.RFC3339Nano)})
	pipeline.Expire(ctx, redisKey, time.Second)
	_, err = pipeline.Exec(ctx)
	if err != nil {
		log.Println("redis: error on pipeline", err)
		return false, 0, err
	}
	return true, count.Val() + 1, nil
}

func (s *Redis) Search(ctx context.Context, key string, value string) (*time.Time, error) {
	redisKey := formatRedisKey("block", key, value)
	blocked, err := s.client.Get(ctx, redisKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	blockingTime, err := strconv.ParseInt(blocked, 10, 64)
	if err != nil {
		log.Println("redis: error when analyzing blocking time", err)
		return nil, err
	}
	blockingTimeUnix := time.Unix(0, blockingTime)
	return &blockingTimeUnix, nil
}

func (s *Redis) Block(ctx context.Context, key string, value string, block int) (*time.Time, error) {
	redisKey := formatRedisKey("block", key, value)
	blockingTime := time.Now().Add(time.Duration(block))
	err := s.client.Set(ctx, redisKey, blockingTime.Unix(), time.Duration(block)*time.Second).Err()
	if err != nil {
		log.Println("redis: error registering a block", err)
		return nil, err
	}
	return &blockingTime, nil
}

func formatRedisKey(prefix string, key string, value string) string {
	return fmt.Sprintf("%s-%s-%s", prefix, key, value)
}
