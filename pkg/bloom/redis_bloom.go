package bloom

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Filter interface {
	Add(ctx context.Context, key string, value interface{}) error
	Exists(ctx context.Context, key string, value interface{}) (bool, error)
}

type redisBloom struct {
	client *redis.Client
}

func NewRedisBloom(client *redis.Client) Filter {
	return &redisBloom{client: client}
}

func (r *redisBloom) Add(ctx context.Context, key string, value interface{}) error {
	// 使用redis的BF.ADD命令
	return r.client.Do(ctx, "BF.ADD", key, value).Err()
}

func (r *redisBloom) Exists(ctx context.Context, key string, value interface{}) (bool, error) {
	res, err := r.client.Do(ctx, "BF.EXISTS", key, value).Int()
	if err != nil {
		return false, err
	}

	return res == 1, nil
}
