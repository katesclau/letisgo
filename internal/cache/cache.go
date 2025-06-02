package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
    client *redis.Client
    ctx    context.Context
}

func NewRedisCache(addr, password string, db int) *RedisCache {
    rdb := redis.NewClient(&redis.Options{
        Addr:     addr,
        Password: password,
        DB:       db,
    })
    return &RedisCache{
        client: rdb,
        ctx:    context.Background(),
    }
}

func (r *RedisCache) Get(pk, sk string) (string, error) {
    key := pk + ":" + sk
    val, err := r.client.Get(r.ctx, key).Result()
    if err == redis.Nil {
        return "", nil
    }
    if err != nil {
        return "", err
    }
    return val, nil
}

func (r *RedisCache) Set(pk, sk, value string) error {
    key := pk + ":" + sk
    err := r.client.Set(r.ctx, key, value, 0).Err()
    if err != nil {
        return err
    }
    return nil
}
