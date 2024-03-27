package core

import (
    "context"
    "github.com/redis/go-redis/v9"
)

type Redis struct {
    client *redis.Client
}

// RedisDal allows testing with mock
type RedisDal interface {
    HMGet(ctx context.Context, key string, fields ...string) (results map[string]interface{}, err error)
    HSet(ctx context.Context, key string, values ...interface{}) error
}

func NewRedis(client *redis.Client) *Redis {
    return &Redis{client: client}
}

func (s Redis) HMGet(ctx context.Context, key string, fields ...string) (results map[string]interface{}, err error) {
    list, err := s.client.HMGet(ctx, key, fields...).Result()
    if err != nil {
        return
    }
    results = map[string]interface{}{}
    for i, k := range fields {
        if list[i] != nil {
            results[k] = list[i]
        }
    }
    return
}

func (s Redis) HSet(ctx context.Context, key string, values ...interface{}) error {
    return s.client.HSet(ctx, key, values...).Err()
}
