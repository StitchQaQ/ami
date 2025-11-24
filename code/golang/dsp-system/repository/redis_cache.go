package repository

import (
	"context"
	"dsp-system/config"
	"dsp-system/rpc"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache Redis缓存操作
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建Redis缓存客户端
func NewRedisCache(cfg *config.RedisConfig) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Redis连接失败: %v", err)
	} else {
		log.Printf("Redis连接成功: %s:%s", cfg.Host, cfg.Port)
	}

	return &RedisCache{
		client: client,
	}
}

// Close 关闭Redis连接
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// GetUserProfile 获取用户画像缓存
func (r *RedisCache) GetUserProfile(ctx context.Context, key string) *rpc.UserProfile {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil // 缓存不存在
	} else if err != nil {
		log.Printf("Redis获取失败: %v", err)
		return nil
	}

	var profile rpc.UserProfile
	if err := json.Unmarshal([]byte(val), &profile); err != nil {
		log.Printf("反序列化用户画像失败: %v", err)
		return nil
	}

	return &profile
}

// SetUserProfile 设置用户画像缓存
func (r *RedisCache) SetUserProfile(ctx context.Context, key string, profile *rpc.UserProfile, expiration time.Duration) error {
	data, err := json.Marshal(profile)
	if err != nil {
		return fmt.Errorf("序列化用户画像失败: %v", err)
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

// GetAdCache 获取广告缓存
func (r *RedisCache) GetAdCache(ctx context.Context, adID string) (string, error) {
	key := fmt.Sprintf("ad:%s", adID)
	return r.client.Get(ctx, key).Result()
}

// SetAdCache 设置广告缓存
func (r *RedisCache) SetAdCache(ctx context.Context, adID string, creative string, expiration time.Duration) error {
	key := fmt.Sprintf("ad:%s", adID)
	return r.client.Set(ctx, key, creative, expiration).Err()
}

// IncrBidCount 增加竞价计数
func (r *RedisCache) IncrBidCount(ctx context.Context, date string) (int64, error) {
	key := fmt.Sprintf("stats:bid_count:%s", date)
	return r.client.Incr(ctx, key).Result()
}

// IncrWinCount 增加赢标计数
func (r *RedisCache) IncrWinCount(ctx context.Context, date string) (int64, error) {
	key := fmt.Sprintf("stats:win_count:%s", date)
	return r.client.Incr(ctx, key).Result()
}

// GetStats 获取统计数据
func (r *RedisCache) GetStats(ctx context.Context, date string) (bidCount int64, winCount int64, err error) {
	bidKey := fmt.Sprintf("stats:bid_count:%s", date)
	winKey := fmt.Sprintf("stats:win_count:%s", date)

	bidCount, _ = r.client.Get(ctx, bidKey).Int64()
	winCount, _ = r.client.Get(ctx, winKey).Int64()

	return bidCount, winCount, nil
}

// SetFrequencyCap 设置频次控制
func (r *RedisCache) SetFrequencyCap(ctx context.Context, userID string, adID string, expiration time.Duration) error {
	key := fmt.Sprintf("freq_cap:%s:%s", userID, adID)
	return r.client.Set(ctx, key, "1", expiration).Err()
}

// CheckFrequencyCap 检查频次控制
func (r *RedisCache) CheckFrequencyCap(ctx context.Context, userID string, adID string) bool {
	key := fmt.Sprintf("freq_cap:%s:%s", userID, adID)
	_, err := r.client.Get(ctx, key).Result()
	return err == redis.Nil // 如果不存在，说明没有超过频次限制
}

// CacheBudget 缓存预算信息
func (r *RedisCache) CacheBudget(ctx context.Context, campaignID string, remaining float64, expiration time.Duration) error {
	key := fmt.Sprintf("budget:%s", campaignID)
	return r.client.Set(ctx, key, remaining, expiration).Err()
}

// GetCachedBudget 获取缓存的预算信息
func (r *RedisCache) GetCachedBudget(ctx context.Context, campaignID string) (float64, error) {
	key := fmt.Sprintf("budget:%s", campaignID)
	return r.client.Get(ctx, key).Float64()
}

